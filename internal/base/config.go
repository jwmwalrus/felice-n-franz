package base

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

const (
	// CurrentConfigFileVersion current or default config file version
	CurrentConfigFileVersion = 1
)

// Config application's configuration
type Config struct {
	Version  int           `json:"version"`
	FirstRun bool          `json:"firstRun"`
	Port     int           `json:"port"`
	Envs     []Environment `json:"envs"`
}

func (c *Config) setDefaults() {
	log.Info("Setting config defaults")

	if c.Version == 0 {
		c.Version = CurrentConfigFileVersion
	}

	if c.Port == 0 {
		c.Port = DefaultPort
	}

	if len(c.Envs) < 1 {
		c.Envs = append(c.Envs, Environment{})
		c.Envs[0].setDefaults()
	}
}

// GetPort returns port as a string
func (c *Config) GetPort() string {
	return ":" + strconv.Itoa(c.Port)
}

// GetURL returns application's base URL
func (c *Config) GetURL() string {
	return "http://localhost" + c.GetPort()
}

// GetEnvConfig returns the configuration for the given environment
func (c *Config) GetEnvConfig(envName string) (env Environment) {
	for _, e := range c.Envs {
		if !e.Active {
			continue
		}
		if e.Name == envName {
			env = e
			break
		}
	}
	return
}

// GetEnvsList returns the list of availablr environments
func (c *Config) GetEnvsList() (list []string) {
	for _, e := range c.Envs {
		if !e.Active {
			continue
		}
		list = append(list, e.Name)
	}
	return
}

// GetEnvTopics returns the list of topics for the given environment`
func (c *Config) GetEnvTopics(envName string) (topics []Topic) {
	for _, e := range c.Envs {
		if e.Name == envName {
			topics = e.Topics
			break
		}
	}
	return
}

// GetEnvGroups returns the list of groups in the given environment
func (c *Config) GetEnvGroups(envName string) (groups []Group) {
	for _, e := range c.Envs {
		if e.Name == envName {
			groups = e.Groups
			break
		}
	}
	return
}

// GetEnvGroupTopics returns the list of topics in a given group, for a given environment
func (c *Config) GetEnvGroupTopics(envName, groupID string) (topics []Topic) {
	for _, e := range c.Envs {
		if e.Name != envName {
			continue
		}

		var keys []string
		for _, g := range e.Groups {
			if g.ID == groupID {
				keys = g.Keys
				break
			}
		}

		for _, k := range keys {
			for _, t := range e.Topics {
				if t.Key == k {
					topics = append(topics, t)
					break
				}
			}
		}

		break
	}
	return
}

func (c *Config) validate() (err error) {
	for i := range c.Envs {
		if err = c.Envs[i].validate(); err != nil {
			return
		}
	}
	return
}

// Environment defines a topics environment
// FIXME: DEPRECATE AssignConsumer
type Environment struct {
	Name           string          `json:"name"`
	Active         bool            `json:"active"`
	AssignConsumer bool            `json:"assignConsumer"`
	Configuration  kafka.ConfigMap `json:"configuration"`
	Vars           EnvVars         `json:"vars"`
	HeaderPrefix   string          `json:"headerPrefix"`
	TopicsFrom     string          `json:"inheritFrom"`
	Schemas        EnvSchemas      `json:"schemas"`
	Topics         []Topic         `json:"topics"`
	Groups         []Group         `json:"groups"`
}

// EnvVars defines the Environment Vars type
type EnvVars map[string]string

// EnvSchemas defines common JSON schemas
type EnvSchemas map[string]interface{}

// AllTopicsExist check if all topics exist for the given keys
func (e *Environment) AllTopicsExist(keys []string) (exist bool) {
	_, err := e.FindTopics(keys)
	exist = err == nil
	return
}

// FindTopic looks for topic by key
func (e *Environment) FindTopic(key string) (v Topic, err error) {
	for _, t := range e.Topics {
		if key == t.Key {
			v = t
			break
		}
	}

	if v.Key == "" {
		err = errors.New("Key not found")
	}
	return
}

// FindTopics looks for all topics in the given keys list
func (e *Environment) FindTopics(keys []string) (topics []Topic, err error) {
	for _, k := range keys {
		var v Topic
		v, err = e.FindTopic(k)
		if err != nil {
			return
		}
		topics = append(topics, v)
	}
	return
}

// FindTopicValues looks for all topics in the given keys list
func (e *Environment) FindTopicValues(keys []string) (values []string, err error) {
	var topics []Topic
	if topics, err = e.FindTopics(keys); err != nil {
		return
	}
	values = getValues(topics)
	return
}

// GetTopic returns the topic matching the given value
func (e *Environment) GetTopic(value string) (v Topic, err error) {
	for _, t := range e.Topics {
		if value == t.Value {
			v = t
			break
		}
	}

	if v.Key == "" {
		err = errors.New("Value not found")
	}
	return
}

func (e *Environment) setDefaults() {
	log.Info("Setting env defaults")

	if e.Name == "" {
		e.Name = "default"
	}

	if e.Configuration == nil {
		e.Configuration = kafka.ConfigMap{}
	}

	if e.Vars == nil {
		e.Vars = EnvVars{}
		e.Vars.setDefaults()
	}

	if e.Schemas == nil {
		e.Schemas = EnvSchemas{}
	}

	if len(e.Topics) < 1 {
		e.Topics = append(e.Topics, Topic{})
		e.Topics[0].setDefaults()
	}

	if len(e.Groups) < 1 {
		e.Groups = append(e.Groups, Group{})
		e.Groups[0].setGroupDefaults()
	}
}

func (e *Environment) setup() {
	log.Info("Setting environment " + e.Name)

	e.Name = os.ExpandEnv(e.Name)
	e.TopicsFrom = os.ExpandEnv(e.TopicsFrom)

	if _, ok := e.Configuration["group.id"]; !ok {
		e.Configuration["group.id"] = "${USER}.${HOSTNAME}"
	}

	for k, v := range e.Configuration {
		if reflect.TypeOf(v).Name() != "string" {
			continue
		}
		e.Configuration[k] = os.ExpandEnv(v.(string))
	}

	fromIdx := -1
	if e.TopicsFrom != "" {
		for j := range Conf.Envs {
			if Conf.Envs[j].Name == e.TopicsFrom {
				fromIdx = j
				break
			}
		}
	}

	if fromIdx < 0 {
		return
	}

	e.Topics = nil
	e.Topics = make([]Topic, len(Conf.Envs[fromIdx].Topics))
	copy(e.Topics, Conf.Envs[fromIdx].Topics)
	for i := range e.Topics {
		if e.Topics[i].Headers == nil {
			e.Topics[i].Headers = []Header{}
		}
	}

	e.Groups = nil
	e.Groups = make([]Group, len(Conf.Envs[fromIdx].Groups))
	copy(e.Groups, Conf.Envs[fromIdx].Groups)
}

func (e *Environment) validate() (err error) {
	if e.Name == "" {
		err = errors.New("Environment name cannot be empty")
		return
	}

	if e.TopicsFrom != "" {
		return
	}

	if e.Schemas == nil {
		e.Schemas = EnvSchemas{}
	}

	for i := range e.Topics {
		if err = e.Topics[i].validate(); err != nil {
			return
		}
	}

	for i := range e.Groups {
		if err = e.Groups[i].validate(); err != nil {
			return
		}
	}

	unique := make(map[string]string)

	n := 0
	for _, t := range e.Topics {
		n++
		if u, ok := unique[t.Key]; ok {
			err = fmt.Errorf("Key %v is already in use by topic/group (#%v)", t.Key, u)
			return
		}
		unique[t.Key] = t.Name
	}
	log.Infof("Found %v unique topic keys for environment %v", n, e.Name)

	n = 0
	for _, g := range e.Groups {
		n++
		if u, ok := unique[g.ID]; ok {
			err = fmt.Errorf("ID %v is already in use by topic/group (#%v)", g.ID, u)
			return
		}
		unique[g.ID] = g.Name
	}
	log.Infof("Found %v unique group IDs for environment %v", n, e.Name)
	return
}

func (g *EnvVars) setDefaults() {
	*g = EnvVars{"myVar": "some-value"}
}

// Topic defines a topic structure
type Topic struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Key         string      `json:"key"`
	Value       string      `json:"value"`
	GroupID     string      `json:"groupId"`
	Headers     []Header    `json:"headers"`
	Schema      TopicSchema `json:"schema"`
}

// TopicSchema defines the applicable JSON schema
type TopicSchema map[string]interface{}

func (t *Topic) expandVars(vars EnvVars) {
	for k, v := range vars {
		t.Value = strings.ReplaceAll(t.Value, "{{"+k+"}}", v)
		if t.GroupID != "" {
			t.GroupID = strings.ReplaceAll(t.GroupID, "{{"+k+"}}", v)
		}
	}
}

func (t *Topic) setDefaults() {
	*t = Topic{
		Name:        "PaymentTopic",
		Description: "Handles payment events",
		Key:         "payment",
		Value:       "{{myVar}}.division.department.section.subsection.payment-type",
		Headers:     []Header{},
		Schema:      TopicSchema{},
	}
}

func (t *Topic) validate() (err error) {
	if t.Key == "" || t.Value == "" {
		err = errors.New("Topic key and value mus be non-empty")
		return
	}

	if t.Headers == nil {
		t.Headers = []Header{}
	}

	if t.Schema == nil {
		t.Schema = TopicSchema{}
	}
	return
}

func getValues(topics []Topic) (values []string) {
	for _, t := range topics {
		values = append(values, t.Value)
	}
	return
}

// Group defines a group of topics
type Group struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	ID          string   `json:"id"`
	Keys        []string `json:"keys"`
}

func (g *Group) setGroupDefaults() {
	*g = Group{
		Name:        "GroupOne",
		Description: "Group One",
		Category:    "Category A",
		ID:          "group1",
		Keys:        []string{"payment"},
	}
}

func (g *Group) validate() (err error) {
	if g.ID == "" || len(g.Keys) == 0 {
		err = errors.New("Group id and keys mus be non-empty")
		return
	}
	return
}

// Header defines the Headers type
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

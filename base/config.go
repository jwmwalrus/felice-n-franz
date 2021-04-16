package base

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Config Application's configuration
type Config struct {
	FirstRun      bool          `json:"firstRun"`
	Port          int           `json:"port"`
	MaxTailOffset int           `json:"maxTailOffset"`
	Envs          []Environment `json:"envs"`
}

func (rec *Config) setDefaults() {
	if rec.Port == 0 {
		rec.Port = DefaultPort
	}

	if rec.MaxTailOffset == 0 {
		rec.MaxTailOffset = DefaultTailOffset
	}

	if len(rec.Envs) < 1 {
		rec.Envs = append(rec.Envs, Environment{})
		rec.Envs[0].setDefaults()
	}
}

// GetPort Returns port as a string
func (rec *Config) GetPort() string {
	return ":" + strconv.Itoa(rec.Port)
}

// GetURL Returns application's base URL
func (rec *Config) GetURL() string {
	return "http://localhost" + rec.GetPort()
}

// GetEnvConfig returns the configuration for the given environment
func (rec *Config) GetEnvConfig(envName string) (env Environment) {
	for _, e := range rec.Envs {
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
func (rec *Config) GetEnvsList() (list []string) {
	for _, e := range rec.Envs {
		if !e.Active {
			continue
		}
		list = append(list, e.Name)
	}
	return
}

// GetEnvTopics returns the list of topics for the given environment`
func (rec *Config) GetEnvTopics(envName string) (topics []Topic) {
	for _, e := range rec.Envs {
		if e.Name == envName {
			topics = e.Topics
			break
		}
	}
	return
}

// GetEnvGroups returns the list of groups in the given environment
func (rec *Config) GetEnvGroups(envName string) (groups []Group) {
	for _, e := range rec.Envs {
		if e.Name == envName {
			groups = e.Groups
			break
		}
	}
	return
}

// GetEnvGroupTopics returns the list of topics in a given group, for a given environment
func (rec *Config) GetEnvGroupTopics(envName, groupID string) (topics []Topic) {
	for _, e := range rec.Envs {
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

// Environment Defines a topics environment
type Environment struct {
	Name          string          `json:"name"`
	Active        bool            `json:"active"`
	Configuration kafka.ConfigMap `json:"configuration"`
	Vars          EnvVars         `json:"vars"`
	TopicsFrom    string          `json:"inheritFrom"`
	Topics        []Topic         `json:"topics"`
	Groups        []Group         `json:"grups"`
}

// EnvVars defines the Environment Vars type
type EnvVars map[string]string

// AllTopicsExist check if all topics exist for the given keys
func (rec *Environment) AllTopicsExist(keys []string) (exist bool) {
	_, err := rec.FindTopics(keys)
	exist = err == nil
	return
}

// FindTopic looks for topic by key
func (rec *Environment) FindTopic(key string) (v Topic, err error) {
	for _, t := range rec.Topics {
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
func (rec *Environment) FindTopics(keys []string) (topics []Topic, err error) {
	v := Topic{}
	for _, k := range keys {
		v, err = rec.FindTopic(k)
		if err != nil {
			return
		}
		topics = append(topics, v)
	}
	return
}

// FindTopicValues looks for all topics in the given keys list
func (rec *Environment) FindTopicValues(keys []string) (values []string, err error) {
	topics := []Topic{}
	if topics, err = rec.FindTopics(keys); err != nil {
		return
	}
	values = getValues(topics)
	return
}

func (rec *Environment) setDefaults() {
	if rec.Configuration == nil {
		rec.Configuration = kafka.ConfigMap{}
	}

	if rec.Vars == nil {
		rec.Vars = EnvVars{}
		rec.Vars.setDefaults()
	}

	if len(rec.Topics) < 1 {
		rec.Topics = append(rec.Topics, Topic{})
		rec.Topics[0].setDefaults()
	}

	if len(rec.Groups) < 1 {
		rec.Groups = append(rec.Groups, Group{})
		rec.Groups[0].setGroupDefaults()
	}
}

func (rec *Environment) setup() {
	rec.Name = os.ExpandEnv(rec.Name)
	rec.TopicsFrom = os.ExpandEnv(rec.TopicsFrom)

	if _, ok := rec.Configuration["group.id"]; !ok {
		rec.Configuration["group.id"] = "${USER}.${HOSTNAME}"
	}

	for k, v := range rec.Configuration {
		if reflect.TypeOf(v).Name() != "string" {
			continue
		}
		rec.Configuration[k] = os.ExpandEnv(v.(string))
	}

	fromIdx := -1
	if rec.TopicsFrom != "" {
		for j := range Conf.Envs {
			if Conf.Envs[j].Name == rec.TopicsFrom {
				fromIdx = j
				break
			}
		}
	}

	if fromIdx < 0 {
		// TODO: warn?
		return
	}

	rec.Topics = nil
	rec.Topics = make([]Topic, len(Conf.Envs[fromIdx].Topics))
	copy(rec.Topics, Conf.Envs[fromIdx].Topics)
	for i := range rec.Topics {
		if rec.Topics[i].Headers == nil {
			rec.Topics[i].Headers = TopicHeaders{}
		}
	}

	rec.Groups = nil
	rec.Groups = make([]Group, len(Conf.Envs[fromIdx].Groups))
	copy(rec.Groups, Conf.Envs[fromIdx].Groups)
}

// Group defines a group of topics
type Group struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Keys        []string `json:"keys"`
}

func (rec *Group) setGroupDefaults() {
	*rec = Group{
		ID:          "group1",
		Description: "Group One",
		Category:    "Category A",
		Keys:        []string{"payment"},
	}
}

func (rec *EnvVars) setDefaults() {
	*rec = EnvVars{"myVar": "some-value"}
}

// Topic Defines a key-value pair
type Topic struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Key         string       `json:"key"`
	Value       string       `json:"value"`
	Headers     TopicHeaders `json:headers`
	Schema      string       `json:"schema"`
}

// TopicHeaders defines the Topic Headers type
type TopicHeaders map[string]string

func (rec *Topic) expandVars(vars EnvVars) {
	for k, v := range vars {
		rec.Value = strings.ReplaceAll(rec.Value, "{{"+k+"}}", v)
	}
}

func (rec *Topic) setDefaults() {
	*rec = Topic{
		Name:        "PaymentTopic",
		Description: "Handles payment events",
		Key:         "payment",
		Value:       "{{myVar}}.division.department.section.subsection.payment-type",
		Headers:     TopicHeaders{},
	}
}

func getKeys(topics []Topic) (keys []string) {
	for _, t := range topics {
		keys = append(keys, t.Key)
	}
	return
}

func getValues(topics []Topic) (values []string) {
	for _, t := range topics {
		values = append(values, t.Value)
	}
	return
}

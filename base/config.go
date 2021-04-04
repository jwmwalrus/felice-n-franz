package base

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
func (rec *Config) GetEnvTopics(envName string) (topics []KeyVal) {
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
func (rec *Config) GetEnvGroupTopics(envName, groupName string) (topics []KeyVal) {
	for _, e := range rec.Envs {
		if e.Name != envName {
			continue
		}

		var keys []string
		for _, g := range e.Groups {
			if g.Name == groupName {
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
	Name              string   `json:"name"`
	Active            bool     `json:"active"`
	Broker            string   `json:"broker"`
	SecurityProtocol  string   `json:"securityProtocol"`
	SASLUsername      string   `json:"saslUsername"`
	SASLPassword      string   `json:"saslPassword"`
	SASLMechanism     string   `json:"saslMechanism"`
	SSLCALocation     string   `json:"sslCaLocation"`
	APIVersionRequest bool     `json:"apiVersionRequest"`
	Vars              []KeyVal `json:"vars"`
	TopicsFrom        string   `json:"inheritFrom"`
	Topics            []KeyVal `json:"topics"`
	Groups            []Group  `json:"grups"`
}

func (rec *Environment) setDefaults() {
	if rec.SecurityProtocol == "" {
		rec.SecurityProtocol = "SASL_SSL"
	}
	if rec.SASLMechanism == "" {
		rec.SASLMechanism = "PLAIN"
	}

	if rec.SSLCALocation == "" {
		caFile := filepath.Join(ConfigDir, DefaultCA)
		if _, err := os.Stat(caFile); !os.IsNotExist(err) {
			rec.SSLCALocation = caFile
		}

		if rec.SSLCALocation == "" {
			caFile = filepath.Join("$HOME", DefaultCA)
			if _, err := os.Stat(caFile); !os.IsNotExist(err) {
				rec.SSLCALocation = caFile
			}
		}

		if len(rec.Vars) < 1 {
			rec.Vars = append(rec.Vars, KeyVal{})
			rec.Vars[0].setVarDefaults()
		}

		if len(rec.Topics) < 1 {
			rec.Topics = append(rec.Topics, KeyVal{})
			rec.Topics[0].setTopicDefaults()
		}

		if len(rec.Groups) < 1 {
			rec.Groups = append(rec.Groups, Group{})
			rec.Groups[0].setGroupDefaults()
		}
	}
}

func (rec *Environment) setup() {
	rec.Name = os.ExpandEnv(rec.Name)
	rec.Broker = os.ExpandEnv(rec.Broker)
	rec.SecurityProtocol = os.ExpandEnv(rec.SecurityProtocol)
	rec.SASLUsername = os.ExpandEnv(rec.SASLUsername)
	rec.SASLPassword = os.ExpandEnv(rec.SASLPassword)
	rec.SASLMechanism = os.ExpandEnv(rec.SASLMechanism)
	rec.SSLCALocation = os.ExpandEnv(rec.SSLCALocation)
	rec.TopicsFrom = os.ExpandEnv(rec.TopicsFrom)

	fromIdx := -1
	if rec.TopicsFrom != "" {
		for j := 0; j < len(Conf.Envs); j++ {
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
	rec.Topics = make([]KeyVal, len(Conf.Envs[fromIdx].Topics))
	copy(rec.Topics, Conf.Envs[fromIdx].Topics)

	rec.Groups = nil
	rec.Groups = make([]Group, len(Conf.Envs[fromIdx].Groups))
	copy(rec.Groups, Conf.Envs[fromIdx].Groups)
}

// Group defines a group of topics
type Group struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Keys        []string `json:"keys"`
}

func (rec *Group) setGroupDefaults() {
	*rec = Group{
		Name:        "group1",
		Description: "Group One",
		Category:    "Category A",
		Keys:        []string{"payment"},
	}
}

// KeyVal Defines a key-value pair
type KeyVal struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (rec *KeyVal) expandVars(vars []KeyVal) {
	for _, kv := range vars {
		rec.Value = strings.ReplaceAll(rec.Value, "{{"+kv.Key+"}}", kv.Value)
	}
}

func (rec *KeyVal) setTopicDefaults() {
	*rec = KeyVal{
		Key:   "payment",
		Value: "{{myVar}}.division.department.section.subsection.payment-type",
	}
}

func (rec *KeyVal) setVarDefaults() {
	*rec = KeyVal{
		Key:   "myVar",
		Value: "some-value",
	}
}

func getKeys(kv []KeyVal) (keys []string) {
	for _, x := range kv {
		keys = append(keys, x.Key)
	}
	return
}

func getValues(kv []KeyVal) (values []string) {
	for _, x := range kv {
		values = append(values, x.Value)
	}
	return
}

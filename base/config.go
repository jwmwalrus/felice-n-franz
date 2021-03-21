package base

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jwmwalrus/bnp"
	"github.com/nightlyone/lockfile"
)

// Conf Application's global configuration
var Conf Config

// Load Loads application's configuration
func Load() (args []string) {
	args = ParseArgs()

	err := SetEnvDirs()
	bnp.PanicOnError(err)

	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		Conf.FirstRun = true
		err := Save()
		bnp.PanicOnError(err)
	}

	// var jsonFile *os.File
	jsonFile, err := os.Open(configFile)
	bnp.PanicOnError(err)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &Conf)

	return
}

// Save Saves application's configuration
func Save() (err error) {
	SetDefaults()

	var lock lockfile.Lockfile
	lock, err = lockfile.New(lockFile)
	if err != nil {
		return
	}

	if err = lock.TryLock(); err != nil {
		return
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			fmt.Printf("Cannot unlock %q, reason: %v", lock, err)
		}
	}()

	var file []byte
	file, err = json.MarshalIndent(Conf, "", " ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(configFile, file, 0644)

	return
}

// SetDefaults sets configuration defaults
func SetDefaults() {
	Conf.setDefaults()
}

// Unload Cleans up application before exit
func Unload() {
	Log.Info("Unloading application")
	if Conf.FirstRun {
		Conf.FirstRun = false

		// NOTE: ignore error
		Save()
	}
}

// Config Application's configuration
type Config struct {
	FirstRun bool          `json:"firstRun"`
	Port     int           `json:"port"`
	Envs     []Environment `json:"envs"`
}

func (rec *Config) setDefaults() {
	if rec.Port == 0 {
		rec.Port = DefaultPort
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

// Environment Defines a topics environment
type Environment struct {
	Name              string   `json:"name"`
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

// Group defines a group of topics
type Group struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keys        []string `json:"keys"`
}

func (rec *Group) setGroupDefaults() {
	*rec = Group{
		Name:        "group1",
		Description: "Group One",
		Keys:        []string{"payment"},
	}
}

// KeyVal Defines a key-value pair
type KeyVal struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

func (rec *KeyVal) setTopicDefaults() {
	*rec = KeyVal{
		Key: "payment",
		Val: "{{myVar}}.division.department.section.subsection.payment-type",
	}
}

func (rec *KeyVal) setVarDefaults() {
	*rec = KeyVal{
		Key: "myVar",
		Val: "some-value",
	}
}

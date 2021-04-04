package base

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jwmwalrus/bnp"
	"github.com/nightlyone/lockfile"
)

// Conf Application's global configuration
var Conf, userConf Config

// Load Loads application's configuration
func Load() (args []string) {
	args = ParseArgs()

	err := SetEnvDirs()
	bnp.PanicOnError(err)

	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		userConf.FirstRun = true
		err := Save()
		bnp.PanicOnError(err)
	}

	// var jsonFile *os.File
	jsonFile, err := os.Open(configFile)
	bnp.PanicOnError(err)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &userConf)
	bnp.PanicOnError(err)

	if userConf.FirstRun {
		userConf.FirstRun = false
		err := Save()
		bnp.PanicOnError(err)
	}

	copyUserConfig()

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
	file, err = json.MarshalIndent(userConf, "", "  ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(configFile, file, 0644)

	return
}

// SetDefaults sets configuration defaults
func SetDefaults() {
	userConf.setDefaults()
}

// Unload Cleans up application before exit
func Unload() {
	Log.Info("Unloading application")
	if userConf.FirstRun {
		userConf.FirstRun = false

		// NOTE: ignore error
		Save()
	}
}

func copyUserConfig() {
	Conf.FirstRun = userConf.FirstRun
	Conf.Port = userConf.Port
	Conf.Envs = make([]Environment, len(userConf.Envs))
	copy(Conf.Envs, userConf.Envs)

	for i := 0; i < len(Conf.Envs); i++ {
		Conf.Envs[i].setup()
	}

	for i := 0; i < len(Conf.Envs); i++ {
		for j := 0; j < len(Conf.Envs[i].Topics); j++ {
			Conf.Envs[i].Topics[j].expandVars(Conf.Envs[i].Vars)
		}
	}
}

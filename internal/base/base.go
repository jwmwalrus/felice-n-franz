package base

import (
	_ "embed" //required to import version.json
	"path/filepath"

	"github.com/jwmwalrus/bumpy/pkg/version"
	"github.com/jwmwalrus/onerror"
	rtc "github.com/jwmwalrus/rtcycler"
)

const (
	// AppDirName application's directory name
	AppDirName = "felice-n-franz"

	// AppName application name
	AppName = "Felice & Franz"

	// DefaultServerAPIVersion idem
	DefaultServerAPIVersion = "/api/v1"

	// DefaultCA Kafka Certificate's filename
	DefaultCA = "kafkacert.pem"

	// DefaultPort application's default port
	DefaultPort = 9191
)

var (
	//go:embed version.json
	versionJSON []byte

	// AppVersion application's version
	AppVersion version.Version

	// Conf application's global configuration
	Conf Config

	// UserConf application's global configuration, read during load
	UserConf Config
)

// Validate validates application's configuration
func Validate() *Config {
	err := AppVersion.Read(versionJSON)
	onerror.Panic(err)

	err = UserConf.validate()
	onerror.Panic(err)

	configFile := filepath.Join(rtc.ConfigDir(), rtc.ConfigFilename())
	if UserConf.FirstRun {
		UserConf.FirstRun = false
		err := rtc.SaveConfig(configFile)
		onerror.Panic(err)
	}

	copyUserConfig()

	return &Conf
}

// SetDefaults sets configuration defaults
func SetDefaults() {
	UserConf.SetDefaults()
}

func copyUserConfig() {
	Conf.Version = UserConf.Version
	Conf.FirstRun = UserConf.FirstRun
	Conf.Port = UserConf.Port
	Conf.Envs = make([]Environment, len(UserConf.Envs))
	copy(Conf.Envs, UserConf.Envs)

	for i := range Conf.Envs {
		Conf.Envs[i].setup()
	}

	for i := range Conf.Envs {
		for j := range Conf.Envs[i].Topics {
			Conf.Envs[i].Topics[j].expandVars(Conf.Envs[i].Vars)
		}
	}
}

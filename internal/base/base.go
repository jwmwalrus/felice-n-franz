package base

import (
	_ "embed" //required to import version.json
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	"github.com/jwmwalrus/bnp/env"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bumpy/pkg/version"
	"github.com/nightlyone/lockfile"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	configFilename = "config.json"
	lockFilename   = configFilename + ".lock"

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
	configFile  string
	lockFile    string
	logFile     *lumberjack.Logger
	logFilename = "app.log"

	flagHelp bool

	// OS operating system's name
	OS string

	//go:embed version.json
	versionJSON []byte

	// AppVersion application's version
	AppVersion version.Version

	// CacheDir home directory for cache
	CacheDir string

	// ConfigDir home directory for config
	ConfigDir string

	// DataDir home directory for data
	DataDir string

	// RuntimeDir run (volatile) directory
	RuntimeDir string

	// FlagVerbose logger severity on
	FlagVerbose bool

	// FlagSeverity logger severity level
	FlagSeverity string

	// FlagEchoLogging echo logs to stderr
	FlagEchoLogging bool

	// Conf application's global configuration
	Conf Config

	// FlagUseConfig user-provided configuration file
	FlagUseConfig string

	userConf Config
)

// Load loads application's configuration
func Load() (args []string) {
	args = env.ParseArgs(logFile, &FlagEchoLogging, &FlagVerbose, &FlagSeverity)

	if flagHelp {
		getopt.Usage()
		os.Exit(0)
	}

	err := AppVersion.Read(versionJSON)
	onerror.Panic(err)

	if FlagUseConfig != "" {
		configFile = FlagUseConfig
	} else {
		configFile = filepath.Join(ConfigDir, configFilename)
	}
	lockFile = filepath.Join(RuntimeDir, lockFilename)

	err = env.SetDirs(
		CacheDir,
		ConfigDir,
		DataDir,
		RuntimeDir,
	)
	onerror.Panic(err)

	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		if FlagUseConfig != "" {
			log.WithFields(log.Fields{
				"--config": FlagUseConfig,
			}).Fatal("No user-provided configuration file was found")
		}
		log.Info(configFilename + " was not found. Generating one")
		userConf.FirstRun = true
		err := Save()
		onerror.Panic(err)
	}

	jsonFile, err := os.Open(configFile)
	onerror.Panic(err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	onerror.Panic(err)

	err = json.Unmarshal(byteValue, &userConf)
	onerror.Panic(err)

	err = userConf.validate()
	onerror.Panic(err)

	configFile = filepath.Join(ConfigDir, configFilename)
	if userConf.FirstRun {
		userConf.FirstRun = false
		err := Save()
		onerror.Panic(err)
	}

	copyUserConfig()

	return
}

// Save saves application's configuration
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
			log.Warningf("Cannot unlock %q, reason: %v", lock, err)
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

// Unload cleans up application before exit
func Unload() {
	log.Info("Unloading application")
	if userConf.FirstRun {
		userConf.FirstRun = false

		// NOTE: ignore error
		Save()
	}
}

func copyUserConfig() {
	Conf.Version = userConf.Version
	Conf.FirstRun = userConf.FirstRun
	Conf.Port = userConf.Port
	Conf.Envs = make([]Environment, len(userConf.Envs))
	copy(Conf.Envs, userConf.Envs)

	for i := range Conf.Envs {
		Conf.Envs[i].setup()
	}

	for i := range Conf.Envs {
		for j := range Conf.Envs[i].Topics {
			Conf.Envs[i].Topics[j].expandVars(Conf.Envs[i].Vars)
		}
	}
}

func init() {
	OS = runtime.GOOS

	// XDG-related
	DataDir = filepath.Join(xdg.DataHome, AppDirName)
	ConfigDir = filepath.Join(xdg.ConfigHome, AppDirName)
	CacheDir = filepath.Join(xdg.CacheHome, AppDirName)
	RuntimeDir = filepath.Join(xdg.RuntimeDir, AppDirName)

	// Define global flags
	getopt.FlagLong(&flagHelp, "help", 'h', "Display this help")
	getopt.FlagLong(&FlagVerbose, "verbose", 'v', "Bump logging severity")
	getopt.FlagLong(&FlagSeverity, "severity", 0, "Logging severity")
	getopt.FlagLong(&FlagEchoLogging, "echo-logging", 'e', "Echo logs to stderr")
	getopt.FlagLong(&FlagUseConfig, "config", 'c', "Use provided config file")

	// log-related
	logFilename = filepath.Base(os.Args[0]) + ".log"
	logFilePath := filepath.Join(DataDir, logFilename)
	logFile = &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     30,    //days
		Compress:   false, // disabled by default
	}
	log.SetOutput(logFile)
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.ErrorLevel)
}

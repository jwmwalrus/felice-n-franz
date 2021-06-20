package base

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	"github.com/jwmwalrus/bnp"
	"github.com/jwmwalrus/bumpy-ride/version"
	"github.com/nightlyone/lockfile"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	configFilename = "config.json"
	lockFilename   = configFilename + ".lock"

	// AppDirName Application's directory name
	AppDirName = "felice-n-franz"

	// AppName Application name
	AppName = "Felice & Franz"

	// DefaultServerAPIVersion idem
	DefaultServerAPIVersion = "/api/v1"

	// DefaultCA Kafka Certificate's filename
	DefaultCA = "kafkacert.pem"

	// DefaultPort Application's default port
	DefaultPort = 9191

	// DefaultTailOffset Maximum number of messages to keep per topic
	DefaultTailOffset = 100
)

var (
	configFile  string
	lockFile    string
	logFile     *lumberjack.Logger
	logFilename = "app.log"

	flagHelp bool

	// OS Operating system's name
	OS string

	// AppVersion application's version
	AppVersion version.Version

	// CacheDir Home directory for cache
	CacheDir string

	// ConfigDir Home directory for config
	ConfigDir string

	// DataDir Home directory for data
	DataDir string

	// RuntimeDir Run (volatile) directory
	RuntimeDir string

	// FlagVerbose Logger severity on
	FlagVerbose bool

	// FlagSeverity Logger severity level
	FlagSeverity string

	// FlagEchoLogging Echo logs to stderr
	FlagEchoLogging bool
)

// Conf Application's global configuration
var Conf, userConf Config

// Load Loads application's configuration
func Load(version []byte) (args []string) {
	args = bnp.ParseArgs(logFile, &FlagEchoLogging, &FlagVerbose, &FlagSeverity)

	if flagHelp {
		getopt.Usage()
		os.Exit(0)
	}

	err := AppVersion.Read(version)
	bnp.PanicOnError(err)

	configFile = filepath.Join(ConfigDir, configFilename)
	lockFile = filepath.Join(RuntimeDir, lockFilename)

	err = bnp.SetEnvDirs(
		configFile,
		lockFile,
		CacheDir,
		ConfigDir,
		DataDir,
		RuntimeDir,
	)
	bnp.PanicOnError(err)

	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		log.Info(configFilename + " was not found. Generating one")
		userConf.FirstRun = true
		err := Save()
		bnp.PanicOnError(err)
	}

	// var jsonFile *os.File
	jsonFile, err := os.Open(configFile)
	bnp.PanicOnError(err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	bnp.PanicOnError(err)

	err = json.Unmarshal(byteValue, &userConf)
	bnp.PanicOnError(err)

	err = userConf.validate()
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

// Unload Cleans up application before exit
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

package base

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	configFilename = "config.json"
	lockFilename   = configFilename + ".lock"

	// AppDirName Application's directory name
	AppDirName = "felice+franz"

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

	// OS Operating system's name
	OS string

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

	// Log Application's log
	Log = log.New()
)

// ParseArgs parses the given command line arguments
func ParseArgs() (args []string) {
	getopt.Parse()
	args = getopt.Args()
	arg0 := []string{os.Args[0]}
	args = append(arg0, args...)

	resolveSeverity()

	if FlagEchoLogging {
		mw := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(mw)
	}

	return
}

// SetEnvDirs Ensure that environment directories exist
func SetEnvDirs() (err error) {
	configFile = filepath.Join(ConfigDir, configFilename)
	lockFile = filepath.Join(RuntimeDir, lockFilename)

	if _, err = os.Stat(CacheDir); os.IsNotExist(err) {
		err = os.MkdirAll(CacheDir, 0755)
		if err != nil {
			return
		}
	}

	if _, err = os.Stat(ConfigDir); os.IsNotExist(err) {
		err = os.MkdirAll(ConfigDir, 0755)
		if err != nil {
			return
		}
	}

	if _, err = os.Stat(DataDir); os.IsNotExist(err) {
		err = os.MkdirAll(DataDir, 0755)
		if err != nil {
			return
		}
	}

	if _, err = os.Stat(RuntimeDir); os.IsNotExist(err) {
		err = os.MkdirAll(RuntimeDir, 0755)
	}

	return
}

func resolveSeverity() {
	givenSeverity := FlagSeverity

	if givenSeverity == "" {
		if FlagVerbose {
			FlagSeverity = "info"
		} else {
			FlagSeverity = "error"
		}
	} else {
		if _, err := log.ParseLevel(givenSeverity); err != nil {
			FlagSeverity = "error"
		} else {
			FlagSeverity = givenSeverity
		}
	}

	level, _ := log.ParseLevel(FlagSeverity)
	log.SetLevel(level)
	log.SetReportCaller(FlagSeverity == "debug")

	return
}

func init() {
	OS = runtime.GOOS

	// XDG-related
	DataDir = filepath.Join(xdg.DataHome, AppDirName)
	ConfigDir = filepath.Join(xdg.ConfigHome, AppDirName)
	CacheDir = filepath.Join(xdg.CacheHome, AppDirName)
	RuntimeDir = filepath.Join(xdg.RuntimeDir, AppDirName)

	// Define global flags
	getopt.FlagLong(&FlagVerbose, "verbose", 'v', "Bump logging severity")
	getopt.FlagLong(&FlagSeverity, "severity", 0, "Logging severity")
	getopt.FlagLong(&FlagEchoLogging, "echo-logging", 'e', "Echo logs to stderr")

	// Log-related
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

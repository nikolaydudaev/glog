// config package sets up scaffolding for the applocation configuration
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// [APP] variable defines the name of the application
var APP = ".glog"

// [HOME] gets the user home directory
var HOME, _ = os.UserHomeDir()

// [DATA] defines application data directory
var DATA = "data"
var STAGING = "staged"
var PROCESSED = "processed"
var DB = "db"

// [CONFIG] defines application config directory
var CONFIG = "etc"

// [LOGGIN] defines application internal logs directory
var LOGGING = "log"

// [CONFIG_FILE] defines the configuration base name
var CONFIG_FILE = "glog.toml"

type Options struct {
	Clean bool
}

type LogTypeMeta struct {
	LogName            string `toml:"logName"`
	LogSourceDir       string `toml:"logSourceDir"`
	LogFilePattern     string `toml:"logFilePattern"`
	LogFileCompression string `toml:"logFileCompression"`
}

// [AppConfig] keeps all application configurations in the struct
type AppConfig struct {
	AppHome          string        `toml:"AppHome"`
	AppConfig        string        `toml:"AppConfig"`
	AppDataStaging   string        `toml:"AppDataStaging"`
	AppDataProcessed string        `toml:"AppDataProcessed"`
	AppDataDB        string        `toml:"AppDataDB"`
	AppLog           string        `toml:"AppLog"`
	DefinedLogs      []LogTypeMeta `toml:"DefinedLogs"`
}

var Placeholder = LogTypeMeta{
	LogName:            "nginx",
	LogSourceDir:       "nginx_logs",
	LogFilePattern:     "access",
	LogFileCompression: ".gz",
}

var TestConfig = AppConfig{
	AppHome:          "app home",
	AppConfig:        "app home",
	AppDataStaging:   "app home",
	AppDataProcessed: "app home",
	AppDataDB:        "app home",
	AppLog:           "app home",
	DefinedLogs: []LogTypeMeta{
		Placeholder,
	},
}

func parseSourceDir(fpath string, home string) (string, error) {
	var err error
	err = errors.New("`sourceDir` is empty")
	if fpath == "" {
		return "", fmt.Errorf("%w\n", err)
	}

	replacer := strings.NewReplacer("~/$HOME", os.ExpandEnv(home), "~", os.ExpandEnv(home), "$HOME", os.ExpandEnv(home))
	fpath = replacer.Replace(fpath)
	fpath = filepath.Clean(fpath)
	err = errors.New("`sourceDir` is not an absolute path")
	if !filepath.IsAbs(fpath) {
		return "", fmt.Errorf("%w\n", err)
	}
	_, err = os.Stat(fpath)
	if errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("The `sourceDir` you provided in glog.toml does not exist;\n%w\n", err)
	}

	return fpath, err
}

func CleanAppDirectory() error {
	err := os.RemoveAll(filepath.Join(HOME, APP))
	if errors.Is(err, &fs.PathError{}) {
		return fmt.Errorf("Nothing to clean %w", err)
	}
	return err
}

func ReadConfig(HOME string, APP string, CONFIG string, CONFIG_FILE string, appConfig *AppConfig) error {
	var err error
	config_file, err := os.Open(filepath.Join(HOME, APP, CONFIG, CONFIG_FILE))
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("Config file does not exist. Have you ran 'glog init'?")
	}
	defer config_file.Close()

	decoder := toml.NewDecoder(config_file)
	meta, err := decoder.Decode(appConfig)
	if err != nil {
		log.Fatal(err)
	}
	ErrIllegalConfigFileds := errors.New("There are illegal fields in configuration file. Aborting!")
	if len(meta.Undecoded()) > 0 {
		return fmt.Errorf("%w", ErrIllegalConfigFileds)
	}
	for _, log := range appConfig.DefinedLogs {
		log.LogSourceDir, err = parseSourceDir(log.LogSourceDir, HOME)
	}
	return err
}

func Init(HOME string, APP string, CONFIG string, DATA string, STAGING string, PROCESSED string, DB string, LOGGING string, CONFIG_FILE string, placeholders []LogTypeMeta) error {
	var err error
	var config AppConfig
	config.AppHome = filepath.Join(HOME, APP)
	config.AppConfig = filepath.Join(HOME, APP, CONFIG)
	config.AppDataStaging = filepath.Join(HOME, APP, DATA, STAGING)
	config.AppDataProcessed = filepath.Join(HOME, APP, DATA, PROCESSED)
	config.AppDataDB = filepath.Join(HOME, APP, DATA, DB)
	config.AppLog = filepath.Join(HOME, APP, LOGGING)

	if err = os.MkdirAll(config.AppHome, 0o755); err != nil {
		return fmt.Errorf("Was not able to create %s due to %w\n", config.AppHome, err)
	}
	if err = os.MkdirAll(config.AppConfig, 0o755); err != nil {
		return fmt.Errorf("Was not able to create %s due to %w\n", config.AppConfig, err)
	}
	if err = os.MkdirAll(config.AppDataStaging, 0o755); err != nil {
		return fmt.Errorf("Was not able to create %s due to %w\n", config.AppDataStaging, err)
	}
	if err = os.MkdirAll(config.AppDataProcessed, 0o755); err != nil {
		return fmt.Errorf("Was not able to create %s due to %w\n", config.AppDataProcessed, err)
	}
	if err = os.MkdirAll(config.AppDataDB, 0o755); err != nil {
		return fmt.Errorf("Was not able to create %s due to %w\n", config.AppDataDB, err)
	}
	if err = os.MkdirAll(config.AppLog, 0o755); err != nil {
		return fmt.Errorf("Was not able to create %s due to %w\n", config.AppLog, err)
	}
	config_file, err := os.OpenFile(filepath.Join(HOME, APP, CONFIG, CONFIG_FILE), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("Was not able to open config file due to %w\n", err)
	}
	defer config_file.Close()

	encoder := toml.NewEncoder(config_file)

	config.DefinedLogs = placeholders

	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("Was not able to write to config file due to %w\n", err)
	}
	config_file.Sync()

	return err
}

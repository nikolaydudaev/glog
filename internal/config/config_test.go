package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestParseSourceDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	os.MkdirAll(filepath.Join(dir, "my-raw-logs"), 0o755)
	srcDirPaths := map[string]struct {
		in      string
		out     string
		wantErr bool
	}{
		"tilde": {
			in:      "~/my-raw-logs",
			out:     filepath.Join(dir, "my-raw-logs"),
			wantErr: false,
		},
		"$HOME": {
			in:      "$HOME/my-raw-logs",
			out:     filepath.Join(dir, "my-raw-logs"),
			wantErr: false,
		},
		"Not Abs Path": {
			in:      "./my-raw-logs",
			out:     "",
			wantErr: true,
		},
		"tilde+$HOME+doubleslash": {
			in:      "~/$HOME//my-raw-logs",
			out:     filepath.Join(dir, "my-raw-logs"),
			wantErr: false,
		},
		"empty string": {
			in:      "",
			out:     "",
			wantErr: true,
		},
	}

	for name, tt := range srcDirPaths {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s, err := parseSourceDir(tt.in, dir)
			if tt.wantErr {
				if err != nil {
					fmt.Println(err)
				} else {
					if s != tt.out {
						t.Errorf("got %q (err: %v), want %q", s, err, tt.out)
					}
				}
			} else {
				if s != tt.out {
					t.Errorf("got %q (err: %v), want %q", s, err, tt.out)
				}
			}

		})
	}
}

func TestInit(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	var APP = ".glog"
	var DATA = "data"
	var STAGING = "staged"
	var PROCESSED = "processed"
	var DB = "db"
	var CONFIG = "etc"
	var LOGGING = "log"
	var CONFIG_FILE = "glog.toml"
	var config AppConfig
	config.AppHome = filepath.Join(dir, APP)
	config.AppConfig = filepath.Join(dir, APP, CONFIG)
	config.AppDataStaging = filepath.Join(dir, APP, DATA, STAGING)
	config.AppDataProcessed = filepath.Join(dir, APP, DATA, PROCESSED)
	config.AppDataDB = filepath.Join(dir, APP, DATA, DB)
	config.AppLog = filepath.Join(dir, APP, LOGGING)

	err := Init(
		dir,
		APP,
		CONFIG,
		DATA,
		STAGING,
		PROCESSED,
		DB,
		LOGGING,
		CONFIG_FILE,
		[]LogTypeMeta{Placeholder, Placeholder},
	)
	if err != nil {
		t.Errorf("Init error: %v\n", err)
	}
}

func TestConfigFile(t *testing.T) {

	dir := t.TempDir()
	t.Setenv("HOME", dir)
	var APP = ".glog"
	var DATA = "data"
	var STAGING = "staged"
	var PROCESSED = "processed"
	var DB = "db"
	var CONFIG = "etc"
	var LOGGING = "log"
	var CONFIG_FILE = "glog.toml"
	var config AppConfig

	content := fmt.Sprintf(`AppHome = "%s/.glog"
	AppConfig = "%s/.glog/etc"
	AppDataStaging = "%s/.glog/data/staged"
	AppDataProcessed = "%s/.glog/data/processed"
	AppDataDB = "%s/.glog/data/db"
	AppLog = "%s/.glog/log"

	[[DefinedLogs]]
  logName = "nginx"
  logSourceDir = "%s/nginx_logs"
  logFilePattern = "access"
  logFileCompression = ".gz"

	[[DefinedLogs]]
  logName = "nginx"
  logSourceDir = "%s/nginx_logs"
  logFilePattern = "access"
  logFileCompression = ".gz"`, dir, dir, dir, dir, dir, dir, dir, dir)

	if err := os.MkdirAll(filepath.Join(dir, "/nginx_logs"), 0o755); err != nil {
		t.Errorf("Was not able to create due to %v\n", err)
	}
	fmt.Println(filepath.Join(dir, "/nginx_logs"))

	err := Init(
		dir,
		APP,
		CONFIG,
		DATA,
		STAGING,
		PROCESSED,
		DB,
		LOGGING,
		CONFIG_FILE,
		[]LogTypeMeta{Placeholder, Placeholder},
	)
	if err != nil {
		t.Errorf("Init error: %v\n", err)
	}

	config_path := filepath.Join(dir, APP, CONFIG, CONFIG_FILE)

	if err = os.WriteFile(config_path, []byte(content), 0o600); err != nil {
		t.Errorf("Writing failed: %v", err)
	}
	err = ReadConfig(
		dir,
		APP,
		CONFIG,
		CONFIG_FILE,
		&config,
	)
	if err != nil {
		t.Errorf("Reading failed: %v", err)
	}
	if config.AppHome != filepath.Join(dir, APP) {
		t.Errorf("AppHome error: %v", err)
	}
	if config.AppConfig != filepath.Join(dir, APP, CONFIG) {
		t.Errorf("AppConfig error: %v", err)
	}
	if config.AppDataStaging != filepath.Join(dir, APP, DATA, STAGING) {
		t.Errorf("AppDataStaging error: %v", err)
	}
	if config.AppDataProcessed != filepath.Join(dir, APP, DATA, PROCESSED) {
		t.Errorf("AppDataProcessed error: %v", err)
	}
	if config.AppDataDB != filepath.Join(dir, APP, DATA, DB) {
		t.Errorf("AppDataDB error: %v", err)
	}
	if config.AppLog != filepath.Join(dir, APP, LOGGING) {
		t.Errorf("AppLog error: %v", err)
	}
	if config.DefinedLogs[0].LogName != "nginx" {
		t.Errorf("DefinedLogs error: %v", err)
	}
	if config.DefinedLogs[0].LogSourceDir != filepath.Join(dir, "/nginx_logs") {
		t.Errorf("DefinedLogs error: %v", err)
	}
	if config.DefinedLogs[0].LogFilePattern != "access" {
		t.Errorf("DefinedLogs error: %v", err)
	}
	if config.DefinedLogs[0].LogFileCompression != ".gz" {
		t.Errorf("DefinedLogs error: %v", err)
	}

}

package extract

import (
	"path/filepath"
	"testing"

	"git.nikolaydudaev.com/ndudaev/glog/internal/config"
)

func TestNewExtract(t *testing.T) {
	var cfg config.AppConfig
	var logType config.LogTypeMeta

	cfg.AppHome = "app_home"
	cfg.AppConfig = "app_config"
	cfg.AppDataStaging = "app_data/staging"
	cfg.AppDataProcessed = "app_data/processed"
	cfg.AppDataDB = "app_data/db"
	cfg.AppLog = "app_log"

	logType.LogName = "test log"
	logType.LogSourceDir = "/log_source_dir"
	logType.LogFilePattern = "access"
	logType.LogFileCompression = ".zip"

	logExtr, err := NewExtractor(&cfg, &logType)
	if err != nil {
		t.Errorf("NewExtractor error: %v", err)
	}
	if logExtr.SourceDir != "/log_source_dir" {
		t.Errorf("Extractor error: got %q, want %q", logExtr.SourceDir, "/log_source_dir")
	}
	if logExtr.StageFile != filepath.Join(cfg.AppDataStaging, logType.LogName+"_staged.log") {
		t.Errorf("Extractor error: got %q, want %q", logExtr.StageFile, filepath.Join(cfg.AppDataStaging, logType.LogName+"_staged.log"))
	}
	if logExtr.LogName != "test log" {
		t.Errorf("Extractor error: got %q, want %q", logExtr.LogName, "test log")
	}
	if logExtr.LogFilePattern != "access" {
		t.Errorf("Extractor error: got %q, want %q", logExtr.LogFilePattern, "access")
	}
	if logExtr.LogFileCompression != ".zip" {
		t.Errorf("Extractor error: got %q, want %q", logExtr.LogFileCompression, ".zip")
	}
}

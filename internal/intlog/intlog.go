package intlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ExtractLogs struct {
	LogType              string
	LogDate              time.Time
	FileName             string
	SourceFilesTarget    int
	SourceFilesProcessed int
	RowsProcessed        int
}

func WriteExtractionLogs(intlogPath string, extrLogs *ExtractLogs) (string, error) {
	var err error
	extractionLogsPath := filepath.Join(intlogPath, "extraction-logs.log")
	f, err := os.OpenFile(extractionLogsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("Couldn't open `extraction-logs.log` for writing: %w", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	if err = encoder.Encode(extrLogs); err != nil {
		return "", fmt.Errorf("Error when writing JSON to `extraction-log.log`: %w", err)
	}

	f.Sync()

	return extractionLogsPath, err
}

func ReadExtractionLogs(extractionLogsPath string) (ExtractLogs, error) {
	var extractionLogs ExtractLogs
	var err error
	f, err := os.Open(extractionLogsPath)
	if err != nil {
		return ExtractLogs{}, fmt.Errorf("Could not open a extraction logs files due to %w", err)
	}

	decoder := json.NewDecoder(f)
	if err = decoder.Decode(&extractionLogs); err != nil {
		return ExtractLogs{}, fmt.Errorf("Could not decode extraction logs due to %w", err)
	}

	return extractionLogs, err
}

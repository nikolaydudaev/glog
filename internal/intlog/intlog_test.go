package intlog

import (
	"testing"
	"time"
)

func TestWriteExtractionLogs(t *testing.T) {
	var extracted ExtractLogs
	var err error

	dir := t.TempDir()

	extracted.LogType = "test"
	extracted.LogDate = time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC)
	extracted.FileName = "test_extracted"
	extracted.SourceFilesTarget = 15
	extracted.SourceFilesProcessed = 15
	extracted.RowsProcessed = 1000

	_, err = WriteExtractionLogs(dir, &extracted)
	if err != nil {
		t.Errorf("WriteExtractionLogs error: %v", err)
	}
}

func TestReadExtractionLogs(t *testing.T) {
	var written ExtractLogs
	var read ExtractLogs
	var err error

	dir := t.TempDir()

	written.LogType = "test"
	written.LogDate = time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC)
	written.FileName = "test_extracted"
	written.SourceFilesTarget = 15
	written.SourceFilesProcessed = 15
	written.RowsProcessed = 1000

	f, err := WriteExtractionLogs(dir, &written)
	if err != nil {
		t.Errorf("WriteExtractionLogs error: %v", err)
	}

	read, err = ReadExtractionLogs(f)
	if err != nil {
		t.Errorf("ReadExtractionLogs error: %v", err)
	}

	if read.LogType != "test" {
		t.Errorf("Read: got %q, want %q", read.LogType, "test")
	}
	if read.LogDate != time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC) {
		t.Errorf("Read: got %q, want %q", read.LogDate, time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC))
	}
	if read.FileName != "test_extracted" {
		t.Errorf("Read: got %q, want %q", read.FileName, "test_extracted")
	}
	if read.SourceFilesTarget != 15 {
		t.Errorf("Read: got %d, want %d", read.SourceFilesTarget, 15)
	}
	if read.SourceFilesProcessed != 15 {
		t.Errorf("Read: got %d, want %d", read.SourceFilesProcessed, 15)
	}
	if read.RowsProcessed != 1000 {
		t.Errorf("Read: got %d, want %d", read.RowsProcessed, 15)
	}
}

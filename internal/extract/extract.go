package extract

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.nikolaydudaev.com/ndudaev/glog/internal/config"
	"git.nikolaydudaev.com/ndudaev/glog/internal/intlog"
)

// Extractor is an interface to cover data extraction
// from source log files into one stage file
type Extractor interface {
	ReadDir() (int, []os.DirEntry, error)      // Return number of relevant files to extract and directory list
	Extract() (int, error)                     // Return number of files processed
	Stat() (int, error)                        // Return File Stat + number of lines in a file
	RunExtractor() (intlog.ExtractLogs, error) // Return ExtractLogs object
}

// BaseExtractor object
// Identifies key properties of nginx logs Extractor
type BaseExtractor struct {
	SourceDir          string
	StageFile          string
	LogName            string
	LogFilePattern     string
	LogFileCompression string
}

func NewExtractor(cfg *config.AppConfig, defined *config.LogTypeMeta) (BaseExtractor, error) {
	var err error
	extractor := BaseExtractor{}
	extractor.SourceDir = defined.LogSourceDir
	extractor.StageFile = filepath.Join(cfg.AppDataStaging, defined.LogName+"_staged.log")
	extractor.LogName = defined.LogName
	extractor.LogFilePattern = defined.LogFilePattern
	extractor.LogFileCompression = defined.LogFileCompression

	defined.LogStageFile = extractor.StageFile

	if extractor.LogName == "" || extractor.SourceDir == "" || extractor.StageFile == "" || extractor.LogFilePattern == "" || extractor.LogFileCompression == "" {
		return BaseExtractor{}, fmt.Errorf("NewExtractor failed with %s: error %w", defined.LogName, err)
	}

	return extractor, err
}

func (e *BaseExtractor) ReadDir() (int, []os.DirEntry) {
	dir, err := os.ReadDir(e.SourceDir)
	if err != nil {
		log.Fatal(err)
	}
	files_count := 0
	for _, file := range dir {
		if strings.Contains(file.Name(), e.LogFilePattern) || (strings.Contains(file.Name(), e.LogFilePattern) && strings.Contains(file.Name(), e.LogFileCompression)) {
			files_count++
		}
	}

	return files_count, dir
}

func (e *BaseExtractor) Extract() int {
	_, dir := e.ReadDir()

	foutput, err := os.OpenFile(e.StageFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer foutput.Close()

	writer := bufio.NewWriter(foutput)

	files_count := 0
	for _, file := range dir {
		if strings.Contains(file.Name(), e.LogFilePattern) && strings.Contains(file.Name(), e.LogFileCompression) {
			fungzip := filepath.Join(e.SourceDir, file.Name())
			gzf, err := os.Open(fungzip)
			if err != nil {
				log.Fatal(err)
			}
			defer gzf.Close()

			gzreader, err := gzip.NewReader(gzf)
			if err != nil {
				log.Fatal(err)
			}
			defer gzreader.Close()

			if _, err := io.Copy(writer, gzreader); err != nil {
				log.Fatal(err)
			}
		} else if strings.Contains(file.Name(), e.LogFilePattern) {
			fregular := filepath.Join(e.SourceDir, file.Name())
			regf, err := os.Open(fregular)
			if err != nil {
				log.Fatal(err)
			}
			defer regf.Close()

			regreader := bufio.NewReader(regf)

			if _, err := io.Copy(writer, regreader); err != nil {
				log.Fatal(err)
			}
		}
		files_count++
	}
	return files_count
}

func (e *BaseExtractor) Stat() int {
	file, err := os.Open(e.StageFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	num_lines := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		num_lines++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading file:", err)
	}

	return num_lines
}

func (e *BaseExtractor) RunExtractor() intlog.ExtractLogs {
	internal_logs := intlog.ExtractLogs{}
	if e.SourceDir != "" {
		relevant_files, _ := e.ReadDir()
		processed := e.Extract()
		lines := e.Stat()

		internal_logs = intlog.ExtractLogs{
			LogType:              e.LogName,
			LogDate:              time.Now(),
			FileName:             e.StageFile,
			SourceFilesTarget:    relevant_files,
			SourceFilesProcessed: processed,
			RowsProcessed:        lines,
		}
	} else {
		fmt.Printf("Skipping %s, no source directory provided...\n", e.LogName)
	}

	return internal_logs
}

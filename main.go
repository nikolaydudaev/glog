package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"git.nikolaydudaev.com/ndudaev/glog/internal/config"
	"git.nikolaydudaev.com/ndudaev/glog/internal/extract"
	"git.nikolaydudaev.com/ndudaev/glog/internal/intlog"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

func initRun(args []string) error {
	var initOpts config.Options
	var err error
	init := flag.NewFlagSet("init", flag.ContinueOnError)
	init.BoolVar(&initOpts.Clean, "clean", false, "Removes app directory and all its content. Run `glog init` to initiate glog again.")
	init.Parse(args)

	if os.Args[1] == "init" && init.NFlag() == 0 {
		err = config.Init(
			config.HOME,
			config.APP,
			config.CONFIG,
			config.DATA,
			config.STAGING,
			config.PROCESSED,
			config.DB,
			config.LOGGING,
			config.CONFIG_FILE,
			[]config.LogTypeMeta{config.Placeholder, config.Placeholder},
		)
		if err != nil {
			return fmt.Errorf("Could not initiate app directories; %w\n", err)
		}
	}
	if initOpts.Clean {
		err = config.CleanAppDirectory()
		if err != nil {
			return fmt.Errorf("Could not clean up app directories; %w\n", err)
		}
	}
	if (init.NFlag() > 1 && initOpts.Clean == false) || len(args) > 1 {
		return fmt.Errorf("To many flags or flags are used incorrectly. Try `glog init -h` for help\n")
	}

	return err
}

func extractRun(args []string) error {
	var err error
	extractCmd := flag.NewFlagSet("extract", flag.ContinueOnError)
	extractCmd.Parse(args)

	if os.Args[1] == "extract" && extractCmd.NFlag() == 0 {
		// Read config file
		var cfg config.AppConfig
		err = config.ReadConfig(
			config.HOME,
			config.APP,
			config.CONFIG,
			config.CONFIG_FILE,
			&cfg,
		)
		if err != nil {
			log.Fatal(err)
		}

		bar := progressbar.NewOptions(len(cfg.DefinedLogs),
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()), //you should install "github.com/k0kubun/go-ansi"
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(100),
			progressbar.OptionSetDescription("[purple][1/2][reset] Extracting raw log files..."),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[purple]=[reset]",
				SaucerHead:    "[purple]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}))

		for _, l := range cfg.DefinedLogs {
			bar.Add(1)
			logExtr, err := extract.NewExtractor(&cfg, l)
			if err != nil {
				log.Fatal(err)
			}
			intExtractionLogs := logExtr.RunExtractor()
			_, err = intlog.WriteExtractionLogs(cfg.AppLog, &intExtractionLogs)
			time.Sleep(5 * time.Millisecond)
		}
		// Loop over DefinedLogs, create Extractor object, run extractor and write to extraction logs
	}
	if extractCmd.NFlag() > 1 || len(args) > 1 {
		return fmt.Errorf("To many flags or flags are used incorrectly. Try `glog init -h` for help\n")
	}

	return err
}

func main() {
	// Initialize
	if len(os.Args) < 2 {
		fmt.Println("expected 'init', 'extract' or 'export' subcommands")
	}

	switch os.Args[1] {
	case "init":
		err := initRun(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "extract":
		err := extractRun(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("expected 'init', 'extract' or 'export' subcommands")
	}
}

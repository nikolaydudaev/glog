package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"git.nikolaydudaev.com/ndudaev/glog/internal/config"
	"git.nikolaydudaev.com/ndudaev/glog/internal/export"
	"git.nikolaydudaev.com/ndudaev/glog/internal/extract"
	"git.nikolaydudaev.com/ndudaev/glog/internal/intlog"
	"git.nikolaydudaev.com/ndudaev/glog/internal/parse"
	"git.nikolaydudaev.com/ndudaev/glog/internal/utils"
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

func extractRun(args []string, cfg *config.AppConfig) error {
	var err error
	extractCmd := flag.NewFlagSet("extract", flag.ContinueOnError)
	extractCmd.Parse(args)

	if os.Args[1] == "extract" && extractCmd.NFlag() == 0 {
		bar := utils.CreateProgressBar(len(cfg.DefinedLogs), "1", "1", "Extracting raw log files ...")

		for _, l := range cfg.DefinedLogs {
			bar.Add(1)
			logExtr, err := extract.NewExtractor(cfg, &l)
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

func exportRun(args []string, cfg *config.AppConfig) error {
	var err error
	exportCmd := flag.NewFlagSet("export", flag.ContinueOnError)
	exportCmd.Parse(args)

	if os.Args[1] == "export" && exportCmd.NFlag() == 0 {
		serverRecords := make(chan parse.ServerEvent, 1000)
		firewallRecords := make(chan parse.FirewallEvent, 1000)
		authRecords := make(chan parse.AuthEvent, 1000)
		fail2banRecords := make(chan parse.FailbanEvent, 1000)

		var reader sync.WaitGroup
		var writer sync.WaitGroup

		bar := utils.CreateProgressBar(len(cfg.DefinedLogs), "1", "2", "Reading events stage files...")
		for i := range cfg.DefinedLogs {
			switch cfg.DefinedLogs[i].LogName {
			case "nginx", "caddy":
				reader.Add(1)
				go func() {
					defer reader.Done()
					export.ReadServerEvents(cfg.DefinedLogs[i].LogStageFile, cfg.DefinedLogs[i].LogName, serverRecords)
					bar.Add(1)
				}()
			case "ufw":
				reader.Add(1)
				go func() {
					defer reader.Done()
					export.ReadFirewallEvents(cfg.DefinedLogs[i].LogStageFile, firewallRecords)
					bar.Add(1)
				}()
			case "fail2ban":
				reader.Add(1)
				go func() {
					defer reader.Done()
					export.ReadFailbanEvents(cfg.DefinedLogs[i].LogStageFile, fail2banRecords)
					bar.Add(1)
				}()
			case "auth":
				reader.Add(1)
				go func() {
					defer reader.Done()
					export.ReadAuthEvents(cfg.DefinedLogs[i].LogStageFile, authRecords)
					bar.Add(1)
				}()
			}
			time.Sleep(5 * time.Millisecond)
		}

		go func() {
			reader.Wait()
			close(serverRecords)
			close(firewallRecords)
			close(fail2banRecords)
			close(authRecords)
		}()

		writer.Add(1)
		go func() {
			defer writer.Done()
			export.WriteServerEvents(filepath.Join(cfg.AppDataProcessed, "server-events.json"), serverRecords)
			bar.Add(1)
		}()
		writer.Add(1)
		go func() {
			defer writer.Done()
			export.WriteFirewallEvents(filepath.Join(cfg.AppDataProcessed, "firewall-events.json"), firewallRecords)
			bar.Add(1)
		}()
		writer.Add(1)
		go func() {
			defer writer.Done()
			export.WriteFailbanEvents(filepath.Join(cfg.AppDataProcessed, "fail2ban-events.json"), fail2banRecords)
			bar.Add(1)
		}()
		writer.Add(1)
		go func() {
			defer writer.Done()
			export.WriteAuthEvents(filepath.Join(cfg.AppDataProcessed, "auth-events.json"), authRecords)
			bar.Add(1)
		}()

		writer.Wait()

		err := export.WriteToDuckDB(cfg)
		if err != nil {
			log.Fatal(err)
		}

	}
	if exportCmd.NFlag() > 1 || len(args) > 1 {
		return fmt.Errorf("To many flags or flags are used incorrectly. Try `glog init -h` for help\n")
	}

	return err
}

func main() {
	var err error
	// Initialize
	if len(os.Args) < 2 {
		fmt.Println("expected 'init', 'extract' or 'export' subcommands")
	}

	// Read config file
	var cfg config.AppConfig

	switch os.Args[1] {
	case "init":
		err := initRun(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "extract":
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

		for i := range cfg.DefinedLogs {
			_, err := extract.NewExtractor(&cfg, &cfg.DefinedLogs[i])
			if err != nil {
				log.Fatal(err)
			}
		}
		err := extractRun(os.Args[2:], &cfg)
		if err != nil {
			log.Fatal(err)
		}
	case "export":
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

		for i := range cfg.DefinedLogs {
			_, err := extract.NewExtractor(&cfg, &cfg.DefinedLogs[i])
			if err != nil {
				log.Fatal(err)
			}
		}
		err := exportRun(os.Args[2:], &cfg)
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("expected 'init', 'extract' or 'export' subcommands")
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"git.nikolaydudaev.com/ndudaev/glog/internal/config"
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
	default:
		fmt.Println("expected 'init', 'extract' or 'export' subcommands")
	}
}

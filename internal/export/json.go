package export

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/schollz/progressbar/v3"

	"git.nikolaydudaev.com/ndudaev/glog/internal/parse"
)

var writeFlags int = os.O_TRUNC | os.O_CREATE | os.O_WRONLY

func ReadServerEvents(fpath string, server string, records chan<- parse.ServerEvent) {
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	bar := progressbar.Default(int64(100), "Reading server events stage files...")
	for scanner.Scan() {
		switch server {
		case "nginx":
			parse.NginxParserInstance.Parse(scanner.Text(), records)
		case "caddy":
			parse.CaddyParserInstance.Parse(scanner.Text(), records)
		}
		bar.Add(1)
	}
	if scanerr := scanner.Err(); scanerr != nil {
		log.Fatal(err)
	}
}

func ReadFirewallEvents(fpath string, records chan<- parse.FirewallEvent) {
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	bar := progressbar.Default(int64(100), "Reading ufw stage files...")
	for scanner.Scan() {
		parse.UfwParserInstance.FirewallParse(scanner.Text(), records)
		bar.Add(1)
	}
	if scanerr := scanner.Err(); scanerr != nil {
		log.Fatal(err)
	}
}

func ReadAuthEvents(fpath string, records chan<- parse.AuthEvent) {
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	bar := progressbar.Default(int64(100), "Reading auth stage files...")
	for scanner.Scan() {
		parse.AuthParserInstance.AuthParse(scanner.Text(), records)
		bar.Add(1)
	}
	if scanerr := scanner.Err(); scanerr != nil {
		log.Fatal(err)
	}
}

func ReadFailbanEvents(fpath string, records chan<- parse.FailbanEvent) {
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	bar := progressbar.Default(int64(100), "Reading fail2ban stage files...")
	for scanner.Scan() {
		parse.FailbanParserInstance.FailbanParse(scanner.Text(), records)
		bar.Add(1)
	}
	if scanerr := scanner.Err(); scanerr != nil {
		log.Fatal(err)
	}
}

func WriteServerEvents(fpath string, records <-chan parse.ServerEvent) {
	f, err := os.OpenFile(fpath, writeFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	bar := progressbar.Default(int64(len(records)), "Writing to server events JSON...")
	for sr := range records {
		if err := enc.Encode(sr); err != nil {
			log.Fatal(err)
		}
		bar.Add(1)
	}
}

func WriteFirewallEvents(fpath string, records chan parse.FirewallEvent) {
	f, err := os.OpenFile(fpath, writeFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	bar := progressbar.Default(int64(len(records)), "Writing to firwall events JSON...")
	for sr := range records {
		if err := enc.Encode(sr); err != nil {
			log.Fatal(err)
		}
		bar.Add(1)
	}
}

func WriteAuthEvents(fpath string, records chan parse.AuthEvent) {
	f, err := os.OpenFile(fpath, writeFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	bar := progressbar.Default(int64(len(records)), "Writing to auth events JSON...")
	for sr := range records {
		if err := enc.Encode(sr); err != nil {
			log.Fatal(err)
		}
		bar.Add(1)
	}
}

func WriteFailbanEvents(fpath string, records chan parse.FailbanEvent) {
	f, err := os.OpenFile(fpath, writeFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	bar := progressbar.Default(int64(len(records)), "Writing to fail2ban events JSON...")
	for sr := range records {
		if err := enc.Encode(sr); err != nil {
			log.Fatal(err)
		}
		bar.Add(1)
	}
}

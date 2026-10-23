package export

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

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
	for scanner.Scan() {
		switch server {
		case "nginx":
			parse.NginxParserInstance.Parse(scanner.Text(), records)
		case "caddy":
			parse.CaddyParserInstance.Parse(scanner.Text(), records)
		}
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
	for scanner.Scan() {
		parse.UfwParserInstance.FirewallParse(scanner.Text(), records)
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
	for scanner.Scan() {
		parse.AuthParserInstance.AuthParse(scanner.Text(), records)
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
	for scanner.Scan() {
		parse.FailbanParserInstance.FailbanParse(scanner.Text(), records)
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
	for sr := range records {
		if err := enc.Encode(sr); err != nil {
			log.Fatal(err)
		}
	}
}

func WriteFirewallEvents(fpath string, records chan parse.FirewallEvent) {
	f, err := os.OpenFile(fpath, writeFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for sr := range records {
		if err := enc.Encode(sr); err != nil {
			log.Fatal(err)
		}
	}
}

func WriteAuthEvents(fpath string, records chan parse.AuthEvent) {
	f, err := os.OpenFile(fpath, writeFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for sr := range records {
		if err := enc.Encode(sr); err != nil {
			log.Fatal(err)
		}
	}
}

func WriteFailbanEvents(fpath string, records chan parse.FailbanEvent) {
	f, err := os.OpenFile(fpath, writeFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for sr := range records {
		if err := enc.Encode(sr); err != nil {
			log.Fatal(err)
		}
	}
}

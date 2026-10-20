# Extract, Parse, Transform and Load Your logs

`glog` is a simple tool written in Go. 
Its main purpose is to extract logs, parse and transform them and load into DuckDB persistent database file for further exploration and analysis.

## Initial Setup

First step is to initiate `glog`. This can be done with the following command:
```
glog init
```
This creates a directory named `.glog` in your `$HOME` directory. `.glog` directory has the following structure:
```
.glog
├── data
│   ├── db
│   ├── processed
│   └── staged
├── etc
│   └── glog.toml
└── log
```
Folders are intended for the following use:
- `etc` contains `glog.toml` configuration file
- `data` folder is split into three parts:
  - `staged` - raw logs are axtracted into this folder, unarchived if needed and prepared for parsing
  - `processed` - logs are parsed and transformed/enriched and stored in this folder as JSON files (for easier processing by DuckDB)
  - `db` - final DuckDB `.duckdb` file is stored here
- `log` - this folder contains `glog` internal logs that track extraction, parsing and loading

### Configuration File

`glog.toml` looks like this:

```toml
AppHome = "/home/user/.glog"
AppConfig = "/home/user/.glog/etc"
AppDataStaging = "/home/user/.glog/data/staged"
AppDataProcessed = "/home/user/.glog/data/processed"
AppDataDB = "/home/user/.glog/data/db"
AppLog = "/home/user/.glog/log"

[[DefinedLogs]]
  logName = "nginx"
  sourceDir = "nginx_logs"
  logFilePattern = "access"
  logFileCompression = ".gz"

[[DefinedLogs]]
  logName = "nginx"
  sourceDir = "nginx_logs"
  logFilePattern = "access"
  logFileCompression = ".gz"
```

It is generated automatically with `glog init`. Fields starting with `App` contain the ingormation about `glog` directory paths.
The main action happens in `[[DefinedLogs]]` section. Placeholder log specifications are added automatically. You can edit them to adjust to your needs.
The meaning of the fileds is as follows:
- `logName` - identifies the name of a service. E.g. `nginx`, `caddy`, `ufw`
- `sourceDir` - this must be an absolute path to the **folder** where your raw logs are, e.g. `"/home/user/my-nginx-logs"`, `"~/my-nginx-logs"`, `"$HOME/my-nginx-logs"` all should work. Relative paths won't work. `glog` will cleanup any trailing spaces and `/`
- `logFilePattern` - this is used to tell `glog` which files it must look for when extracting. It should be **basename without file extension**. E.g. for `nginx` or `caddy` this field should be `"access"`
- `logFileCompression` - contains an compression extension, like `".gz"`

You can create additional `[[DefinedLogs]]` sections, just copy the existing ones. Placeholders will be silently skipped by `glog` so you can keep them as a template for new sections.

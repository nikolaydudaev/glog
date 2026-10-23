# Extract, Parse, Transform and Load Your logs

`glog` is a simple tool written in Go. 
Its main purpose is to extract logs, parse and transform them and load into DuckDB persistent database file for further exploration and analysis.

**IMPORTANT**:  
It currently only works with `nginx`, `caddy`, `ufw`, `auth` anf `fail2ban` logs. The readers and writers, and parsers are somewhat _hard-coded_.
I will probably revisit this at a later time, but I am too eager to finish some working version of this tool to start actual analysis sooner!

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

### Extraction

The extraction process is based on an idea that you store your logs somewhere on your machine. `glog` does not bother to check - as of this moment - if user has access to files and what kind of access that is. Consider copying files into the directory to which you have full read, write access to avoid any problems.

Each log record has `logSourceDir` defined. `glog` will scan that directory, looking for files identified by `logFilePattern` and by `logFileCompression`, will unpack archived files if necessary and concatenate the content of all those files into one staging file which will appear in `.glog/data/staged` folder with the file name consisting of `logName_staged.log`.
The extraction is also logged into `glog` internal logs, you can find it in `.glog/log/extraction-logs.log`. Internal logs are just JSON, so you can export it to e.g. DuckDB database.

I am thinking of adding more internal logging but as with this entire project, I am eager to start working with data itself, so I will revisit `glog` enhancements later.

### Parsing

Parsing is structured like this. I split logs into different types of events, like `ServerEvent`, `FirewallEvent` etc. Each such event has an `interface` which, so far, defines a dedicated `Parse` method. `Parse` methods are different depending on log type, e.g. `nginx` events use `regexp` to extract relevant information while `caddy` logs `Parse` is simply extracting `caddy` JSON into appropriate Go structs. Parsers for `ufw`, `fail2ban` and `auth` are a combination of `regexp` and strings splitting.

### Export

Export process is split into two parts conceptually, by converting parsed logs into JSON format and persistent DuckDB database.

#### JSON

As I am using Observable notebooks for my data exploration and analysis, JSON can be an easy format to work with in JS, `arquero` or `D3.js`. JSON files are stored in `.glog/data/processed` directory.

#### DB

`glog` also creates persistent DuckDB file `logs.duckdb` extension and stores it in `.glog/data/db`.
Generated `logs.duckdb` contains following tables:
- `glog_axtraction` - contains `glog` internal extraction logs
- `server_logs` - contains `nginx` and `caddy` logs
- `firewall_logs` - contains `ufw` logs
- `fail2ban_logs` - contains `fail2ban` logs
- `auth_logs` - contains `auth` logs

## Installation

Clone the repository and build like this:
```
git clone https://git.nikolaydudaev.com/ndudaev/glog.git
cd glog
go build -o ./build/glog
```

Then call:
```
./build/glog init
```

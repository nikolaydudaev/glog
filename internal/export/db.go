package export

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"git.nikolaydudaev.com/ndudaev/glog/internal/config"
	"git.nikolaydudaev.com/ndudaev/glog/internal/utils"
	_ "github.com/duckdb/duckdb-go/v2"
)

func WriteToDuckDB(cfg *config.AppConfig) (err error) {
	db, err := sql.Open("duckdb", filepath.Join(cfg.AppDataDB, "logs.duckdb"))
	if err != nil {
		return err
	}
	defer db.Close()

	var queries []string = []string{
		fmt.Sprintf(`COPY (SELECT DISTINCT * FROM read_ndjson_auto('%s')) TO '%s' (FORMAT parquet, COMPRESSION zstd);`, filepath.Join(cfg.AppDataProcessed, "server-events.json"), filepath.Join(cfg.AppDataProcessed, "server-logs.parquet")),
		fmt.Sprintf(`COPY (SELECT DISTINCT * FROM read_ndjson_auto('%s')) TO '%s' (FORMAT parquet, COMPRESSION zstd);`, filepath.Join(cfg.AppDataProcessed, "firewall-events.json"), filepath.Join(cfg.AppDataProcessed, "firewall-logs.parquet")),
		fmt.Sprintf(`COPY (SELECT DISTINCT * FROM read_ndjson_auto('%s')) TO '%s' (FORMAT parquet, COMPRESSION zstd);`, filepath.Join(cfg.AppDataProcessed, "fail2ban-events.json"), filepath.Join(cfg.AppDataProcessed, "fail2ban-logs.parquet")),
		fmt.Sprintf(`COPY (SELECT DISTINCT * FROM read_ndjson_auto('%s')) TO '%s' (FORMAT parquet, COMPRESSION zstd);`, filepath.Join(cfg.AppDataProcessed, "auth-events.json"), filepath.Join(cfg.AppDataProcessed, "auth-logs.parquet")),
		fmt.Sprintf(`CREATE OR REPLACE TABLE glog_extraction AS SELECT DISTINCT * FROM read_ndjson_auto('%s');`, filepath.Join(cfg.AppLog, "extraction-logs.log")),
		fmt.Sprintf(`CREATE OR REPLACE TABLE server_logs AS SELECT DISTINCT * FROM read_ndjson_auto('%s');`, filepath.Join(cfg.AppDataProcessed, "server-events.json")),
		fmt.Sprintf(`CREATE OR REPLACE TABLE firewall_logs AS SELECT DISTINCT * FROM read_ndjson_auto('%s');`, filepath.Join(cfg.AppDataProcessed, "firewall-events.json")),
		fmt.Sprintf(`CREATE OR REPLACE TABLE fail2ban_logs AS SELECT DISTINCT * FROM read_ndjson_auto('%s');`, filepath.Join(cfg.AppDataProcessed, "fail2ban-events.json")),
		fmt.Sprintf(`CREATE OR REPLACE TABLE auth_logs AS SELECT DISTINCT * FROM read_ndjson_auto('%s');`, filepath.Join(cfg.AppDataProcessed, "auth-events.json")),
	}
	bar := utils.CreateProgressBar(len(cfg.DefinedLogs), "2", "2", "Writing to DB...")
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			return err
		}
		bar.Add(1)
	}

	fmt.Println("Persistent DuckDB datavase created.")

	return err
}

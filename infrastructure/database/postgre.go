package database

import (
	"database/sql"
	"eth_fetcher/infrastructure/config"
	"eth_fetcher/infrastructure/logger"
)

func Init(cfg *config.Application, l logger.ILogger) *sql.DB {
	connStr := cfg.DSN
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		l.Fatalw("connecting to database", "error", err, "connection_string", connStr)
	}

	return db
}

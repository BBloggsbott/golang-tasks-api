package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/BBloggsbott/task-api/internal/config"
	"github.com/go-sql-driver/mysql"
)

func NewMySQLDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	mysql_cfg := mysql.NewConfig()
	mysql_cfg.User = cfg.User
	mysql_cfg.Passwd = cfg.Password
	mysql_cfg.Net = "tcp"
	mysql_cfg.Addr = fmt.Sprintf("%q:%q", cfg.Host, cfg.Port)
	mysql_cfg.DBName = cfg.Database

	db, err := sql.Open("mysql", mysql_cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxLifetime)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the db: %w", err)
	}

	return db, nil
}

func Close(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close mysql connection: %w", err)
	}
	return nil
}

func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("errro during mysql healthcheck: %w", err)
	}
	return nil
}

package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"sitecrawler/newgo/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDatabase(ctx context.Context, params *config.DBConnectionParams) (*sql.DB, error) {
	if params == nil {
		return nil, fmt.Errorf("db connection params are required")
	}

	user := url.QueryEscape(params.User)
	password := url.QueryEscape(params.Password)
	host := params.Host
	port := params.DBPort
	dbName := params.DBName

	query := url.Values{}
	query.Set("sslmode", "disable")
	if params.DBSchema != "" {
		query.Set("search_path", params.DBSchema)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", user, password, host, port, dbName, query.Encode())

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

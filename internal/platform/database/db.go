package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/CALEXCO/minisocial/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib" // register driver
)

func Init(config config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.Config.DSN)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

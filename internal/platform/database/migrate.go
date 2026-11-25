package database

import (
	"database/sql"

	"github.com/CALEXCO/minisocial/migrations"
	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations.EmbededMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(db, ".")
}

func Down(db *sql.DB, dir string) error {
	return goose.Down(db, dir)
}

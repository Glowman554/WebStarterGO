package database

import (
	"database/sql"
	"log/slog"
	"os"
	"strings"

	_ "github.com/tursodatabase/go-libsql"
)

var DB *sql.DB

func WithDatabase(next func(db *sql.DB) error) error {
	db, err := sql.Open("libsql", "file:local.db")
	if err != nil {
		return err
	}
	defer db.Close()

	DB = db

	return next(db)
}

func ApplyMigration(db *sql.DB, file string) error {
	_, err := db.Exec("create table if not exists __migrations (id text primary key)")
	if err != nil {
		return err
	}

	var count int
	err = db.QueryRow("select count(*) from __migrations where id = $1", file).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		slog.Info("Applying migartion " + file)
		migration, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		migrations := strings.Split(string(migration), ";")

		for _, v := range migrations {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}

			_, err = db.Exec(v)
			if err != nil {
				return err
			}
		}

		_, err = db.Exec("insert into __migrations (id) values ($1)", file)
		if err != nil {
			return err
		}
	}
	return nil
}

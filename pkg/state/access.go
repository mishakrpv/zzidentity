package state

import (
	"embed"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zzidentity/zzidentity/pkg/config"
)

// Embed migrations into the Go binary
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

func Open(postgresCfg *config.Postgres) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", postgresCfg.URL)
	if err != nil {
		return nil, err
	}

	files, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		m, err := migrationFiles.ReadFile("migrations/" + file.Name())
		if err != nil {
			return nil, err
		}

		db.MustExec(string(m))
	}

	return db, nil
}
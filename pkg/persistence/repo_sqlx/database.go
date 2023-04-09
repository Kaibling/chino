package repo_sqlx

import (
	"context"
	"time"

	"chino/pkg/config"
	"chino/pkg/log"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func InitDatabase(cfg config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", cfg.DatabaseFilePath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB(ctx context.Context, db *sqlx.DB) error {
	schemaMovie := `CREATE TABLE movies (
		id text PRIMARY KEY,
		name text UNIQUE,
		created_at datetime,
		created_by text,
		genre text,
		release_date datetime,
		notified bool
		);`

	_, err := db.Exec(schemaMovie)
	if err != nil {
		return err
	}
	log.Info(ctx, "movies created")
	return nil
}

type dbModel struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_by"`
}

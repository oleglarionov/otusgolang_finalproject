package sql

import (
	"log"

	"github.com/jmoiron/sqlx"

	// init driver.
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func NewDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Println("connection to database established")
	return db, nil
}

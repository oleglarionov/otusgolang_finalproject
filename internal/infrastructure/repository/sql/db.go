package sql

import (
	"log"

	"github.com/jmoiron/sqlx"

	// init driver.
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type DBConnector interface {
	GetConn() (*sqlx.DB, error)
	CloseConn() error
}

type DBConnectorImpl struct {
	dsn string
	db  *sqlx.DB
}

func NewDBConnectorImpl(dsn string) *DBConnectorImpl {
	return &DBConnectorImpl{dsn: dsn}
}

func (c *DBConnectorImpl) GetConn() (*sqlx.DB, error) {
	if c.db == nil {
		var err error
		c.db, err = c.reconnect()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return c.db, nil

}

func (c *DBConnectorImpl) CloseConn() error {
	if c.db != nil {
		err := c.db.Close()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (c *DBConnectorImpl) reconnect() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", c.dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Println("connection to database established")
	return db, nil
}

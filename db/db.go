package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DB struct {
	pool   *pgxpool.Pool
	logger *logrus.Entry
}

func New() (*DB, error) {
	db := &DB{
		logger: logrus.WithField("module", "db"),
	}

	pgxCfg, err := db.pgxPoolConfig()
	if err != nil {
		return nil, errors.Wrap(err, "could not build pgx config")
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), pgxCfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to pgx pool")
	}

	db.pool = pool

	return db, nil
}

func (db *DB) Connection() *pgxpool.Pool {
	return db.pool
}

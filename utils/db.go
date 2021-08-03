package utils

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/db"
)

func GetHighestBlock(ctx context.Context, db *db.DB) (*int64, error) {
	var number int64

	err := db.Connection().QueryRow(ctx, `select number from blocks order by number desc limit 1;`).Scan(&number)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "could not get highest block")
	}

	return &number, nil
}

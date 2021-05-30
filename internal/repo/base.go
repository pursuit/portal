package repo

import (
	"context"
	"database/sql"

	"github.com/pursuit/portal/internal"
)

//go:generate mockgen -source=base.go -destination=mock/base.go

type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func Transaction(ctx context.Context, db *sql.DB, fn func(DB) *internal.E) *internal.E {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return &internal.E{
			Err:    err,
			Status: 503,
		}
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return &internal.E{
			Err:    err,
			Status: 503,
		}
	}

	return nil
}

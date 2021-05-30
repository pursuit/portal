package repo

import (
	"context"
	"database/sql"
)

//go:generate mockgen -source=base.go -destination=mock/base.go

type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

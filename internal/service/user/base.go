package user

import (
	"context"
	"database/sql"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/repo"
)

//go:generate mockgen -source=base.go -destination=mock/base.go

type Service interface {
	Create(ctx context.Context, username string, password []byte) (int, *internal.E)
	Login(ctx context.Context, username string, password []byte) (string, *internal.E)
}

type Svc struct {
	DB       *sql.DB
	UserRepo repo.User
}

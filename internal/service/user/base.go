package user

import (
	"context"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/repo"
)

//go:generate mockgen -source=base.go -destination=mock/base.go

type Service interface {
	Create(ctx context.Context, username string, password []byte) *internal.E
}

type Svc struct {
	DB       repo.DB
	UserRepo repo.User
}

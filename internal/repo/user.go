package repo

import (
	"context"
	"time"

	"github.com/pursuit/portal/internal"
)

//go:generate mockgen -source=user.go -destination=mock/user.go

type User interface {
	Create(ctx context.Context, db DB, username string, hashedPassword []byte, now time.Time) (int, *internal.E)
}

type UserRepo struct {
}

func (this UserRepo) Create(ctx context.Context, db DB, username string, hashedPassword []byte, now time.Time) (int, *internal.E) {
	var id int
	if err := db.QueryRowContext(ctx, "INSERT INTO users (username,hashed_password,created_at) VALUES($1,$2,$3) RETURNING id", username, hashedPassword, now).Scan(&id); err != nil {
		return 0, &internal.E{
			Err:    err,
			Status: internal.EDbProblem,
		}
	}

	return id, nil
}

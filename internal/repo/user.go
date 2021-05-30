package repo

import (
	"context"
	"time"
)

//go:generate mockgen -source=user.go -destination=mock/user.go

type User interface {
	Create(ctx context.Context, db DB, username string, hashedPassword []byte, now time.Time) (int, error)
}

type UserRepo struct {
}

func (this UserRepo) Create(ctx context.Context, db DB, username string, hashedPassword []byte, now time.Time) (int, error) {
	var id int
	err := db.QueryRowContext(ctx, "INSERT INTO users (username,hashed_password,created_at) VALUES($1,$2,$3) RETURNING id", username, hashedPassword, now).Scan(&id)
	return id, err
}

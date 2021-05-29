package repo

import (
	"context"
)

//go:generate mockgen -source=user.go -destination=mock/user.go

type User interface {
	Create(ctx context.Context, db DB, username string, hashedPassword []byte) (int, error)
}

type UserRepo struct {
}

func (this UserRepo) Create(ctx context.Context, db DB, username string, hashedPassword []byte) (int, error) {
	var id int
	err := db.QueryRowContext(ctx, "INSERT INTO users (username,hashed_password) VALUES($1,$2) RETURNING id", username, hashedPassword).Scan(&id)
	return id, err
}

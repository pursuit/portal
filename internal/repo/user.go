package repo

import (
	"context"
)

//go:generate mockgen -source=user.go -destination=mock/user.go

type User interface {
	Create(ctx context.Context, db DB, username string, hashedPassword []byte) error
}

type UserRepo struct {
}

func (this UserRepo) Create(ctx context.Context, db DB, username string, hashedPassword []byte) error {
	_, err := db.ExecContext(ctx, "INSERT INTO users (username,hashed_password) VALUES($1,$2)", username, hashedPassword)
	return err
}

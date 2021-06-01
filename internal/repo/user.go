package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/model"
)

//go:generate mockgen -source=user.go -destination=mock/user.go

type User interface {
	Create(ctx context.Context, db DB, username string, hashedPassword []byte, now time.Time) (int, *internal.E)
	GetByUsername(ctx context.Context, db DB, username string) (model.User, *internal.E)
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

func (this UserRepo) GetByUsername(ctx context.Context, db DB, username string) (model.User, *internal.E) {
	var user model.User
	if err := db.QueryRowContext(ctx, "SELECT id,hashed_password FROM users WHERE username = $1", username).Scan(&user.ID, &user.HashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return user, &internal.E{
				Err:    err,
				Status: internal.EUserNotFound,
			}
		}

		return user, &internal.E{
			Err:    err,
			Status: internal.EDbProblem,
		}
	}

	user.Username = username
	return user, nil
}

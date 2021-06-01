package repo

import (
	"context"
	"time"
)

//go:generate mockgen -source=mutation.go -destination=mock/mutation.go

type Mutation interface {
	Create(ctx context.Context, db DB, userID int, referenceID int, referenceType string, amount int, createdAt time.Time) (int, error)
	GetBalance(ctx context.Context, db DB, userID int) (int, error)
}

type MutationRepo struct {
}

func (this MutationRepo) Create(ctx context.Context, db DB, userID int, referenceID int, referenceType string, amount int, createdAt time.Time) (int, error) {
	var id int
	err := db.QueryRowContext(ctx, "INSERT INTO mutations (user_id,reference_id,reference_type,amount,created_at) VALUES($1,$2,$3,$4,$5) RETURNING id", userID, referenceID, referenceType, amount, createdAt).Scan(&id)
	return id, err
}

func (this MutationRepo) GetBalance(ctx context.Context, db DB, userID int) (int, error) {
	var balance int
	err := db.QueryRowContext(ctx, "select coalesce(sum(amount), 0) as balance from mutations where user_id = $1", userID).Scan(&balance)
	return balance, err
}

package repo

import (
	"context"
	"time"

	"github.com/pursuit/portal/internal"
)

//go:generate mockgen -source=mutation.go -destination=mock/mutation.go

type Mutation interface {
	Create(ctx context.Context, db DB, userID int, referenceID int, referenceType string, amount int, createdAt time.Time) (int, *internal.E)
	GetBalance(ctx context.Context, db DB, userID int) (int, *internal.E)
}

type MutationRepo struct {
}

func (this MutationRepo) Create(ctx context.Context, db DB, userID int, referenceID int, referenceType string, amount int, createdAt time.Time) (int, *internal.E) {
	var id int
	if err := db.QueryRowContext(ctx, "INSERT INTO mutations (user_id,reference_id,reference_type,amount,created_at) VALUES($1,$2,$3,$4,$5) RETURNING id", userID, referenceID, referenceType, amount, createdAt).Scan(&id); err != nil {
		return 0, &internal.E{
			Err:    err,
			Status: internal.EDbProblem,
		}
	}

	return id, nil
}

func (this MutationRepo) GetBalance(ctx context.Context, db DB, userID int) (int, *internal.E) {
	var balance int
	if err := db.QueryRowContext(ctx, "SELECT COALESCE(SUM(amount), 0) as balance from mutations where user_id = $1", userID).Scan(&balance); err != nil {
		return 0, &internal.E{
			Err:    err,
			Status: internal.EDbProblem,
		}
	}

	return balance, nil
}

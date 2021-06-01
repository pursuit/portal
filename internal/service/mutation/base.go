package mutation

import (
	"context"
	"errors"
	"time"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/config"
	"github.com/pursuit/portal/internal/repo"
)

//go:generate mockgen -source=base.go -destination=mock/base.go

type Service interface {
	Create(ctx context.Context, userID int, referenceID int, referenceType string, amount int) *internal.E
	GetBalance(ctx context.Context, userID int) (int, *internal.E)
}

type Svc struct {
	DB           repo.DB
	MutationRepo repo.Mutation
}

func (this Svc) Create(ctx context.Context, userID int, referenceID int, referenceType string, amount int) *internal.E {
	now := time.Now().UTC()
	_, validType := config.MutationReferences[referenceType]
	if !validType {
		return &internal.E{
			Err:    errors.New("invalid reference"),
			Status: internal.EMutationRefNotFound,
		}
	}

	_, err := this.MutationRepo.Create(ctx, this.DB, userID, referenceID, referenceType, amount, now)
	if err != nil {
		return err
	}

	return nil
}

func (this Svc) GetBalance(ctx context.Context, userID int) (int, *internal.E) {
	balance, err := this.MutationRepo.GetBalance(ctx, this.DB, userID)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

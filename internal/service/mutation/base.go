package mutation

import (
	"context"
	"errors"
	"time"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/config"
	"github.com/pursuit/portal/internal/repo"
)

type Service interface {
	Create(ctx context.Context, userID int, referenceID int, referenceType string, amount int) *internal.E
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

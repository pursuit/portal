package mutation

import (
	"context"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/repo"
)

type Service interface {
	Create(ctx context.Context, userID int, referenceID int, referenceType string, amount int) *internal.E
}

type Svc struct {
	DB           repo.DB
	MutationRepo repo.Mutation
}

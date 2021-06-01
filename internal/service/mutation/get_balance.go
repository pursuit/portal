package mutation

import (
	"context"

	"github.com/pursuit/portal/internal"
)

func (this Svc) GetBalance(ctx context.Context, userID int) (int, *internal.E) {
	balance, err := this.MutationRepo.GetBalance(ctx, this.DB, userID)
	if err != nil {
		return 0, &internal.E{
			Err:    err,
			Status: 503,
		}
	}

	return balance, nil
}

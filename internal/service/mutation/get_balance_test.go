package mutation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pursuit/portal/internal/repo/mock"
	"github.com/pursuit/portal/internal/service/mutation"

	"github.com/golang/mock/gomock"
)

func TestGetBalance(t *testing.T) {
	for _, testcase := range []struct {
		tName string

		balance         int
		connErr    error
		outputErr  error
	}{
		{
			tName:         "failed conn",
			connErr:    errors.New("failed connection"),
			outputErr:     errors.New("failed connection"),
		},
		{
			tName:         "success",
			balance: 50,
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			mocker := gomock.NewController(t)

			db := mock_repo.NewMockDB(mocker)
			repo := mock_repo.NewMockMutation(mocker)

			repo.EXPECT().GetBalance(gomock.Any(), db, 2).Return(testcase.balance, testcase.connErr)

			balance, err := mutation.Svc{
				DB:           db,
				MutationRepo: repo,
			}.GetBalance(context.Background(), 2)

			if (testcase.outputErr == nil && err != nil) ||
				(testcase.outputErr != nil && err == nil) ||
				(err != nil && testcase.outputErr.Error() != err.Error()) {
				t.Errorf("Test %s, err is %v, should be %v", testcase.tName, err, testcase.outputErr)
			}

			if testcase.balance != balance {
				t.Errorf("Test %s, balance is %d, should be %d", testcase.tName, balance, testcase.balance)
			}
		})
	}
}

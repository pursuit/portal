package mutation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/repo/mock"
	"github.com/pursuit/portal/internal/service/mutation"

	"github.com/golang/mock/gomock"
)

func TestCreate(t *testing.T) {
	for _, testcase := range []struct {
		tName string

		referenceType string

		validType bool

		persistErr *internal.E
		outputErr  error
	}{
		{
			tName:         "invalid reference",
			referenceType: "invalid type",
			outputErr:     errors.New("invalid reference"),
		},
		{
			tName:         "failed persist",
			referenceType: "free_coins",
			validType:     true,
			persistErr:    &internal.E{Err: errors.New("failed persist")},
			outputErr:     errors.New("failed persist"),
		},
		{
			tName:         "success",
			referenceType: "free_coins",
			validType:     true,
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			mocker := gomock.NewController(t)

			db := mock_repo.NewMockDB(mocker)
			repo := mock_repo.NewMockMutation(mocker)

			if testcase.validType {
				repo.EXPECT().Create(gomock.Any(), db, 2, 5, testcase.referenceType, 50, gomock.Any()).Return(1, testcase.persistErr)
			}

			err := mutation.Svc{
				DB:           db,
				MutationRepo: repo,
			}.Create(context.Background(), 2, 5, testcase.referenceType, 50)

			if (testcase.outputErr == nil && err != nil) ||
				(testcase.outputErr != nil && err == nil) ||
				(err != nil && testcase.outputErr.Error() != err.Error()) {
				t.Errorf("Test %s, err is %v, should be %v", testcase.tName, err, testcase.outputErr)
			}
		})
	}
}

func TestGetBalance(t *testing.T) {
	for _, testcase := range []struct {
		tName string

		dbErr     *internal.E
		outputErr error
		balance   int
	}{
		{
			tName:     "failed db",
			dbErr:     &internal.E{Err: errors.New("failed conn")},
			outputErr: errors.New("failed conn"),
		},
		{
			tName:   "success",
			balance: 500,
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			mocker := gomock.NewController(t)

			db := mock_repo.NewMockDB(mocker)
			repo := mock_repo.NewMockMutation(mocker)

			repo.EXPECT().GetBalance(gomock.Any(), db, 2).Return(testcase.balance, testcase.dbErr)

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

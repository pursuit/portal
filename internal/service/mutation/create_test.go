package mutation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pursuit/portal/internal/repo/mock"
	"github.com/pursuit/portal/internal/service/mutation"

	"github.com/golang/mock/gomock"
)

func TestCreate(t *testing.T) {
	for _, testcase := range []struct {
		tName string

		referenceType string

		validType bool

		persistErr error
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
			persistErr:    errors.New("failed persist"),
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

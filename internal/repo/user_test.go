package repo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pursuit/portal/internal/repo"
	"github.com/pursuit/portal/internal/repo/mock"
	"github.com/pursuit/portal/mock/database/sql"

	"github.com/golang/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	username := "username"
	password := []byte("password")

	for _, testcase := range []struct {
		tName string

		connErr error
		idErr   error

		outputErr error
	}{
		{
			tName:     "cant access db",
			connErr:   errors.New("failed timeout"),
			outputErr: errors.New("failed timeout"),
		},
		{
			tName: "success",
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			mocker := gomock.NewController(t)
			defer mocker.Finish()

			db := mock_repo.NewMockDB(mocker)
			result := mock_sql.NewMockResult(mocker)

			db.EXPECT().ExecContext(gomock.Any(), gomock.Any(), username, password).Return(result, testcase.connErr)

			err := repo.UserRepo{}.Create(context.Background(), db, username, password)
			if (testcase.outputErr == nil && err != nil) ||
				(testcase.outputErr != nil && err == nil) ||
				(err != nil && testcase.outputErr.Error() != err.Error()) {
				t.Errorf("Test %s, err is %v, should be %v", testcase.tName, err, testcase.outputErr)
			}
		})
	}
}

package repo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pursuit/portal/internal/repo"
	"github.com/pursuit/portal/internal/repo/mock"

	"github.com/golang/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	username := "username"
	password := []byte("password")

	for _, testcase := range []struct {
		tName string

		connErr error
		id      int

		outputErr error
	}{
		{
			tName:     "cant access db",
			connErr:   errors.New("failed timeout"),
			outputErr: errors.New("failed timeout"),
		},
		{
			tName: "success",
			id: 8,
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			mocker := gomock.NewController(t)
			defer mocker.Finish()

			db := mock_repo.NewMockDB(mocker)
			row := mock_repo.NewMockRow(mocker)

			db.EXPECT().QueryRowContext(gomock.Any(), "INSERT INTO users (username,hashed_password) VALUES($1,$2) RETURNING id", username, password).Return(row)
			row.EXPECT().Scan(gomock.Any()).DoAndReturn(func (id *int) error {
				*id = testcase.id
				return testcase.connErr
			})

			id, err := repo.UserRepo{}.Create(context.Background(), db, username, password)
			if (testcase.outputErr == nil && err != nil) ||
				(testcase.outputErr != nil && err == nil) ||
				(err != nil && testcase.outputErr.Error() != err.Error()) {
				t.Errorf("Test %s, err is %v, should be %v", testcase.tName, err, testcase.outputErr)
			}

			if testcase.id != id {
				t.Errorf("Test %s, id is %d, should be %d", testcase.tName, id, testcase.id)
			}
		})
	}
}

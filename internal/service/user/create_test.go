package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/repo/mock"
	"github.com/pursuit/portal/internal/service/user"

	"github.com/golang/mock/gomock"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	for _, testcase := range []struct {
		tName string

		username string
		password []byte

		validInput bool

		persistErr *internal.E

		outputErr error
		outputID  int
	}{
		{
			tName:     "username len is not too small",
			username:  "a",
			password:  []byte("valids"),
			outputErr: errors.New("invalid input"),
		},
		{
			tName:     "username len is not too big",
			username:  "a123456678902213123",
			password:  []byte("valids"),
			outputErr: errors.New("invalid input"),
		},
		{
			tName:     "username contain special char",
			username:  "a123?4566",
			password:  []byte("valids"),
			outputErr: errors.New("invalid input"),
		},
		{
			tName:     "password is not valid char",
			username:  "a1234566",
			password:  []byte{'?'},
			outputErr: errors.New("invalid input"),
		},
		{
			tName:     "password len is too small",
			username:  "a1234566",
			password:  []byte("a"),
			outputErr: errors.New("invalid input"),
		},
		{
			tName:     "password len is too big",
			username:  "a1234566",
			password:  []byte("aqwezcxadqwewqwedas"),
			outputErr: errors.New("invalid input"),
		},
		{
			tName:      "fail to persist",
			username:   "a1234566",
			password:   []byte("valids"),
			validInput: true,
			persistErr: &internal.E{Err: errors.New("timeout persist")},
			outputErr:  errors.New("timeout persist"),
		},
		{
			tName:      "success",
			username:   "a1234566",
			password:   []byte("valids"),
			validInput: true,
			outputID:   8,
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			mocker := gomock.NewController(t)
			defer mocker.Finish()

			db, mock, dbErr := sqlmock.New()
			if dbErr != nil {
				panic(dbErr)
			}
			defer db.Close()

			mock.ExpectBegin()
			mock.ExpectExec(`INSERT INTO events \(topic,payload\) VALUES\(\$1,\$2\)`).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			userRepo := mock_repo.NewMockUser(mocker)

			if testcase.validInput {
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any(), testcase.username, gomock.Any(), gomock.Any()).Return(testcase.outputID, testcase.persistErr)
			}

			id, err := user.Svc{
				UserRepo: userRepo,
				DB:       db,
			}.Create(context.Background(), testcase.username, testcase.password)

			if (testcase.outputErr == nil && err != nil) ||
				(testcase.outputErr != nil && err == nil) ||
				(err != nil && testcase.outputErr.Error() != err.Error()) {
				t.Errorf("Test %s, err is %v, should be %v", testcase.tName, err, testcase.outputErr)
			}

			if testcase.outputID != id {
				t.Errorf("Test %s, id is %d, should be %d", testcase.tName, id, testcase.outputID)
			}

			for _, ch := range testcase.password {
				if ch != 0 {
					t.Errorf("Test %s, password is not cleared", testcase.tName)
				}
			}
		})
	}
}

package server_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/proto/out"
	"github.com/pursuit/portal/internal/proto/server"
	"github.com/pursuit/portal/internal/service/user/mock"

	"github.com/golang/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	for _, testcase := range []struct {
		tName string

		input     *proto.CreateUserPayload
		createErr *internal.E

		outputErr error
	}{
		{
			tName: "failed to create, client error",
			createErr: &internal.E{
				Err:    errors.New("invalid username"),
				Status: 422,
			},
			outputErr: errors.New("rpc error: code = InvalidArgument desc = invalid username"),
		},
		{
			tName: "failed to create, server error",
			createErr: &internal.E{
				Err:    errors.New("database down"),
				Status: 503,
			},
			outputErr: errors.New("rpc error: code = Unavailable desc = Please try again in a few moment"),
		},
		{
			tName: "success",
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			mocker := gomock.NewController(t)
			defer mocker.Finish()

			svc := mock_user.NewMockService(mocker)

			svc.EXPECT().Create(gomock.Any(), testcase.input.GetName(), testcase.input.GetPassword()).Return(testcase.createErr)

			_, err := server.UserServer{
				UserService: svc,
			}.Create(context.Background(), testcase.input)

			if (testcase.outputErr == nil && err != nil) ||
				(testcase.outputErr != nil && err == nil) ||
				(err != nil && testcase.outputErr.Error() != err.Error()) {
				t.Errorf("Test %s, err is %v, should be %v", testcase.tName, err, testcase.outputErr)
			}
		})
	}
}

package server_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pursuit/portal/internal"
	"github.com/pursuit/portal/internal/proto/out/api/portal"
	"github.com/pursuit/portal/internal/proto/server"
	"github.com/pursuit/portal/internal/service/mutation/mock"
	"github.com/pursuit/portal/internal/service/user/mock"

	"github.com/golang/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	for _, testcase := range []struct {
		tName string

		input     *portal_proto.CreateUserPayload
		createErr *internal.E

		outputID  int
		outputErr error
	}{
		{
			tName: "failed to create, client error",
			createErr: &internal.E{
				Err:    errors.New("invalid username"),
				Status: 422_000,
			},
			outputErr: errors.New("rpc error: code = InvalidArgument desc = invalid username"),
		},
		{
			tName: "failed to create, server error",
			createErr: &internal.E{
				Err:    errors.New("database down"),
				Status: 503_000,
			},
			outputErr: errors.New("rpc error: code = Unavailable desc = Please try again in a few moment"),
		},
		{
			tName:    "success",
			outputID: 7,
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			mocker := gomock.NewController(t)
			defer mocker.Finish()

			svc := mock_user.NewMockService(mocker)

			svc.EXPECT().Create(gomock.Any(), testcase.input.GetUsername(), testcase.input.GetPassword()).Return(testcase.outputID, testcase.createErr)

			resp, err := server.UserServer{
				UserService: svc,
			}.Create(context.Background(), testcase.input)

			if (testcase.outputErr == nil && err != nil) ||
				(testcase.outputErr != nil && err == nil) ||
				(err != nil && testcase.outputErr.Error() != err.Error()) {
				t.Errorf("Test %s, err is %v, should be %v", testcase.tName, err, testcase.outputErr)
			}

			if int64(testcase.outputID) != resp.GetId() {
				t.Errorf("Test %s, id is %d, should be %d", testcase.tName, resp.GetId(), testcase.outputID)
			}
		})
	}
}

func TestGetBalance(t *testing.T) {
	for _, testcase := range []struct {
		tName string

		input  *portal_proto.GetUserBalancePayload
		svcErr *internal.E

		outputErr error
		balance   int
	}{
		{
			tName: "failed to get, client error",
			svcErr: &internal.E{
				Err:    errors.New("not found user"),
				Status: 422_000,
			},
			outputErr: errors.New("rpc error: code = InvalidArgument desc = not found user"),
		},
		{
			tName: "failed to get, server error",
			svcErr: &internal.E{
				Err:    errors.New("database down"),
				Status: 503_000,
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

			svc := mock_mutation.NewMockService(mocker)

			svc.EXPECT().GetBalance(gomock.Any(), int(testcase.input.GetUserId())).Return(testcase.balance, testcase.svcErr)

			resp, err := server.UserServer{
				MutationService: svc,
			}.GetBalance(context.Background(), testcase.input)

			if (testcase.outputErr == nil && err != nil) ||
				(testcase.outputErr != nil && err == nil) ||
				(err != nil && testcase.outputErr.Error() != err.Error()) {
				t.Errorf("Test %s, err is %v, should be %v", testcase.tName, err, testcase.outputErr)
			}

			if int64(testcase.balance) != resp.GetAmount() {
				t.Errorf("Test %s, balance is %d, should be %d", testcase.tName, resp.GetAmount(), testcase.balance)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	for _, testcase := range []struct {
		tName string

		input  *portal_proto.LoginPayload
		svcErr *internal.E

		outputToken string
		outputErr   error
	}{
		{
			tName: "failed to create, client error",
			svcErr: &internal.E{
				Err:    errors.New("invalid username"),
				Status: 422_000,
			},
			outputErr: errors.New("rpc error: code = InvalidArgument desc = invalid username"),
		},
		{
			tName: "failed to create, server error",
			svcErr: &internal.E{
				Err:    errors.New("database down"),
				Status: 503_000,
			},
			outputErr: errors.New("rpc error: code = Unavailable desc = Please try again in a few moment"),
		},
		{
			tName:       "success",
			outputToken: "7a",
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			mocker := gomock.NewController(t)
			defer mocker.Finish()

			svc := mock_user.NewMockService(mocker)

			svc.EXPECT().Login(gomock.Any(), testcase.input.GetUsername(), testcase.input.GetPassword()).Return(testcase.outputToken, testcase.svcErr)

			resp, err := server.UserServer{
				UserService: svc,
			}.Login(context.Background(), testcase.input)

			if (testcase.outputErr == nil && err != nil) ||
				(testcase.outputErr != nil && err == nil) ||
				(err != nil && testcase.outputErr.Error() != err.Error()) {
				t.Errorf("Test %s, err is %v, should be %v", testcase.tName, err, testcase.outputErr)
			}

			if testcase.outputToken != resp.GetToken() {
				t.Errorf("Test %s, token is %s, should be %s", testcase.tName, resp.GetToken(), testcase.outputToken)
			}
		})
	}
}

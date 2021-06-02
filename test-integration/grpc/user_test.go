package grpc_test

import (
	"errors"
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/pursuit/portal/internal/proto/out"
)

func TestUser(t *testing.T) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":5001", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var successID int

	for _, testcase := range []struct {
		tName string

		username string
		password []byte

		outputErr error
	}{
		{
			tName:    "success",
			username: "Bambang",
			password: []byte("ZXCqweqwe"),
		},
		{
			tName:     "too small username",
			username:  "B",
			password:  []byte("ZXCqweqwe"),
			outputErr: errors.New("rpc error: code = InvalidArgument desc = invalid input"),
		},
		{
			tName:     "too big username",
			username:  "Bqweqewqweqweqdsa",
			password:  []byte("ZXCqweqwe"),
			outputErr: errors.New("rpc error: code = InvalidArgument desc = invalid input"),
		},
		{
			tName:     "username contain special char",
			username:  "Bamb?ang",
			password:  []byte("ZXCqweqwe"),
			outputErr: errors.New("rpc error: code = InvalidArgument desc = invalid input"),
		},
		{
			tName:     "too small password",
			username:  "ZXCqweqwe",
			password:  []byte("B"),
			outputErr: errors.New("rpc error: code = InvalidArgument desc = invalid input"),
		},
		{
			tName:     "too big password",
			username:  "ZXCqweqwe",
			password:  []byte("Bqweqewqweqweqdsa"),
			outputErr: errors.New("rpc error: code = InvalidArgument desc = invalid input"),
		},
		{
			tName:     "password contain special char",
			username:  "ZXCqweqwe",
			password:  []byte("Bamb?ang"),
			outputErr: errors.New("rpc error: code = InvalidArgument desc = invalid input"),
		},
	} {
		t.Run(testcase.tName, func(t *testing.T) {
			c := proto.NewUserClient(conn)
			resp, err := c.Create(context.Background(), &proto.CreateUserPayload{
				Name:     testcase.username,
				Password: testcase.password,
			})

			if (testcase.outputErr == nil && err != nil) ||
				(testcase.outputErr != nil && err == nil) ||
				(err != nil && testcase.outputErr.Error() != err.Error()) {
				t.Errorf("Test %s, err is %v, should be %v", testcase.tName, err, testcase.outputErr)
			}

			if err == nil {
				successID = int(resp.GetId())
			}
		})
	}

	time.Sleep(5 * time.Second)
	testGetUserBalanceValid(t, successID)
}

func testGetUserBalanceValid(t *testing.T, userID int) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":5001", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewUserClient(conn)
	resp, err := c.GetBalance(context.Background(), &proto.GetUserBalancePayload{
		UserId: int64(userID),
	})

	if err != nil {
		t.Errorf("Test created get user balance got error %v", err)
	}

	if resp.GetAmount() != int64(10) {
		t.Errorf("Test created get user balance is %d, should be 10", resp.GetAmount())
	}
}

func TestGetNotExistingUserBalance(t *testing.T) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":5001", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewUserClient(conn)
	resp, err := c.GetBalance(context.Background(), &proto.GetUserBalancePayload{
		UserId: int64(0),
	})

	if err != nil {
		t.Errorf("Test not exist get user balance got error %v", err)
	}

	if resp.GetAmount() != int64(0) {
		t.Errorf("Test not exists get user balance is %d, should be 0", resp.GetAmount())
	}
}

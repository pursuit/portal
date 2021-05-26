package test_integration_test

import (
	"errors"
	"testing"

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
			_, err = c.Create(context.Background(), &proto.CreateUserPayload{
				Name:     testcase.username,
				Password: testcase.password,
			})

			if (testcase.outputErr == nil && err != nil) ||
				(testcase.outputErr != nil && err == nil) ||
				(err != nil && testcase.outputErr.Error() != err.Error()) {
				t.Errorf("Test %s, err is %v, should be %v", testcase.tName, err, testcase.outputErr)
			}
		})
	}
}

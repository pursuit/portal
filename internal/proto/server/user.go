package server

import (
	"context"

	"github.com/pursuit/portal/internal/proto/out"
	"github.com/pursuit/portal/internal/service/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
)

type UserServer struct {
	proto.UnimplementedUserServer

	UserService user.Service
}

func (this UserServer) Create(ctx context.Context, in *proto.CreateUserPayload) (*empty.Empty, error) {
	if err := this.UserService.Create(ctx, in.GetName(), in.GetPassword()); err != nil {
		if err.Status >= 400 && err.Status < 500 {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}

		return nil, grpc.Errorf(codes.Unavailable, "Please try again in a few moment")
	}

	return &empty.Empty{}, nil
}

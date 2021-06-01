package server

import (
	"context"

	"github.com/pursuit/portal/internal/proto/out"
	"github.com/pursuit/portal/internal/service/mutation"
	"github.com/pursuit/portal/internal/service/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type UserServer struct {
	proto.UnimplementedUserServer

	UserService     user.Service
	MutationService mutation.Service
}

func (this UserServer) Create(ctx context.Context, in *proto.CreateUserPayload) (*proto.CreateUserResponse, error) {
	id, err := this.UserService.Create(ctx, in.GetName(), in.GetPassword())
	if err != nil {
		httpStatus := err.Status / 1_000
		if httpStatus >= 400 && httpStatus < 500 {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}

		return nil, grpc.Errorf(codes.Unavailable, "Please try again in a few moment")
	}

	return &proto.CreateUserResponse{Id: int64(id)}, nil
}

func (this UserServer) GetBalance(ctx context.Context, in *proto.GetUserBalancePayload) (*proto.GetUserBalanceResponse, error) {
	balance, err := this.MutationService.GetBalance(ctx, int(in.GetUserId()))
	if err != nil {
		httpStatus := err.Status / 1_000
		if httpStatus >= 400 && httpStatus < 500 {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}

		return nil, grpc.Errorf(codes.Unavailable, "Please try again in a few moment")
	}

	return &proto.GetUserBalanceResponse{Amount: int64(balance)}, nil
}

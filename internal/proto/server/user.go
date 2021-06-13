package server

import (
	"context"

	"github.com/pursuit/portal/internal/proto/out/api/portal"
	"github.com/pursuit/portal/internal/service/mutation"
	"github.com/pursuit/portal/internal/service/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	portal_proto.UnimplementedUserServer

	UserService     user.Service
	MutationService mutation.Service
}

func (this UserServer) Create(ctx context.Context, in *portal_proto.CreateUserPayload) (*portal_proto.CreateUserResponse, error) {
	id, err := this.UserService.Create(ctx, in.GetUsername(), in.GetPassword())
	if err != nil {
		httpStatus := err.Status / 1_000
		if httpStatus >= 400 && httpStatus < 500 {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	return &portal_proto.CreateUserResponse{Id: uint64(id)}, nil
}

func (this UserServer) GetBalance(ctx context.Context, in *portal_proto.GetUserBalancePayload) (*portal_proto.GetUserBalanceResponse, error) {
	balance, err := this.MutationService.GetBalance(ctx, int(in.GetUserId()))
	if err != nil {
		httpStatus := err.Status / 1_000
		if httpStatus >= 400 && httpStatus < 500 {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	return &portal_proto.GetUserBalanceResponse{Amount: uint64(balance)}, nil
}

func (this UserServer) Login(ctx context.Context, in *portal_proto.LoginPayload) (*portal_proto.LoginResponse, error) {
	token, err := this.UserService.Login(ctx, in.GetUsername(), in.GetPassword())
	if err != nil {
		httpStatus := err.Status / 1_000
		if httpStatus >= 400 && httpStatus < 500 {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	return &portal_proto.LoginResponse{Token: token}, nil
}

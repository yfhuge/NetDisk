// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package server

import (
	"context"

	"filestore-server/service/user/internal/logic"
	"filestore-server/service/user/internal/svc"
	"filestore-server/service/user/user"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) Signup(ctx context.Context, in *user.SignupRequest) (*user.SignupResponse, error) {
	l := logic.NewSignupLogic(ctx, s.svcCtx)
	return l.Signup(in)
}

func (s *UserServer) Signin(ctx context.Context, in *user.SigninRequest) (*user.SigninResponse, error) {
	l := logic.NewSigninLogic(ctx, s.svcCtx)
	return l.Signin(in)
}

func (s *UserServer) UserInfo(ctx context.Context, in *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	l := logic.NewUserInfoLogic(ctx, s.svcCtx)
	return l.UserInfo(in)
}
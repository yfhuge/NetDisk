package logic

import (
	"context"
	"filestore-server/model"

	"filestore-server/service/user/internal/svc"
	"filestore-server/service/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	username := in.Name

	// 查询用户信息
	userInfo, err := model.GetUserInfo(username)
	resp := new(user.UserInfoResponse)
	if err != nil {
		return nil, err
	}

	return &user.UserInfoResponse{
		Name:         userInfo.UserName,
		Email:        userInfo.Email,
		Phone:        userInfo.Phone,
		SignupAt:     userInfo.SignupAt,
		LastActiveAt: userInfo.LastActiveAt,
		Status:       userInfo.Status,
	}, nil
}

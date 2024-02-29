package logic

import (
	"context"
	"errors"
	"filestore-server/response"
	util "filestore-server/utils"

	"filestore-server/model"
	"filestore-server/service/user/internal/svc"
	"filestore-server/service/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SignupLogic) Signup(in *user.SignupRequest) (*user.SignupResponse, error) {
	// 获取数据
	username := in.Name
	passwd := in.Password
	encPwd := util.Sha1([]byte(passwd + l.svcCtx.Config.Salt))
	// 将数据插入到数据库中
	ok := model.UserSignup(username, encPwd)
	if !ok {
		response.Fail(c)
		return nil, errors.New("用户名已经存在")
	}
	return &user.SignupResponse{
		Code:    0,
		Message: "成功",
	}, nil
}

package logic

import (
	"context"
	"filestore-server/model"
	util "filestore-server/utils"

	"filestore-server/service/user/internal/svc"
	"filestore-server/service/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SigninLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSigninLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SigninLogic {
	return &SigninLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SigninLogic) Signin(in *user.SigninRequest) (*user.SigninResponse, error) {
	username := in.Name
	passwd := in.Password
	encPwd := util.Sha1([]byte(passwd + l.svcCtx.Config.Salt))
	resp := new(user.SigninResponse)

	// 校验用户名及密码
	ok := model.UserSignIn(username, encPwd)
	if !ok {
		resp.Code = 1
		resp.Message = "登录失败"
		return resp, nil
	}
	resp.Code = 0
	resp.Message = "success"
	return resp, nil
}

package user

import (
	"filestore-server/model"
	"filestore-server/response"
	util "filestore-server/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const PWD_SALT = "yf&huge"

type UserApi struct {
}

// UserSignUpPage 用户注册页面
func (u *UserApi) UserSignUpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

// UserSignInPage 用户登录页面
func (u *UserApi) UserSignInPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", gin.H{})
}

// UserSignUp 用户注册
func (u *UserApi) UserSignUp(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(username) < 3 || len(password) < 5 {
		response.FailWithMessage(c, "Invalid parameter!")
		return
	}

	// 将用户密码+密钥计算sha1
	encPwd := util.Sha1([]byte(password + PWD_SALT))
	// 进行落库
	ok := model.UserSignup(username, encPwd)
	if !ok {
		response.Fail(c)
		return
	}
	response.Success(c)
}

// UserSignIn 用户登录
func (u *UserApi) UserSignIn(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	encPwd := util.Sha1([]byte(password + PWD_SALT))

	// 1. 校验用户名及密码
	ok := model.UserSignIn(username, encPwd)
	if !ok {
		response.Fail(c)
		return
	}

	// 2. 生成访问凭证
	token, err := util.GetToken(username)
	if err != nil {
		response.Fail(c)
		return
	}
	// 更新用户token记录进行落库
	ok = model.UpdateToken(username, token)
	if !ok {
		response.Fail(c)
		return
	}

	// 3. 登录成功重定向到首页
	location := fmt.Sprintf("http://%s/static/view/home.html", c.Request.Host)
	resp := model.UserInfo{
		Location: location,
		UserName: username,
		Token:    token,
	}

	response.SuccessWithDetailed(c, "success", resp)
}

func (u *UserApi) UserInfo(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Query("username")

	// 2. 查询用户信息
	userInfo, err := model.GetUserInfo(username)
	if err != nil {
		c.Status(http.StatusForbidden)
		response.Fail(c)
		return
	}
	// 3. 返回用户详情
	response.SuccessWithDetailed(c, "ok", userInfo)
}

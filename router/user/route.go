package user

import (
	"filestore-server/api"
	"filestore-server/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (u *UserRouter) InitUserRouter(r *gin.Engine) {
	// 用户相关接口 登录、注册不需要auth拦截器验证
	group := r.Group("/user")
	apiGroup := api.ApiGroupApp.UserApiGroup

	// 注册
	group.GET("/signup", apiGroup.UserSignUpPage)
	group.POST("/signup", apiGroup.UserSignUp)

	// 登录
	group.GET("/signin", apiGroup.UserSignInPage)
	group.POST("/signin", apiGroup.UserSignIn)

	// 加入拦截器验证，后面接口都会进行auth中间件校验
	r.Use(middleware.Auth)
	group.POST("/info", apiGroup.UserInfo)
}

package initiallize

import (
	"filestore-server/middleware"
	"filestore-server/router"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	// 设置静态资源
	r.Static("/static", "./static")
	// 设置加载模板
	r.LoadHTMLGlob("static/view/**")

	r.Use(middleware.Cors)
	// 用户路由管理
	userRouter := router.RouterGroupApp.UserRouterGroup
	userRouter.InitUserRouter(r)

	// 文件路由管理
	fileStoreRouter := router.RouterGroupApp.FileStoreRouterGroup
	fileStoreRouter.InitFileStoreRouter(r)

	// 分块上传路由管理
	mpuploadRouter := router.RouterGroupApp.MpUploadRouterGroup
	mpuploadRouter.InitMpUploadRouter(r)

	// 注册pprof相关路由
	pprof.Register(r)

	return r
}

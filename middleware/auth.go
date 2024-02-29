package middleware

import (
	"filestore-server/response"
	util "filestore-server/utils"
	"github.com/gin-gonic/gin"
)

// Auth token拦截器验证
func Auth(c *gin.Context) {
	username := c.Query("username")
	token := c.Query("token")
	// 验证token是否有效
	if len(username) < 3 || !util.IsTokenValid(token) {
		// 后续的流程不在执行
		c.Abort()
		response.FailWithMessage(c, "token无效")
		return
	}
	// 流程到下一个流程
	c.Next()
}

package mpupload

import (
	"filestore-server/api"
	"github.com/gin-gonic/gin"
)

type MpUploadRouter struct {
}

func (m *MpUploadRouter) InitMpUploadRouter(r *gin.Engine) {
	group := r.Group("/file/mpupload")
	apiGroup := api.ApiGroupApp.MpUploadApiGroup
	group.POST("/init", apiGroup.InitMultipartUpload)
	group.POST("/uppart", apiGroup.UploadPart)
	group.POST("complete", apiGroup.CompleteUpload)
}

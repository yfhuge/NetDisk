package filestore

import (
	"filestore-server/api"
	"github.com/gin-gonic/gin"
)

type FileStoreRouter struct {
}

func (f *FileStoreRouter) InitFileStoreRouter(r *gin.Engine) {
	group := r.Group("/file")
	apiGroup := api.ApiGroupApp.FileStoreApiGroup
	// 文件存取接口
	group.GET("/upload", apiGroup.FileUploadPage)
	group.POST("/upload", apiGroup.FileUpload)
	group.POST("/meta", apiGroup.GetFileMeta)
	group.POST("/query", apiGroup.GetUserFileMetas)
	group.POST("/download", apiGroup.FileDownload)
	group.POST("/update", apiGroup.FileMetaUpdate)
	group.DELETE("/delete", apiGroup.FileDelete)
	group.POST("/downloadurl", apiGroup.DownloadURL)

	// 秒传接口
	group.POST("/fastupload", apiGroup.TryFastUpload)
}

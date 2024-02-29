package filestore

import (
	"encoding/json"
	"filestore-server/global"
	"filestore-server/model"
	"filestore-server/mq"
	"filestore-server/response"
	"filestore-server/store/oss"
	util "filestore-server/utils"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type FileStoreApi struct {
}

// FileUploadPage 文件上传页面
func (f *FileStoreApi) FileUploadPage(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{})
}

// FileUpload 文件上传
func (f *FileStoreApi) FileUpload(c *gin.Context) {
	// 接收文件流存储在本地
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(c, "Failed to get date, err:"+err.Error())
		return
	}
	src, err := file.Open()
	if err != nil {
		response.FailWithMessage(c, "upload file failed, err:"+err.Error())
		return
	}
	defer src.Close()
	fileMeta := model.FileMeta{
		FileName: file.Filename,
		Location: "./data/" + file.Filename,
		UploadAt: time.Now().Format("2000-01-02 12:00:00"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		response.FailWithMessage(c, "upload file failed, err:"+err.Error())
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, src)
	if err != nil {
		response.FailWithMessage(c, "upload file failed, err:"+err.Error())
		return
	}

	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)

	newFile.Seek(0, 0)
	// 同时将文件写入到oss
	ossPath := "oss/" + fileMeta.FileSha1
	//err = oss.Bucket().PutObject(ossPath, newFile)
	//if err != nil {
	//	log.Error("upload to oss failed, err:" + err.Error())
	//	response.FailWithMessage(c, "upload failed!")
	//	return
	//}
	fileMeta.Location = ossPath

	// 加入消息队列异步转存到oss
	data, _ := json.Marshal(fileMeta)
	err = mq.Send(global.Conf.RBMQConf, data)
	if err != nil {
		response.Fail(c)
		return
	}

	// 更新文件元信息，写入到数据库
	if ok := model.UpdateFileMetaDB(fileMeta); !ok {
		response.FailWithMessage(c, "update to sql failed, err:"+err.Error())
		return
	}

	// 更新用户文件元信息
	username := c.Query("username")
	if ok := model.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize); !ok {
		response.FailWithMessage(c, "Upload Failed")
	} else {
		response.Success(c)
	}
}

// GetFileMeta 获取文件元信息
func (f *FileStoreApi) GetFileMeta(c *gin.Context) {
	fileHash := c.Query("filehash")
	fileMeta, err := model.GetFileMetaDB(fileHash)
	if err != nil {
		response.Fail(c)
		return
	}
	data, err := json.Marshal(fileMeta)
	if err != nil {
		response.Fail(c)
		return
	}
	response.SuccessWithDetailed(c, "success!", data)
}

// GetUserFileMetas 批量获取用户文件元信息
func (f *FileStoreApi) GetUserFileMetas(c *gin.Context) {
	username := c.Query("username")
	limit, _ := strconv.Atoi(c.PostForm("limit"))
	userFiles, err := model.QueryUserFileMetas(username, limit)
	if err != nil {
		response.Fail(c)
		return
	}
	response.SuccessWithDetailed(c, "ok", userFiles)
}

// FileDownload 文件下载
func (f *FileStoreApi) FileDownload(c *gin.Context) {
	fileHash := c.Query("filehash")
	fm, err := model.GetFileMetaDB(fileHash)
	if err != nil {
		response.Fail(c)
		return
	}

	file, err := os.Open(fm.Location)
	if err != nil {
		response.Fail(c)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		response.Fail(c)
		return
	}
	c.Header("Content-Type", "application/octect-stream")
	c.Header("content-disposition", "attachment;filename="+fm.FileName)
	response.SuccessWithDetailed(c, "success", data)
}

// FileMetaUpdate 更新文件元信息
func (f *FileStoreApi) FileMetaUpdate(c *gin.Context) {
	op := c.Query("op")
	fileHash := c.Query("filehash")
	newFileName := c.Query("filename")

	if op != "0" {
		c.Status(http.StatusForbidden)
		response.Fail(c)
		return
	}
	if c.Request.Method != "POST" {
		c.Status(http.StatusMethodNotAllowed)
		response.Fail(c)
		return
	}

	currFileMeta, err := model.GetFileMetaDB(fileHash)
	if err != nil {
		response.Fail(c)
		return
	}
	currFileMeta.FileName = newFileName
	model.UpdateFileMetaDB(currFileMeta)
	data, err := json.Marshal(currFileMeta)
	if err != nil {
		response.Fail(c)
		return
	}
	response.SuccessWithDetailed(c, "success", data)
}

// FileDelete 删除文件
func (f *FileStoreApi) FileDelete(c *gin.Context) {
	fileHash := c.Query("filehash")
	fileName := c.Query("filename")
	fm, err := model.GetFileMetaDB(fileHash)
	if err != nil {
		response.Fail(c)
		return
	}
	// 删除文件
	os.Remove(fm.Location)
	// 删除文件原信息
	model.DeleteUserFileMeta(fileHash, fileName)
	response.Success(c)
}

// TryFastUpload 尝试秒传
func (f *FileStoreApi) TryFastUpload(c *gin.Context) {
	// 1. 解析请求参数
	username := c.PostForm("username")
	filehash := c.PostForm("filehash")
	filename := c.PostForm("filename")
	filesize, _ := strconv.Atoi(c.PostForm("filesize"))

	// 2. 从文件表中查询相同的hash的文件记录
	fileMeta, err := model.GetFile(filehash)
	if err != nil {
		log.Fatal(err)
		response.Fail(c)
		return
	}

	// 3. 查不到记录则返回秒传失败
	if fileMeta == nil {
		response.FailWithMessage(c, "秒传失败，请访问普通上传接口")
		return
	}

	// 4. 上传过则将文件信息吸入用户文件表，返回成功
	ok := model.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if !ok {
		response.FailWithMessage(c, "秒传失败，请稍后重试")
		return
	}
	response.FailWithMessage(c, "秒传成功")
}

// DownloadURL 生成oss文件的下载地址
func (f *FileStoreApi) DownloadURL(c *gin.Context) {
	filehash := c.Query("filehash")
	// 从文件表中查找记录
	row, _ := model.GetFile(filehash)
	signedURL := oss.GetClient().DownloadURL(row.FileAddr.String)
	response.SuccessWithDetailed(c, "ok", signedURL)
}

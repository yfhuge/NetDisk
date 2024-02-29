package mpupload

import (
	RDB "filestore-server/dao/redis"
	"filestore-server/global"
	"filestore-server/model"
	"filestore-server/response"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type MpUploadApi struct {
}

// InitMultipartUpload 初始化分块上传接口
func (m *MpUploadApi) InitMultipartUpload(c *gin.Context) {
	// 1. 解析用户请求参数
	username := c.PostForm("username")
	filehash := c.PostForm("filehash")
	filesize, err := strconv.Atoi(c.PostForm("filesize"))
	if err != nil {
		response.FailWithMessage(c, "params invalid")
		return
	}

	// 2. 获得redis的一个连接
	redisConn := global.RDB.GetConn()
	defer redisConn.Close()

	// 3. 生成分块上传的初始化信息
	uploadInfo := model.MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	// 4. 将初始化信息写入到redis缓存
	redisConn.Do("MSET", "MP_"+uploadInfo.UploadID, uploadInfo)

	// 5. 将初始化信息返回给客户端
	response.SuccessWithDetailed(c, "ok", uploadInfo)
}

// UploadPart 上传文件分块
func (m *MpUploadApi) UploadPart(c *gin.Context) {
	// 1. 解析用户请求参数
	uploadID := c.PostForm("uploadid")
	chunkIndex := c.PostForm("index")

	// 2. 获得一个redis连接
	redisConn := RDB.GetRDBInstance().GetConn()
	defer redisConn.Close()

	// 3. 获得文件句柄，用于存储分块内容
	fpath := "./data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		response.FailWithMessage(c, "Upload part failed")
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		if err != nil {
			break
		}
		fd.Write(buf[:n])
	}

	// 4. 更新redis缓存状态
	redisConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 5. 返回处理结果到客户端
	response.SuccessWithMessage(c, "ok")
}

// CompleteUpload 通知上传合并接口
func (m *MpUploadApi) CompleteUpload(c *gin.Context) {
	// 1. 解析请求参数
	uploadID := c.PostForm("uploadid")
	//username := c.PostForm("username")
	//filehash := c.PostForm("filehash")
	//filesize, _ := strconv.Atoi(c.PostForm("filesize"))
	//filename := c.PostForm("filename")

	// 2. 获得一个redis连接
	redisConn := RDB.GetRDBInstance().GetConn()
	defer redisConn.Close()

	// 3. 通过uploadid 查询redis并判断是否所有分开上传完成
	data, err := redis.Values(redisConn.Do("HGETALL", "MP_"+uploadID), nil)
	if err != nil {
		response.FailWithMessage(c, "compete upload failed")
		return
	}
	totalCount, chunkCount := 0, 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		response.FailWithMessage(c, "invalid request")
		return
	}

	// 4. 合并分块

	// 5. 更新唯一文件表及用户文件表

	// 6. 响应处理结果
	response.SuccessWithMessage(c, "ok")
}

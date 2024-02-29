package oss

import (
	"encoding/json"
	"filestore-server/global"
	"filestore-server/model"
	"filestore-server/mq"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	once   sync.Once
	client *OssClient
)

type OssClient struct {
	client *oss.Client
	bucket *oss.Bucket
}

func GetClient() *OssClient {
	once.Do(func() {
		client = &OssClient{}
	})
	return client
}

func (o *OssClient) Init() {
	var err error
	o.client, err = oss.New(global.Conf.OSSConf.EndPoint, global.Conf.OSSConf.AccessKeyID, global.Conf.OSSConf.AccessKeySecret)
	if err != nil {
		panic(err)
	}

	o.bucket, err = o.client.Bucket(global.Conf.OSSConf.Bucket)
	if err != nil {
		panic(err)
	}
}

func (o *OssClient) UploadFile() {
	chanMsg := make(chan []byte, 100)

	go mq.Receive(global.Conf.RBMQConf, chanMsg)

	for {
		select {
		case msg := <-chanMsg:
			fMeta := &model.FileMeta{}
			err := json.Unmarshal(msg, fMeta)
			if err != nil {
				log.Error("oss client upload file: [%s]", string(msg))
			} else {
				o.uploadToOss(fMeta)
			}
		}
	}
}

func (o *OssClient) uploadToOss(fMeta *model.FileMeta) {
	err := o.bucket.PutObjectFromFile("oss/"+fMeta.FileSha1, fMeta.Location)
	if err != nil {
		log.Error(err)
	}
}

func (o *OssClient) DownloadURL(objName string) string {
	signedURL, err := o.bucket.SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return signedURL
}

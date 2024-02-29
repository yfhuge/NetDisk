package main

import (
	"filestore-server/global"
	"filestore-server/initiallize"
	"filestore-server/store/oss"
)

func main() {
	initiallize.Viper()
	global.DB.Init(global.Conf.DBConf)
	global.RDB.Init(global.Conf.RDBConf)
	defer global.RDB.GetConn().Close()
	oss.GetClient().Init()
	go oss.GetClient().UploadFile()
	r := initiallize.Router()
	addr := global.Conf.System.Addr
	err := r.Run(addr)

	if err != nil {
		panic(err)
	}
}

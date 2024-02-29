package global

import (
	"filestore-server/config"
	"filestore-server/dao/db"
	"filestore-server/dao/redis"
)

var (
	Conf = config.GetConfig()
	DB   = db.GetDBInstance()
	RDB  = redis.GetRDBInstance()
)

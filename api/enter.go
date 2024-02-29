package api

import (
	"filestore-server/api/filestore"
	"filestore-server/api/mpupload"
	"filestore-server/api/user"
)

type ApiGroup struct {
	UserApiGroup      user.ApiGroup
	FileStoreApiGroup filestore.ApiGroup
	MpUploadApiGroup  mpupload.ApiGroup
}

var ApiGroupApp = new(ApiGroup)

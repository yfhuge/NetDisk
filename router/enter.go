package router

import (
	"filestore-server/router/filestore"
	"filestore-server/router/mpupload"
	"filestore-server/router/user"
)

type RouteGroup struct {
	UserRouterGroup      user.UserRouter
	FileStoreRouterGroup filestore.FileStoreRouter
	MpUploadRouterGroup  mpupload.MpUploadRouter
}

var RouterGroupApp = new(RouteGroup)

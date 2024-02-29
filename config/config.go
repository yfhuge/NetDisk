package config

import (
	"sync"
)

var (
	once sync.Once
	conf *Server
)

func GetConfig() *Server {
	once.Do(func() {
		conf = &Server{}
	})
	return conf
}

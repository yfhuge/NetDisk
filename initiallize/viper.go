package initiallize

import (
	"filestore-server/global"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Viper() {
	v := viper.New()
	v.SetConfigFile("config.yaml")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(global.Conf)
	if err != nil {
		panic(err)
	}
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Infoln("配置文件被修改了")
		if err := v.Unmarshal(global.Conf); err != nil {
			log.Error("viper unmarshal failed, err:" + err.Error())
		}
	})
}

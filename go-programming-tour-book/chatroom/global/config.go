package global

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var SensitiveWords []string
var MsgChanSize int

func initConfig() {
	// 配置文件名
	viper.SetConfigName("chatroom")
	// 配置文件路径
	viper.AddConfigPath(RootDir + "/config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	SensitiveWords = viper.GetStringSlice("sensitive")
	MsgChanSize = viper.GetInt("msg-chan-size")

	log.Printf("敏感词汇: %v\n", SensitiveWords)

	// 监听文件发送改变
	viper.WatchConfig()
	// 配置文件发生改变后的处理
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 重新读取配置文件更新SensitiveWords
		viper.ReadInConfig()
		SensitiveWords = viper.GetStringSlice("sensitive")
	})
}

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/model"
	"github.com/go-programming-tour-book/blog-service/internal/routers"
	"github.com/go-programming-tour-book/blog-service/pkg/logger"
	"github.com/go-programming-tour-book/blog-service/pkg/setting"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"time"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger  err: %v", err)
	}
}

// @title 博客系统
// @version 1.0
// @description Go 语言编程之旅：一起用 Go 做项目
// @termsOfService https://github.com/go-programming-tour-book
func main() {
	// 测试
	//e := gin.Default()
	//e.GET("/ping", func(context *gin.Context) {
	//	context.JSON(200, gin.H{"message1": "pong1", "message2": "pong2"})
	//})
	//e.Run()

	gin.SetMode(global.ServerSetting.RunMode) // 运行模式
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort, // 端口号
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	// 测试日志输出
	// global.Logger.Infof("%s: go-programming-tour-book/%s", "eddycjy", "blog-service")
	s.ListenAndServe()

}

// setupSetting 将从配置文件中读取的数据进行保存
func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	// 读取配置并保存至全局遍历中
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}

	global.JWTSetting.Expire *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}

// setupDBEngine 初始化数据库连接
func setupDBEngine() error {
	var err error
	// 注意：这里使用=而不是:=是为了直接将赋值到全局变量global.DBEngine
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	// 使用了 lumberjack 作为日志库的 io.Writer
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt, // 输出的目标文件
		MaxSize:   600,                                                                                                // 日志文件最大600M
		MaxAge:    10,                                                                                                 // 日志文件最大生存周期为 10 天
		LocalTime: true,                                                                                               // 使用本地时间
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}

package routers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-programming-tour-book/blog-service/docs" // 初始化doc包
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/middleware"
	"github.com/go-programming-tour-book/blog-service/internal/routers/api"
	v1 "github.com/go-programming-tour-book/blog-service/internal/routers/api/v1"
	"github.com/go-programming-tour-book/blog-service/pkg/limiter"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"time"
)

// 对/auth接口进行限流保护
var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.BucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

// NewRouter 路由对象(本质就是engine)
func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}
	r.Use(middleware.RateLimiter(methodLimiters))                                // 限流
	r.Use(middleware.ContextTimeout(global.ServerSetting.DefaultContextTimeout)) // 超时处理
	r.Use(middleware.Translations())                                             // 翻译
	r.Use(middleware.JWT())                                                      // JWT鉴权

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // 通过 http://ip:port/swagger/index.html 访问接口文档

	tag := v1.NewTag()
	article := v1.NewArticle()
	upload := api.NewUpload()
	r.GET("/auth", api.GetAuth)
	r.POST("/upload/file", upload.UploadFile) // 上传文件
	// http.Handle("/static/*filepath", http.StripPrefix("/static", http.FileServer(http.Dir("storage/uploads/"))))
	// 访问/static/xxxx.jpg 等于是访问到的 storage/uploads/xxxx.jpg路径
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath)) // 设置静态资源能够被访问，其实就是一个处理对应path的handler
	apiv1 := r.Group("api/v1")
	{
		// 绑定 path 和 对应的处理方法
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id/state", article.Update)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.GET("/articles", article.List)
	}

	return r
}

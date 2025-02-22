package routers

import (
	"chainqa_offchain_demo/controller"
	"chainqa_offchain_demo/setting"

	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func SetupRouter() *gin.Engine {

	if setting.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
		//如果setting.Conf.Release的值为true（即处于发布模式），则将gin引擎的工作模式设置为发布模式。这样做可以确保在生产环境中优化应用程序的性能。
	}
	r := gin.Default()
	// 使用CORS中间件
	r.Use(CORSMiddleware())

	// 前端模板
	r.LoadHTMLGlob("dist/index.html") // 加载HTML模板
	r.Static("/assets", "dist/assets")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	// 后端路由
	// api
	apiGroup := r.Group("/api")
	{
		helloGroup := apiGroup.Group("/hello")
		{
			helloGroup.GET("/hello", controller.HelloHandler)
		}

		uploadGroup := apiGroup.Group("/upload")
		{
			uploadGroup.POST("/getAesKey", controller.GetAesKeyHandler)
			uploadGroup.POST("/uploadFile", controller.UploadFileHandler)
		}

		downloadGroup := apiGroup.Group("/download")
		{
			downloadGroup.POST("/downloadIPFSFile", controller.DownloadIPFSFileHandler)
			downloadGroup.POST("/tryDecryptFile", controller.TryDecryptFileHandler)
		}

		queryGroup := apiGroup.Group("/query")
		{
			queryGroup.POST("/queryData", controller.QueryDataHandler)
		}

		logGroup := apiGroup.Group("/log")
		{
			logGroup.POST("/logByUid", controller.LogByUidHandler)
			logGroup.POST("/logByTimeRange", controller.LogByTimeRangeHandler)
		}
	}

	// 通用路由处理，捕获所有请求并返回index.html,后续前端路由处理
	r.NoRoute(func(c *gin.Context) {
		// 对于所有其他请求，返回index.html
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	// v1

	return r
}

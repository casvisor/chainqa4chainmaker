package routers

import (
	"tencent-chainmaker/controller"
	"tencent-chainmaker/setting"

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
	// api
	apiGroup := r.Group("/tencent-chainapi")
	{
		apiGroup.POST("/hello", controller.HelloHandler)
		apiGroup.POST("/exec", controller.ExecChain)
	}

	// v1

	return r
}

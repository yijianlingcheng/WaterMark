package api

import (
	"WaterMark/src/logs"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ServerStart api服务启动
func ServerStart() {

	// release模式
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	// 指定日志输出文件
	gin.DefaultWriter = logs.API.Writer()

	// router
	router := gin.Default()

	// 格式化日志
	router.Use(gin.LoggerWithFormatter(LogFormat))

	// 设置跨域
	router.Use(CORSMiddleware())

	// 加载路由
	loadRouters(router)

	// 监听地址并运行
	router.Run(viper.GetString("server.address"))
}

// CORSMiddleware 中间件处理跨域问题
//
//	@return gin.HandlerFunc
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// LogFormat 格式化日志输出格式
//
//	@param param
//	@return string
func LogFormat(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC3339Nano),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

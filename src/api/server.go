package api

import (
	"WaterMark/src/log"
	"WaterMark/src/tool"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServerStart api服务启动
func ServerStart() {

	log.InfoLogger.Println("Server Start begin")

	// release模式
	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// 设置跨域
	router.Use(CORSMiddleware())

	// 获取快门次数
	router.POST("/server/getShutterByFile", func(c *gin.Context) {

		log.InfoLogger.Println("req Server api /server/getShutterByFile")

		imgPath := c.PostForm("shutterimg")
		r := GetShutter(imgPath, true)

		log.InfoLogger.Println("res Server api /server/getShutterByFile:" + tool.ExifToJson(r))

		c.JSON(http.StatusOK, r)
	})

	// 获取图片预览
	router.GET("/server/getImagePreview", func(c *gin.Context) {
		imgPath := c.Query("imgagePath")
		if runtime.GOOS == "windows" {
			imgPath = strings.ReplaceAll(imgPath, "\\", "/")
		}
		file, _ := os.ReadFile(imgPath) //把要显示的图片读取到变量中
		c.Writer.WriteString(string(file))
	})

	port := "11079"

	log.InfoLogger.Println("Server listen port:" + port)
	log.InfoLogger.Println("Server Start success")

	router.Run(":" + port)
}

// CORSMiddleware 中间件处理跨域问题
//
//	@return gin.HandlerFunc
func CORSMiddleware() gin.HandlerFunc {

	log.InfoLogger.Println("Server set CORS policy")

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

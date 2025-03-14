package api

import (
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServerStart api服务启动
func ServerStart() {

	// release模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// 设置跨域
	router.Use(CORSMiddleware())

	// 获取快门次数
	router.POST("/server/getShutterByFile", func(c *gin.Context) {
		imgPath := c.PostForm("shutterimg")
		r := GetShutter(imgPath, true)
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

	// 获取图片模板列表类型
	router.GET("/server/getTplListType", func(c *gin.Context) {
		r := getTplListType()
		c.JSON(http.StatusOK, r)
	})

	// 获取图片水印预览
	router.POST("/server/getImageWaterMarkPreview", func(c *gin.Context) {
		imgPath := c.PostForm("imgagePath")
		tid := c.DefaultPostForm("tid", "1")
		flag := c.DefaultPostForm("flag", "false")
		color := c.DefaultPostForm("color", "255,255,255,255")
		if runtime.GOOS == "windows" {
			imgPath = strings.ReplaceAll(imgPath, "\\", "/")
		}
		r := getImageWaterMarkPreview(tid, imgPath, color, flag == "true")
		c.JSON(http.StatusOK, r)
	})

	port := "11079"
	router.Run(":" + port)
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

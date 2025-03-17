package controller

import (
	"WaterMark/src/cmd"
	"WaterMark/src/exif"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

// getShutterByFile 根据文件路径获取对应的快门信息
func GetShutterByFile(ctx *gin.Context) {
	c := Container(ctx)

	imgPath := c.PostForm("shutterimg")
	if runtime.GOOS == "windows" {
		imgPath = strings.ReplaceAll(imgPath, "\\", "/")
	}
	m, _ := cmd.CacheLoadExifTool(imgPath)
	r := exif.Getshutter(m)

	c.JSON(r)
}

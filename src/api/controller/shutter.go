package controller

import (
	"WaterMark/src/cmd"
	"WaterMark/src/exif"

	"github.com/gin-gonic/gin"
)

// getShutterByFile 根据文件路径获取对应的快门信息
func GetShutterByFile(ctx *gin.Context) {
	c := Container(ctx)

	imgPath := c.PostForm("shutterimg")
	m, _ := cmd.CacheLoadExifTool(imgPath)
	r := exif.Getshutter(m)

	c.JSON(r)
}

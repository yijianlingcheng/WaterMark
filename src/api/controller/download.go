package controller

import (
	"WaterMark/src/tool"
	"strings"

	"github.com/gin-gonic/gin"
)

// DownloadFile 下载文件
//
//	@param ctx
func DownloadFile(ctx *gin.Context) {
	c := Container(ctx)

	sourcePath := c.PostForm("source")
	previewPath := c.PostForm("preview")

	sourcePath = tool.ReplaceDir(sourcePath)
	previewPath = tool.ReplaceDir(previewPath)

	t := strings.Split(sourcePath, ".")
	t[len(t)-2] = t[len(t)-2] + "_watermark"

	newSourcePath := strings.Join(t, ".")
	tool.CopyFile(previewPath, newSourcePath, 4*1024)

	c.JSON(gin.H{"path": newSourcePath})
}

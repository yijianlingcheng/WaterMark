package controller

import (
	"WaterMark/src/cmd"
	"WaterMark/src/images"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetImagePreview 根据文件路径将对应的文件返回
//
//	@param ctx
func GetImagePreview(ctx *gin.Context) {
	c := Container(ctx)

	imgPath := c.Query("imagePath")
	if runtime.GOOS == "windows" {
		imgPath = strings.ReplaceAll(imgPath, "\\", "/")
	}
	file, _ := os.ReadFile(imgPath) //把要显示的图片读取到变量中
	c.ctx.Writer.WriteString(string(file))
}

// GetImageWaterMarkPreview
//
//	@param ctx
func GetImageWaterMarkPreview(ctx *gin.Context) {
	c := Container(ctx)
	imgPath := c.PostForm("imagePath")
	tid := c.DefaultPostForm("tid", "1")
	flag := c.DefaultPostForm("flag", "false")
	color := c.DefaultPostForm("color", "255,255,255,255")
	if runtime.GOOS == "windows" {
		imgPath = strings.ReplaceAll(imgPath, "\\", "/")
	}

	onlyBottomBorder := flag == "true"
	e := images.NewExternal()
	e.WithBoderColor(color).WithOnlyBottomFlag(onlyBottomBorder).WithPath(imgPath).WithTid(tid)
	r := images.CreatePreviewWaterMark(e)
	r["SaveImgPath"] = cmd.GetPwdPath(strings.TrimLeft(r["SaveImgPath"], "."))

	c.JSON(r)
}

// ImagePreviewSmall 将原图地址转换为预览小图地址并返回图片
//
//	@param ctx
func ImagePreviewSmall(ctx *gin.Context) {
	c := Container(ctx)

	imgPath := c.Query("imagePath")
	if runtime.GOOS == "windows" {
		imgPath = strings.ReplaceAll(imgPath, "\\", "/")
	}

	e := images.NewExternal().WithSmallPreviewPath(imgPath)
	file, _ := os.ReadFile(e.SavePath)
	c.ctx.Writer.WriteString(string(file))
}

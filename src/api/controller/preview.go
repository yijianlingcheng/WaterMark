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
	onlyBottomBorder := flag == "true"

	borderColor := c.DefaultPostForm("borderColor", "255,255,255,255")
	words := c.DefaultPostForm("words", "")
	model := c.DefaultPostForm("model", "")
	lensModel := c.DefaultPostForm("lensModel", "")
	firstWordsColor := c.DefaultPostForm("firstWordsColor", "")
	secondBorderColor := c.DefaultPostForm("secondBorderColor", "")

	if runtime.GOOS == "windows" {
		imgPath = strings.ReplaceAll(imgPath, "\\", "/")
	}

	e := images.NewExternal()
	e.WithBoderColor(borderColor).WithOnlyBottomFlag(onlyBottomBorder).WithPath(imgPath).WithTid(tid)
	if len(words) > 0 {
		e.WithWords(words)
	}
	if len(model) > 0 {
		e.WithModel(model)
	}
	if len(lensModel) > 0 {
		e.WithLensModel(lensModel)
	}
	if len(firstWordsColor) > 0 {
		e.WithFirstWordsColor(firstWordsColor)
	}
	if len(secondBorderColor) > 0 {
		e.WithSecondBorderColor(secondBorderColor)
	}
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

package controller

import (
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
	if runtime.GOOS == "windows" {
		imgPath = strings.ReplaceAll(imgPath, "\\", "/")
	}
	op := c.PostForm("type")

	// 操作类型,load代表加载已经生成好的水印数据
	if op == "load" {
		r, err := imageInfoCacheGet(imgPath)
		if err == nil {
			c.JSON(r)
			return
		}
	}
	// 模板id
	tid := c.DefaultPostForm("tid", "1")

	// 边框类型
	flag := c.DefaultPostForm("flag", "false")
	onlyBottomBorder := flag == "true"

	borderColor := c.DefaultPostForm("borderColor", "255,255,255,255")
	words := c.DefaultPostForm("words", "")
	model := c.DefaultPostForm("model", "")
	lensModel := c.DefaultPostForm("lensModel", "")
	firstWordsColor := c.DefaultPostForm("firstWordsColor", "")
	secondBorderColor := c.DefaultPostForm("secondBorderColor", "")

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
	r = mergeParams(ctx, r)
	// 新增缓存
	imageInfoCacheSet(imgPath, r)
	c.JSON(r)
}

// mergeParams 合并参数
//
//	@param ctx
//	@param r
//	@return map
func mergeParams(ctx *gin.Context, r map[string]string) map[string]string {
	for i, v := range ctx.Request.Form {
		r[i] = strings.Join(v, "")
	}
	return r
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

// ChangeImagePath 将原始文件路径转换成水印预览图路径
//
//	@param ctx
func ChangeImagePath(ctx *gin.Context) {
	c := Container(ctx)

	imgPath := c.Query("imagePath")
	if runtime.GOOS == "windows" {
		imgPath = strings.ReplaceAll(imgPath, "\\", "/")
	}
	e := images.NewExternal().WithPath(imgPath)
	c.JSON(gin.H{
		"path": e.SavePath,
	})
}

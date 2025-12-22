package controller

import (
	"image/jpeg"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/engine"
	"WaterMark/engine/frame"
	"WaterMark/internal"
	"WaterMark/layout"
	"WaterMark/message"
	"WaterMark/pkg"
)

// 检查参数.
func baseCheckShowPhotoFrame(ctx *gin.Context) map[string]any {
	file := ctx.PostForm(paramQueryFile)
	if file == "" {
		return map[string]any{"error": pkg.NewErrors(pkg.REQUEST_PARAM_ERROR, paramFileIsEmpty)}
	}
	photoType := ctx.PostForm("type")
	if photoType == "" {
		return map[string]any{"error": pkg.NewErrors(pkg.REQUEST_PARAM_ERROR, "type参数为空")}
	}
	if photoType != "border" && photoType != "photo" {
		return map[string]any{"error": pkg.NewErrors(pkg.REQUEST_PARAM_ERROR, "type参数类型错误")}
	}
	if !internal.PathExists(file) {
		return map[string]any{"error": pkg.NewErrors(pkg.FILE_NOT_EXIST_ERROR, file+":"+paramFileIsNotExist)}
	}
	layoutStr := ctx.PostForm(paramQueryLayout)
	if layoutStr == "" {
		return map[string]any{"error": pkg.NewErrors(pkg.REQUEST_PARAM_ERROR, paramLayoutIsEmpty)}
	}
	layout, buildErr := buildFramePrams(layoutStr)
	if pkg.HasError(buildErr) {
		return map[string]any{"error": buildErr}
	}
	exifInfo, checkErr := getExifAndCheckPhotoLogoExist(file)
	if pkg.HasError(checkErr) {
		return map[string]any{"error": checkErr}
	}

	return map[string]any{
		"error":    pkg.NoError,
		"file":     file,
		"type":     photoType,
		"layout":   layout,
		"exifInfo": exifInfo,
	}
}

// @Summary 对指定照片生成边框图
// @Description 对指定照片生成边框水印图片,并且直接输出图片内容
// @Tags Frame
// @Produce image/jpeg
// @Param file formData string true "照片路径"
// @Param type formData string true "返回的图片类型,border:只返回边框图,photo:返回照片+边框合成图"
// @Param layout formData string true "布局信息,JSON字符串:必须包含frame_name字段"
// @Router /frame/showPhotoFrame [post]
// @Failure 400 {object} ErrorInfo "错误信息".
func ShowPhotoFrame(ctx *gin.Context) {
	checkResult := baseCheckShowPhotoFrame(ctx)
	if checkErr, ok := checkResult["error"].(pkg.EError); ok {
		if pkg.HasError(checkErr) {
			ctx.JSON(400, checkErr)

			return
		}
	}
	file, fileOk := checkResult["file"].(string)
	photoType, photoOk := checkResult["type"].(string)
	exifInfo, exifInfoOk := checkResult["exifInfo"].(exiftool.FileMetadata)
	layout, layoutOk := checkResult["layout"].(layout.FrameLayout)

	if !fileOk || !photoOk || !exifInfoOk || !layoutOk {
		ctx.JSON(400, pkg.InternalError)

		return
	}
	plug := frame.GetPlugin()
	// 获取指定照片的缩略图
	imageRGBA, frameErr := plug.CreateFrameImageRGBA(map[string]any{
		"sourceImageFile": file,
		"photoType":       photoType,
		"exif":            exifInfo,
		"Params":          layout,
	})
	if pkg.HasError(frameErr) {
		ctx.JSON(400, frameErr)

		return
	}
	// 生成更小的图片,加快前端访问,将jpg图片作为输出直接返回
	err := jpeg.Encode(ctx.Writer, photoFrameResize(imageRGBA), &jpeg.Options{Quality: 75})
	if err != nil {
		message.SendErrorMsg("ShowPhotoFrame 接口出现错误:" + err.Error())
	}
}

// @Summary 重新加载logo文件下的图片
// @Description 重新加载logo文件夹下的图片,此接口用于运行过程中在logo文件夹下面添加了文件
// @Tags Frame
// @Produce json
// @Router /frame/reloadLogoImages [post]
// @Success 200 {object} NoError "成功信息".
// @Failure 400 {object} ErrorInfo "错误信息".
func ReloadLogoImages(ctx *gin.Context) {
	plug := frame.GetPlugin()
	err := plug.ReloadLogoImages()
	if pkg.HasError(err) {
		ctx.JSON(400, err)

		return
	}
	ctx.JSON(200, pkg.NoError)
}

// @Summary 重新加载边框模板文件
// @Description 重新加载边框模板文件,此接口用于运行过程中调整或新增了边框布局文件
// @Tags Frame
// @Produce json
// @Router /frame/reloadFrameTemplate [post]
// @Success 200 {object} NoError "成功信息".
// @Failure 400 {object} ErrorInfo "错误信息".
func ReloadFrameTemplate(ctx *gin.Context) {
	plug := frame.GetPlugin()
	err := plug.ReloadFrameTemplate()
	if pkg.HasError(err) {
		ctx.JSON(400, err)

		return
	}
	ctx.JSON(200, pkg.NoError)
}

// @Summary 获取边框模板信息
// @Description 获取边框模板信息,此接口用于展示边框模板列表
// @Tags Frame
// @Produce json
// @Router /frame/getFrameTemplateInfo [post]
// @Success 200 {object} TemplatesInfo "成功信息".
// @Failure 400 {object} ErrorInfo "错误信息".
func GetFrameTemplateInfo(ctx *gin.Context) {
	plug := frame.GetPlugin()
	err := plug.ReloadFrameTemplate()
	if pkg.HasError(err) {
		ctx.JSON(400, err)

		return
	}
	frame.LoadOrCreateLayoutImage()
	ctx.JSON(200, TemplatesInfo{
		Code: pkg.NO_ERROR,
		List: frame.GetTemplateInfo(),
	})
}

// @Summary 将指定图片导入到内存中
// @Description 将指定图片导入到内存中,提供给生成照片边框的逻辑使用
// @Tags Frame
// @Produce json
// @Param file formData string true "照片路径支持多个路径;多个路径,隔开"
// @Router /frame/importPhotoFiles [post]
// @Success 200 {object} ImportInfo "成功信息".
// @Failure 400 {object} ErrorInfo "错误信息".
func ImportPhotoFiles(ctx *gin.Context) {
	file := ctx.PostForm(paramQueryFile)
	if file == "" {
		ctx.JSON(400, requestParamError(paramFileIsEmpty))

		return
	}
	paths := strings.Split(file, ",")
	for _, path := range paths {
		if !internal.PathExists(path) {
			ctx.JSON(400, requestResoureNotExistError(path, paramFileIsNotExist))

			return
		}
	}
	limit := ctx.PostForm("limit")
	limitNum, _ := strconv.Atoi(limit)

	exifInfos := make([]exiftool.FileMetadata, 0)
	paths = uniqueStr(paths)
	if len(paths) > limitNum {
		paths = paths[:limitNum]
	}

	for _, path := range paths {
		exifInfo, err := engine.CacheGetImageExif(path)
		if pkg.HasError(err) {
			ctx.JSON(400, err)

			return
		}
		exifInfos = append(exifInfos, exifInfo)
	}

	plug := frame.GetPlugin()
	plug.ImportImageFiles(paths, exifInfos)

	ctx.JSON(200, ImportInfo{
		Code:  pkg.NO_ERROR,
		Files: paths,
	})
}

// @Summary 获取照片的exif与边框信息
// @Description 获取照片的exif与边框信息
// @Tags Frame
// @Produce json
// @Param file formData string true "照片路径"
// @Param type formData string true "图片尺寸类型,border:边框尺寸,photo:原始图片尺寸"
// @Param layout formData string true "布局信息,JSON字符串:必须包含frame_name字段"
// @Router /frame/getExifAndBorderInfo [post]
// @Success 200 {object} ExifAndBorderInfo "成功信息"
// @Failure 400 {object} ErrorInfo "错误信息".
func GetPhotoExifAndBorderInfo(ctx *gin.Context) {
	checkResult := baseCheckShowPhotoFrame(ctx)
	if checkErr, ok := checkResult["error"].(pkg.EError); ok {
		if pkg.HasError(checkErr) {
			ctx.JSON(400, checkErr)

			return
		}
	}
	file, fileOk := checkResult["file"].(string)
	photoType, photoOk := checkResult["type"].(string)
	exifInfo, exifInfoOk := checkResult["exifInfo"].(exiftool.FileMetadata)
	layout, layoutOk := checkResult["layout"].(layout.FrameLayout)
	if !fileOk || !photoOk || !exifInfoOk || !layoutOk {
		ctx.JSON(400, pkg.InternalError)

		return
	}
	// 获取指定照片的缩略图
	info, err := frame.GetPlugin().GetFrameImageBorderInfo(map[string]any{
		"sourceImageFile": file,
		"photoType":       photoType,
		"exif":            exifInfo,
		"Params":          layout,
	})
	if pkg.HasError(err) {
		ctx.JSON(400, err)

		return
	}
	size, sizeOk := info["size"].(map[string]int)
	text, textOk := info["text"].([]string)
	if !sizeOk || !textOk {
		ctx.JSON(400, pkg.InternalError)

		return
	}
	ctx.JSON(200, ExifAndBorderInfo{
		Exif: exifInfoTranslatorApi(exifInfo),
		Size: frame.NewPhotoSize(size),
		Text: text,
	})
}

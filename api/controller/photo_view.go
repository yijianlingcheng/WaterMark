package controller

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"WaterMark/engine"
	"WaterMark/internal"
	"WaterMark/pkg"
)

// @Summary 获取照片exif信息
// @Description 返回指定照片的exif信息与摘要
// @Tags view
// @Param file formData string true "照片路径"
// @Produce json
// @Success 200 {object} ExifInfoSuccess "成功"
// @Failure 400 {object} ErrorInfo "错误信息"
// @Router /view/GetPhotosExifInfo [post].
func GetPhotosExifInfo(ctx *gin.Context) {
	path := ctx.PostForm(paramQueryFile)
	if path == "" {
		ctx.JSON(400, requestParamError(paramFileIsEmpty))

		return
	}
	if !internal.PathExists(path) {
		ctx.JSON(400, requestResoureNotExistError(path, paramFileIsNotExist))

		return
	}
	exifInfo, err := engine.CacheGetImageExif(path)
	if pkg.HasError(err) {
		ctx.JSON(400, err)
	} else {
		ctx.JSON(200, exifInfoTranslatorApi(exifInfo))
	}
}

// @Summary 展示照片
// @Description 展示指定照片,直接输出图片信息
// @Tags view
// @Produce image/jpeg
// @Param file query string true "照片路径"
// @Router /view/showImage [get].
// @Failure 400 {object} ErrorInfo "错误信息".
func ShowImage(ctx *gin.Context) {
	path := ctx.Query(paramQueryFile)
	if internal.IsWindows() {
		path = strings.ReplaceAll(path, "\\", "/")
	}
	file, err := os.Open(path)
	if err != nil {
		ctx.JSON(400, requestResoureNotExistError(path, paramFileIsNotExist))

		return
	}
	defer file.Close()

	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		ctx.JSON(400, requestResoureNotExistError(path, "file请求的文件返回失败:"+err.Error()))
	}
}

// 导出照片exif信息对应的csv文件,通过后端直接保存文件实现
// exif导出必须在成功调用GetPhotosExifInfo接口之后

// @Summary 导出指定照片的exif信息
// @Description 导出当前展示照片的exif信息,将其保存在指定路径中
// @Tags view
// @Param file formData string true "照片路径"
// @Param save formData string true "导出文件保存的路径"
// @Produce json
// @Success 200 {object} NoError "导出完成"
// @Failure 400 {object} ErrorInfo "错误信息".
// @Router /view/ExifInfoExportBySaveFile [post].
func ExifInfoExportBySaveFile(ctx *gin.Context) {
	path := ctx.PostForm(paramQueryFile)
	savePath := ctx.PostForm(paramQuerySave)

	if internal.IsWindows() {
		path = strings.ReplaceAll(path, "\\", "/")
		savePath = strings.ReplaceAll(savePath, "\\", "/")
	}

	// 删除同名文件
	if internal.PathExists(savePath) {
		os.Remove(savePath)
	}
	// 判断保存文件夹是否存在并创建
	d := filepath.Dir(savePath)
	if !internal.PathExists(d) {
		err := os.MkdirAll(d, os.ModePerm)
		if err != nil {
			ctx.JSON(400, pkg.NewErrors(pkg.FILE_NOT_OPEN_ERROR, d+":保存的文件的文件夹创建失败"))

			return
		}
	}

	exifInfo, _ := engine.CacheGetImageExif(path)
	csvData := exifInfoTranslatorCsv(exifInfo)

	csv := pkg.CreateCSV(filepath.Base(savePath), d+"/", true)
	csv.AddData(csvData)
	err := csv.Generate()

	ctx.JSON(200, err)
}

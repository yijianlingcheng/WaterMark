package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/internal"
	"WaterMark/pkg"
)

type exifInfoAbstr struct {
	equipment []string
	mode      []string
	params    []string
	focal     []string
	color     []string
}

// 请求参数错误.
func requestParamError(msg string) map[string]any {
	return map[string]any{
		"code":   pkg.REQUEST_PARAM_ERROR,
		"errmsg": msg,
		"file":   "",
	}
}

// 请求的资源不存在.
func requestResoureNotExistError(file, msg string) map[string]any {
	return map[string]any{
		"code":   pkg.FILE_NOT_EXIST_ERROR,
		"errmsg": msg,
		"file":   file,
	}
}

// 转换exif获取结果.
func exifInfoTranslatorApi(exifInfo exiftool.FileMetadata) ExifInfoSuccess {
	if internal.IsWindows() {
		exifInfo.File = strings.ReplaceAll(exifInfo.File, "\\", "/")
	}
	if exifInfo.Err != nil {
		return ExifInfoSuccess{
			Code:   pkg.EXIFTOOL_IMAGE_EXIF_ERROR,
			Errmsg: "获取exif信息失败:" + exifInfo.Err.Error(),
			File:   exifInfo.File,
		}
	}
	exifAbstr := exifInfoAbstrGenerate(exifInfo)
	result := ExifInfoSuccess{
		Code:          pkg.NO_ERROR,
		File:          exifInfo.File,
		Equipment:     strings.Join(exifAbstr.equipment, " "),
		CMode:         strings.Join(exifAbstr.mode, " "),
		Params:        strings.Join(exifAbstr.params, " "),
		Focal:         strings.Join(exifAbstr.focal, " "),
		Color:         strings.Join(exifAbstr.color, " "),
		Time:          pkg.AnyToString(exifInfo.Fields["DateTimeOriginal"]),
		Shutter:       pkg.AnyToString(exifInfo.Fields["ShutterCount"]),
		ImageWidth:    pkg.AnyToString(exifInfo.Fields["ImageWidth"]),
		ImageHeight:   pkg.AnyToString(exifInfo.Fields["ImageHeight"]),
		Make:          pkg.AnyToString(exifInfo.Fields["Make"]),
		Model:         pkg.AnyToString(exifInfo.Fields["Model"]),
		LensModel:     pkg.AnyToString(exifInfo.Fields["LensModel"]),
		FocalLength:   pkg.AnyToString(exifInfo.Fields["FocalLength"]),
		FNumber:       pkg.AnyToString(exifInfo.Fields["FNumber"]),
		ExposureTime:  pkg.AnyToString(exifInfo.Fields["ExposureTime"]),
		ISO:           pkg.AnyToString(exifInfo.Fields["ISO"]),
		FileName:      pkg.AnyToString(exifInfo.Fields["FileName"]),
		ImageSize:     pkg.AnyToString(exifInfo.Fields["ImageSize"]),
		ImageDataSize: pkg.AnyToString(exifInfo.Fields["ImageDataSize"]),
		Orientation:   pkg.AnyToString(exifInfo.Fields["Orientation"]),
	}

	return result
}

// 转换exif获取结果.
func exifInfoTranslatorCsv(exifInfo exiftool.FileMetadata) [][]string {
	if internal.IsWindows() {
		exifInfo.File = strings.ReplaceAll(exifInfo.File, "\\", "/")
	}
	if exifInfo.Err != nil {
		return [][]string{
			{"file", exifInfo.File},
			{"errmsg", "获取exif信息失败:" + exifInfo.Err.Error()},
		}
	}

	// 获取摘要信息
	exifAbstr := exifInfoAbstrGenerate(exifInfo)

	csvData := make([][]string, 0, 400)
	csvData = append(
		csvData,
		[]string{"文件", exifInfo.File},
		[]string{"器材", strings.Join(exifAbstr.equipment, " ")},
		[]string{"模式", strings.Join(exifAbstr.mode, " ")},
		[]string{"参数", strings.Join(exifAbstr.params, " ")},
		[]string{"焦距", strings.Join(exifAbstr.focal, " ")},
		[]string{"色彩", strings.Join(exifAbstr.color, " ")},
		[]string{"时间", pkg.AnyToString(exifInfo.Fields["DateTimeOriginal"])},
		[]string{"快门次数", pkg.AnyToString(exifInfo.Fields["ShutterCount"])},
	)

	for key, value := range exifInfo.Fields {
		csvData = append(csvData, []string{key, pkg.AnyToString(value)})
	}

	return csvData
}

// exif摘要信息生成.
func exifInfoAbstrGenerate(exifInfo exiftool.FileMetadata) exifInfoAbstr {
	// 器材摘要
	equipment := []string{
		// 相机
		pkg.AnyToString(exifInfo.Fields["Model"]),
		// 镜头
		pkg.AnyToString(exifInfo.Fields["LensModel"]),
	}

	// 模式摘要
	mode := []string{
		"曝光模式:" + pkg.AnyToString(exifInfo.Fields["ExposureProgram"]),
		"测光模式:" + pkg.AnyToString(exifInfo.Fields["MeteringMode"]),
	}

	// 曝光补偿
	if exifInfo.Fields["ExposureBiasValue"] == nil {
		exifInfo.Fields["ExposureBiasValue"] = "0"
	}
	mode = append(mode, "曝光补偿:"+pkg.AnyToString(exifInfo.Fields["ExposureBiasValue"]))

	// 参数
	params := []string{
		"光圈:" + pkg.AnyToString(exifInfo.Fields["FNumber"]),
		"快门:" + pkg.AnyToString(exifInfo.Fields["ExposureTime"]),
		"ISO:" + pkg.AnyToString(exifInfo.Fields["ISO"]),
	}

	// 焦段
	focal := []string{
		pkg.AnyToString(exifInfo.Fields["FocalLength35efl"]),
		"视角:" + pkg.AnyToString(exifInfo.Fields["FOV"]),
	}

	// 白平衡
	if pkg.AnyToString(exifInfo.Fields["LightSource"]) == "Unknown" {
		exifInfo.Fields["LightSource"] = "Auto"
	}
	// 色彩
	color := []string{
		"白平衡:" + pkg.AnyToString(exifInfo.Fields["LightSource"]),
		"色彩空间:" + pkg.AnyToString(exifInfo.Fields["ColorSpace"]),
	}

	return exifInfoAbstr{equipment: equipment, mode: mode, params: params, focal: focal, color: color}
}

// 设置SSE头信息.
func setSSEHeader(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
}

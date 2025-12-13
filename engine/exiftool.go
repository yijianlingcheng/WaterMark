package engine

import (
	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/engine/exift"
	"WaterMark/pkg"
)

var (
	// exiftool工具实例.
	et *exiftool.Exiftool

	// 工具初始化标致.
	exiftoolInitFlag = false
)

// 初始化exiftool工具.
func initExiftool() pkg.EError {
	var err pkg.EError
	et, err = exift.InitExiftool()
	if !pkg.IsOk(err) {
		return err
	}
	exiftoolInitFlag = true

	return pkg.NoError
}

// 关闭exiftool工具.
func closeExiftool() {
	if exiftoolInitFlag {
		et.Close()
	}
}

// 获取指定照片的exif信息,支持传递多张图片.
func GetPhotosExifInfos(f ...string) []exiftool.FileMetadata {
	return et.ExtractMetadata(f...)
}

// 获取指定图片的exif信息,只支持单张图片.
func GetPhotosExifInfo(path string) (exiftool.FileMetadata, pkg.EError) {
	exifInfos := GetPhotosExifInfos(path)
	if exifInfos[0].Err != nil {
		return exiftool.FileMetadata{}, pkg.NewErrors(
			pkg.ExiftoolImageError.Code,
			path+":获取exif信息失败:"+exifInfos[0].Err.Error(),
		)
	}

	return exifInfos[0], pkg.NoError
}

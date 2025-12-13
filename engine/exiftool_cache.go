package engine

import (
	"maps"
	"sync"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/pkg"
)

// exif信息缓存池.
var exiftoolCache sync.Map

// 带缓存的获取图片exif信息.
func CacheGetImageExif(path string) (exiftool.FileMetadata, pkg.EError) {
	// 计算文件路径的md5,暂时忽略掉文件路径不变,但是文件内容发生变化的情况
	md5 := pkg.GetStrMD5(path)

	// 返回缓存
	if cache, ok := exiftoolCache.Load(md5); ok {
		if v, vok := cache.(exiftool.FileMetadata); vok {
			exifData := copyExifData(v)

			return exifData, pkg.NoError
		}

		return exiftool.FileMetadata{}, pkg.ExiftoolCacheTypeError
	}

	// 调用exiftool工具获取全量exif信息
	exifInfos := GetPhotosExifInfos(path)
	if exifInfos[0].Err != nil {
		errmsg := path + ":获取exif信息失败:" + exifInfos[0].Err.Error()

		return exiftool.FileMetadata{}, pkg.NewErrors(pkg.ExiftoolImageError.Code, errmsg)
	}
	exiftoolCache.Store(md5, exifInfos[0])

	exifData := copyExifData(exifInfos[0])

	return exifData, pkg.NoError
}

// 拷贝.避免后续修改exif信息导致影响原始exif值.
func copyExifData(exifInfos exiftool.FileMetadata) exiftool.FileMetadata {
	fields := make(map[string]any)
	maps.Copy(fields, exifInfos.Fields)
	exifData := exiftool.FileMetadata{
		File:   exifInfos.File,
		Fields: fields,
		Err:    exifInfos.Err,
	}

	return exifData
}

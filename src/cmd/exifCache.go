package cmd

import (
	"WaterMark/src/exif"
	. "WaterMark/src/logs"
	"WaterMark/src/tool"
	"sync"
)

// exifCache exif缓存
var exifCache sync.Map

// exifCacheErrors
var exifCacheErrors sync.Map

// cacheLoadExifTool  加载图片并缓存
//
//	@param path
//	@return exif.Exif
//	@return error
func CacheLoadExifTool(path string) (exif.Exif, error) {
	// 计算md5
	md5 := tool.StrMD5(path)
	// 查找是否加载错误
	if err, ok := exifCacheErrors.Load(md5); ok {
		return exif.Exif{}, err.(error)
	}

	// 返回缓存
	if cache, ok := exifCache.Load(md5); ok {
		Info.Println("读取exif缓存成功:" + path)
		return cache.(exif.Exif), nil
	}

	exifInfo, err := RunExifTool(path)
	if err != nil {
		exifCacheErrors.Store(md5, err)
		Errors.Println("获取exif信息失败:" + path)
		return exif.Exif{}, err
	}
	exifCache.Store(md5, exifInfo)

	Info.Println("写入exif缓存成功:" + path)
	return exifInfo, nil
}

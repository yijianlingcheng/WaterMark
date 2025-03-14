package cmd

import (
	"WaterMark/src/exif"
	"WaterMark/src/log"
	"WaterMark/src/tool"
)

// exifCache
var exifCache *ExifCache

// ExifCache 本地的exiftool缓存对象,导入的exif信息全部存放到这个缓存中,加速数据访问
type ExifCache struct {

	// cache
	caches map[string]exif.Exif

	// errors
	errors map[string]error
}

// InitExifCache 初始化exif缓存
func InitExifCache() {
	caches := map[string]exif.Exif{}
	errors := map[string]error{}
	exifCache = &ExifCache{
		caches: caches,
		errors: errors,
	}
	log.InfoLogger.Println("初始化exif缓存成功")
}

// cacheLoadExifTool  加载图片并缓存
//
//	@param path
//	@return exif.Exif
//	@return error
func CacheLoadExifTool(path string) (exif.Exif, error) {
	// 计算md5
	md5 := tool.StrMD5(path)
	// 查找是否加载错误
	if err, ok := exifCache.errors[md5]; ok {
		return exif.Exif{}, err
	}

	// 返回缓存
	if cache, ok := exifCache.caches[md5]; ok {
		log.InfoLogger.Println("读取exif缓存成功:" + path)
		return cache, nil
	}

	exifInfo, err := RunExifTool(path)
	if err != nil {
		exifCache.errors[md5] = err
		return exif.Exif{}, err
	}
	exifCache.caches[md5] = exifInfo

	log.InfoLogger.Println("写入exif缓存成功:" + path)
	return exifInfo, nil
}

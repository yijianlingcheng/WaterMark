package images

import (
	"WaterMark/src/log"
	"WaterMark/src/tool"
	"image"
)

// imagesCache
var imagesCache *ImagesCache

// ImagesCache 本地的图片缓存对象,导入的图片全部存放到这个缓存中,加速数据访问
type ImagesCache struct {

	// cache
	caches map[string]image.Image

	// errors
	errors map[string]error
}

// InitImagesCache 初始化图片缓存
func InitImagesCache() {
	caches := map[string]image.Image{}
	errors := map[string]error{}
	imagesCache = &ImagesCache{
		caches: caches,
		errors: errors,
	}
	log.InfoLogger.Println("初始化图片缓存成功")
}

// cacheLoadImage LoadImage 加载图片并缓存
//
//	@param path
//	@return image.Image
//	@return error
func cacheLoadImage(path string) (image.Image, error) {
	// 计算md5
	md5 := tool.StrMD5(path)
	// 查找是否加载错误
	if err, ok := imagesCache.errors[md5]; ok {
		return nil, err
	}

	// 返回缓存
	if cache, ok := imagesCache.caches[md5]; ok {
		log.InfoLogger.Println("读取图片缓存成功:" + path)
		return cache, nil
	}

	image, err := loadImage(path)
	if err != nil {
		imagesCache.errors[md5] = err
		return nil, err
	}
	imagesCache.caches[md5] = image

	log.InfoLogger.Println("写入图片缓存成功:" + path)
	return image, nil
}

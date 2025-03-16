package images

import (
	"WaterMark/src/logs"
	"WaterMark/src/tool"
	"image"
	"sync"
)

// imagesCache 图片缓存
var imagesCache sync.Map

// imagesCacheErrors
var imagesCacheErrors sync.Map

// cacheLoadImage LoadImage 加载图片并缓存
//
//	@param path
//	@return image.Image
//	@return error
func cacheLoadImage(path string) (image.Image, error) {
	// 计算md5
	md5 := tool.StrMD5(path)
	// 查找是否加载错误
	if err, ok := imagesCacheErrors.Load(md5); ok {
		return nil, err.(error)
	}

	// 返回缓存
	if cache, ok := imagesCache.Load(md5); ok {
		logs.Info.Println("读取图片缓存成功:" + path)
		return cache.(image.Image), nil
	}

	image, err := loadImage(path)
	if err != nil {
		imagesCacheErrors.Store(md5, err)
		return nil, err
	}
	imagesCache.Store(md5, image)

	logs.Info.Println("写入图片缓存成功:" + path)
	return image, nil
}

package images

import (
	"WaterMark/src/log"
	"WaterMark/src/tool"
	"fmt"
	"image"
)

// imagesRGBACache
var imagesRGBACache *ImagesRGBACache

// ImagesRGBACache 缓存水印生成过程中产生的对象,加速访问
type ImagesRGBACache struct {
	// cache
	caches map[string]*image.RGBA
}

// InitImagesCache 初始化图片缓存
func InitImagesRGBACache() {
	caches := map[string]*image.RGBA{}
	imagesRGBACache = &ImagesRGBACache{
		caches: caches,
	}
	log.InfoLogger.Println("初始化图片对象缓存成功")
}

// cacheLoadImageRGBA 获取对象缓存
//
//	@param path
//	@param x1
//	@param y1
//	@param x2
//	@param y2
//	@return *image.RGBA
func cacheLoadImageRGBA(path string, x1 int, y1 int, x2 int, y2 int) *image.RGBA {
	str := path + fmt.Sprintf("%d%d%d%d", x1, y1, x2, y2)
	// 计算md5
	md5 := "cacheLoadImageRGBA:" + tool.StrMD5(str)
	// 返回缓存
	if cache, ok := imagesRGBACache.caches[md5]; ok {
		log.InfoLogger.Println("读取图片对象缓存成功:" + str)
		return cache
	}
	borderRect := image.Rect(x1, y1, x2, y2)
	image := image.NewRGBA(borderRect)
	imagesRGBACache.caches[md5] = image

	log.InfoLogger.Println("设置图片对象缓存成功:" + str)
	return image
}

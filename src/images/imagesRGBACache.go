package images

import (
	. "WaterMark/src/logs"
	"WaterMark/src/tool"
	"fmt"
	"image"
	"sync"
)

// imagesRGBACache
var imagesRGBACache sync.Map

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
	if cache, ok := imagesRGBACache.Load(md5); ok {
		Info.Println("读取图片对象缓存成功:" + str)
		return cache.(*image.RGBA)
	}
	borderRect := image.Rect(x1, y1, x2, y2)
	image := image.NewRGBA(borderRect)
	imagesRGBACache.Store(md5, image)

	Info.Println("设置图片对象缓存成功:" + str)
	return image
}

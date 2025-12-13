package internal

import (
	"image"
	"sync"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/message"
	"WaterMark/pkg"
)

// 图片文件decode缓存.
var imagesCache map[string]image.Image

// 加载图片.
func CacheLoadImageWithDecode(path string) (image.Image, pkg.EError) {
	// 计算md5
	md5 := pkg.GetStrMD5(path)

	// 返回缓存
	if cache, ok := imagesCache[md5]; ok {
		return cache, pkg.NoError
	}

	image, err := pkg.LoadImageWithDecode(path)
	if pkg.HasError(err) {
		message.SendErrorMsg(path + ":加载图片文件失败:" + err.String())

		return nil, err
	}

	return image, pkg.NoError
}

// 导入图片.
func ImportImageFiles(paths []string, exifInfos []exiftool.FileMetadata) {
	// 图片导入map
	importImageMaps := make(map[string]image.Image, 10)
	// 互斥锁,防止并发写入map报错
	var mtx sync.Mutex
	var importWg sync.WaitGroup
	// 协程数量限制
	ch := make(chan struct{}, 10)
	for i, path := range paths {
		ch <- struct{}{}
		importWg.Add(1)
		go func(file string, exifInfo exiftool.FileMetadata) {
			defer importWg.Done()
			md5 := pkg.GetStrMD5(file)
			if cache, ok := imagesCache[md5]; ok {
				mtx.Lock()
				importImageMaps[md5] = cache
				mtx.Unlock()
				<-ch

				return
			}
			image, err := pkg.LoadImageWithDecode(path)
			if pkg.HasError(err) {
				<-ch

				return
			}
			orientation := pkg.GetOrientation(pkg.AnyToString(exifInfo.Fields["Orientation"]))
			if orientation > 0 {
				image = pkg.ImageRotate(orientation, image)
			}
			mtx.Lock()
			importImageMaps[md5] = image
			mtx.Unlock()

			<-ch
		}(path, exifInfos[i])
	}
	importWg.Wait()
	imagesCache = importImageMaps
	importImageMaps = nil
}

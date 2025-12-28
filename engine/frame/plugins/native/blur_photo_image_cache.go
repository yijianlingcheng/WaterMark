package native

import (
	"sync"
	"time"
)

type blurImageFileList struct {
	blurImageWriteMap map[string]struct{}
	blurImageMap      map[string]struct{}
	blurImageMtx      sync.Mutex
}

var (
	// 确保不会重复执行.
	blurImageCacheOnce sync.Once

	// 模糊图片列表实例.
	blurImageFileListMaps *blurImageFileList
)

// 获取模糊图片列表实例.
func initBlurImageCache() {
	blurImageCacheOnce.Do(func() {
		blurImageFileListMaps = &blurImageFileList{
			blurImageMap:      make(map[string]struct{}),
			blurImageWriteMap: make(map[string]struct{}),
		}
	})
}

// 添加模糊图片到写入列表.
func addBlurImageToWriteList(path string) {
	initBlurImageCache()
	blurImageFileListMaps.blurImageMtx.Lock()
	blurImageFileListMaps.blurImageWriteMap[path] = struct{}{}
	blurImageFileListMaps.blurImageMtx.Unlock()
}

// 添加模糊图片到列表.
func moveToBlurImageFileList(path string) {
	initBlurImageCache()
	blurImageFileListMaps.blurImageMtx.Lock()
	blurImageFileListMaps.blurImageMap[path] = struct{}{}
	delete(blurImageFileListMaps.blurImageWriteMap, path)
	blurImageFileListMaps.blurImageMtx.Unlock()
}

// 检查模糊图片是否存在.
func checkBlurImageExist(path string) bool {
	initBlurImageCache()
	blurImageFileListMaps.blurImageMtx.Lock()
	_, ok := blurImageFileListMaps.blurImageMap[path]
	blurImageFileListMaps.blurImageMtx.Unlock()

	return ok
}

// 检查模糊图片是否在写入列表中.
func checkBlurImageInWriteList(path string) bool {
	initBlurImageCache()
	blurImageFileListMaps.blurImageMtx.Lock()
	_, ok := blurImageFileListMaps.blurImageWriteMap[path]
	blurImageFileListMaps.blurImageMtx.Unlock()

	return ok
}

// 等待模糊图片写入完成.
func waitBlurImageInList(path string) {
	for !checkBlurImageExist(path) {
		time.Sleep(100 * time.Millisecond)
	}
}

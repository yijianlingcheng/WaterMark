package native

import (
	"os"
	"sync"

	"github.com/golang/freetype/truetype"

	"WaterMark/internal"
	"WaterMark/pkg"
)

// 字体缓存.
var textFontCache sync.Map

// 带缓存的加载字体文件.
func loadTextFontWithCache(fontFilePath string) (*truetype.Font, pkg.EError) {
	// 计算md5
	md5 := pkg.GetStrMD5(fontFilePath)

	// 获取缓存
	if cache, ok := textFontCache.Load(md5); ok {
		if v, vok := cache.(*truetype.Font); vok {
			return v, pkg.NoError
		}
		internal.Log.Error(pkg.ImageTextCacheTypeError.String())
	}

	fontFilePath = internal.GetFontFilePath(fontFilePath)
	fontFile, err := os.ReadFile(fontFilePath)
	if err != nil {
		errMsg := fontFilePath + ":字体文件读取失败:" + err.Error()

		return nil, pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, errMsg)
	}
	fontType, err := truetype.Parse(fontFile)
	if err != nil {
		return nil, pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, fontFilePath+":字体文件解析失败:"+err.Error())
	}

	// 写入缓存
	textFontCache.Store(md5, fontType)

	return fontType, pkg.NoError
}

// 画笔字体文件初始化.
func textBrushInitFontFileToCache() pkg.EError {
	// 读取字体库下面的全部文件,全部提前初始化
	fontDir := internal.GetFontFilePath("")
	list, err := pkg.GetDirFiles(fontDir)
	if pkg.HasError(err) {
		return err
	}
	for _, item := range list {
		_, err = loadTextFontWithCache(item)
		if pkg.HasError(err) {
			return err
		}
	}

	return pkg.NoError
}

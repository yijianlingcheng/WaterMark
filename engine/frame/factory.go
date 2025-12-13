package frame

import (
	"image"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/internal"
	"WaterMark/pkg"
)

type Plugin interface {
	// 初始化
	InitPlugin() pkg.EError
	// 关闭(在APP UI退出前执行的操作)
	ClosePlugin()
	// 获取插件名称
	GetPluginName() string
	// 是否是原生插件(不依赖CGO库)
	IsNavite() bool
	// 创建图片边框对应的*image.RGBA对象
	CreateFrameImageRGBA(opts map[string]any) (*image.RGBA, pkg.EError)
	// 获取图片边框信息
	GetFrameImageBorderInfo(opts map[string]any) (map[string]any, pkg.EError)
	// 重新加载logo
	ReloadLogoImages() pkg.EError
	// 重新加载边框模板文件
	ReloadFrameTemplate() pkg.EError
	// 导入图片资源
	ImportImageFiles(paths []string, exifInfos []exiftool.FileMetadata)
}

// 初始化.
func PluginInitAll() pkg.EError {
	p := GetPlugin()

	return p.InitPlugin()
}

// 获取插件.
func GetPlugin() Plugin {
	return &NativePlugin{
		Name: internal.GetPlugin(),
	}
}

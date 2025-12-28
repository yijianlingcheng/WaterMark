package frame

import (
	"image/draw"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/engine/frame/plugins/native"
	"WaterMark/internal"
	"WaterMark/layout"
	"WaterMark/pkg"
)

type NativePlugin struct {
	Name string
}

// 初始化插件.
func (p *NativePlugin) InitPlugin() pkg.EError {
	return native.InitAllCachaAndTools()
}

// 关闭插件.
func (p *NativePlugin) ClosePlugin() {
}

// 获取插件名称.
func (p *NativePlugin) GetPluginName() string {
	return p.Name
}

// 是否是原生边框生成插件
// 后续考虑性能可能会引入libvips对应的CGO库提升速度.
func (p *NativePlugin) IsNavite() bool {
	return true
}

// 生成照片边框的RGBA数据.
func (p *NativePlugin) CreateFrameImageRGBA(opts map[string]any) (draw.Image, pkg.EError) {
	return native.CreateFrameImageRGBA(opts)
}

func (p *NativePlugin) GetFrameImageBorderInfo(opts map[string]any) (map[string]any, pkg.EError) {
	return native.GetFrameImageBorderInfo(opts)
}

// 重新加载logo路径下的全部logo照片.
func (p *NativePlugin) ReloadLogoImages() pkg.EError {
	return layout.LogosImagesInit()
}

// 导入图片.
func (p *NativePlugin) ImportImageFiles(paths []string, exifInfos []exiftool.FileMetadata) {
	internal.ImportImageFiles(paths, exifInfos)
}

// 重新导入模板.
func (p *NativePlugin) ReloadFrameTemplate() pkg.EError {
	return layout.ReloadandInitLayout()
}

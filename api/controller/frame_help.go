package controller

import (
	"encoding/json"
	"image"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/engine"
	"WaterMark/internal"
	"WaterMark/layout"
	"WaterMark/pkg"
)

// 检查模板中指定的字体文件是否存在.
//
//nolint:gocritic
func checkLayoutTemplateFont(templateLayout layout.FrameLayout) pkg.EError {
	// 检查模板中指定的字体文件是否存在，不存在则返回错误.基于目前的实现，只有4个绘制文字内容的位置
	fontFiles := []string{
		templateLayout.TextOneFontFile,
		templateLayout.TextTwoFontFile,
		templateLayout.TextThreeFontFile,
		templateLayout.TextFourFontFile,
	}

	// 读取字体库下面的全部文件,全部提前初始化
	fontDir := internal.GetFontFilePath("")
	list, err := pkg.GetDirFiles(fontDir)
	if pkg.HasError(err) {
		return err
	}
	for _, font := range fontFiles {
		if font == "" {
			continue
		}
		if !pkg.In(font, list) {
			return pkg.NewErrors(pkg.FILE_NOT_EXIST_ERROR, font+":字体文件不存在")
		}
	}

	return pkg.NoError
}

// 获取照片exif信息并检查照片logo是否配置.
func getExifAndCheckPhotoLogoExist(file string) (exiftool.FileMetadata, pkg.EError) {
	exifInfo, err := engine.CacheGetImageExif(file)
	if pkg.HasError(err) {
		return exiftool.FileMetadata{}, err
	}
	exifMake := pkg.AnyToString(exifInfo.Fields["Make"])
	logoName := layout.GetLogoNameByMake(exifMake)
	if layout.CheckLogoIsUnSupport(logoName) {
		return exiftool.FileMetadata{}, pkg.NewErrors(pkg.IMAGE_LOGO_NOT_FIND_ERROR, exifMake+":不支持的logo,请检查是否配置logo图片")
	}

	return exifInfo, pkg.NoError
}

// 照片边框压缩尺寸函数.
// 根据不同的照片尺寸选择不同展示比例.
func photoFrameResize(imageRGBA image.Image) *image.NRGBA {
	// 默认压缩比例
	ratio := defaultRatio
	// 针对4000W+像素照片特殊处理
	if imageRGBA.Bounds().Dx() > maxImageSize || imageRGBA.Bounds().Dy() > maxImageSize {
		ratio = maxRatio
	}

	return imaging.Resize(imageRGBA, imageRGBA.Bounds().Dx()/ratio, imageRGBA.Bounds().Dy()/ratio, imaging.Lanczos)
}

// 构造生成水印的参数.
func buildFramePrams(layoutStr string) (layout.FrameLayout, pkg.EError) {
	var frameLayout layout.FrameLayout
	// json序列化为布局
	jsonErr := json.Unmarshal([]byte(layoutStr), &frameLayout)
	if jsonErr != nil {
		return frameLayout, pkg.NewErrors(pkg.REQUEST_PARAM_ERROR, layoutStr+":布局信息格式错误,json解析失败")
	}
	// 查找布局
	templateLayout, findErr := layout.FindLayoutByName(frameLayout.Name)
	if pkg.HasError(findErr) {
		return frameLayout, findErr
	}
	// 将外部传递的参数合并到布局中
	jsonErr = json.NewDecoder(strings.NewReader(layoutStr)).Decode(&templateLayout)
	if jsonErr != nil {
		return frameLayout, pkg.NewErrors(pkg.REQUEST_PARAM_ERROR, layoutStr+":布局信息格式错误,json解析失败")
	}

	return templateLayout, checkLayoutTemplateFont(templateLayout)
}

// 去除字符串.
func uniqueStr(paths []string) []string {
	m := make(map[string]string)
	nPath := make([]string, 0)
	for _, path := range paths {
		if _, ok := m[path]; ok {
			continue
		}
		m[path] = path
		nPath = append(nPath, path)
	}

	return nPath
}

package frame

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/internal"
	"WaterMark/layout"
	"WaterMark/pkg"
)

var (
	// 自动生成的模板示例图片的宽.
	templateWidth = 1856

	// 自动生成的模板示例图片的高.
	templateHeight = 1238

	// 原始模板文件.
	templateFile = "template.jpg"

	// 模板名称.
	templateName = "template_%s.jpg"

	// 固定布局.
	fixedLayout = "fixed"

	// 默认厂商.
	defaultMake = "NIKON CORPORATION"
)

// 加载或者创建模板图片.
func LoadOrCreateLayoutImage() pkg.EError {
	path, err := createTemplateImage()
	if pkg.HasError(err) {
		return err
	}
	exifInfo := exiftool.EmptyFileMetadata()
	exifInfo.Fields["ImageWidth"] = float64(templateWidth)
	exifInfo.Fields["ImageHeight"] = float64(templateHeight)
	exifInfo.Fields["Make"] = defaultMake

	layouts := layout.GetAllLayout()
	plugin := GetPlugin()
	for i := range layouts {
		if layouts[i].Layout == fixedLayout {
			continue
		}
		name := layouts[i].Name
		file := internal.GetRuntimePath(fmt.Sprintf(templateName, name))
		if internal.PathExists(file) {
			continue
		}
		isBlue := false
		if layouts[i].Isblur {
			isBlue = true
		}
		go plugin.CreateFrameImageRGBA(map[string]any{
			"sourceImageFile": path,
			"exif":            exifInfo,
			"params":          layouts[i],
			"saveImageFile":   file,
			"isBlur":          isBlue,
		})
	}

	return pkg.NoError
}

// 创建模板图片.
func createTemplateImage() (string, pkg.EError) {
	path := internal.GetRuntimePath(templateFile)
	if internal.PathExists(path) {
		return path, pkg.NoError
	}

	file, err := os.Create(path)
	if err != nil {
		return "", pkg.ImageJpegSaveError
	}
	defer file.Close()

	rgba := image.NewRGBA(image.Rect(0, 0, templateWidth, templateHeight))
	for x := range templateWidth {
		for y := range templateHeight {
			rgba.Set(x, y, color.NRGBA{0, 0, 0, 255})
		}
	}
	err = jpeg.Encode(file, rgba, nil)
	if err != nil {
		return path, pkg.ImageJpegSaveError
	}

	return path, pkg.NoError
}

// 获取模板信息.
func GetTemplateInfo() map[string]string {
	layouts := layout.GetAllLayout()
	list := make(map[string]string)
	for i := range layouts {
		if layouts[i].Layout == fixedLayout {
			continue
		}
		name := layouts[i].Name
		file := internal.GetRuntimePath(fmt.Sprintf(templateName, name))

		list[name] = file
	}

	return list
}

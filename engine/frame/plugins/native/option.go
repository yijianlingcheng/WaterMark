package native

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/internal"
	"WaterMark/layout"
	"WaterMark/pkg"
)

type frameOption struct {
	Exif            exiftool.FileMetadata `mapstructure:"exif"`
	PhotoType       string                `mapstructure:"photoType"`
	SourceImageFile string                `mapstructure:"sourceImageFile"`
	SaveImageFile   string                `mapstructure:"saveImageFile"`
	Params          layout.FrameLayout    `mapstructure:"params"`
	OriginWidth     int
	OriginHeight    int
	IsAutoSave      bool
}

// 是否需要加载原始图片.
func (fp *frameOption) needSourceImage() bool {
	return fp.PhotoType != PHOTO_TYPE_BORDER
}

func (fp *frameOption) getExif() exiftool.FileMetadata {
	return fp.Exif
}

// 获取原始照片路径.
func (fp *frameOption) getSourceImageFile() string {
	return fp.SourceImageFile
}

// 获取照片width.
func (fp *frameOption) getSourceImageX() int {
	width, widthIsOk := fp.Exif.Fields["ImageWidth"].(float64)
	if widthIsOk {
		return int(width)
	}

	return 0
}

// 获取照片height.
func (fp *frameOption) getSourceImageY() int {
	height, heightIsOk := fp.Exif.Fields["ImageHeight"].(float64)
	if heightIsOk {
		return int(height)
	}

	return 0
}

// 获取照片对应的相机厂商.
func (fp *frameOption) getMakeFromExif() string {
	return pkg.AnyToString(fp.Exif.Fields["Make"])
}

// 获取一个结构体
// 利用mapstructure库将map转为frameOption.
func newFrameOption(opts map[string]any) *frameOption {
	var fp frameOption
	// 处理nil opts，确保Decode不会失败
	if len(opts) == 0 {
		opts = make(map[string]any)
	}
	err := mapstructure.Decode(opts, &fp)
	if err != nil {
		internal.Log.Panic("opts转为frameOption失败:" + err.Error())
	}

	// 确保Exif.Fields不是nil，防止nil pointer dereference
	if fp.Exif.Fields == nil {
		fp.Exif.Fields = make(map[string]any)
	}

	return &fp
}

// 重置照片width.
func (fp *frameOption) resetSourceImageX(width int) {
	fp.OriginWidth = fp.getSourceImageX()
	fp.Exif.Fields["ImageWidth"] = float64(width)
}

// 重置照片height.
func (fp *frameOption) resetSourceImageY(height int) {
	fp.OriginHeight = fp.getSourceImageY()
	fp.Exif.Fields["ImageHeight"] = float64(height)
}

// 是否是竖构图照片.
func (fp *frameOption) isVerticalImage() bool {
	if fp.OriginWidth == 0 {
		fp.OriginWidth = fp.getSourceImageX()
	}
	if fp.OriginHeight == 0 {
		fp.OriginHeight = fp.getSourceImageY()
	}

	return fp.OriginWidth < fp.OriginHeight
}

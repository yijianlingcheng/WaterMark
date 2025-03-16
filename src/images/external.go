package images

import (
	"image/color"
)

// External 外部数据,由接口导入,注入到水印模板中控制水印图生成
type External struct {
	OnlyBottom   bool
	BorderColor  color.RGBA
	Tid          string
	BorderColors string
	SourcePath   string
	SavePath     string
}

// newExternal
//
//	@return *External
func NewExternal() *External {
	return &External{}
}

// WithBoderColor 指定边框颜色
//
//	@param color
//	@return *External
func (e *External) WithBoderColor(color string) *External {
	e.BorderColors = color
	e.BorderColor = StrColor2RGBA(color)
	return e
}

// WithDefaultBoderColor 默认边框颜色
//
//	@return *External
func (e *External) WithDefaultBoderColor() *External {
	color := "255,255,255,255"
	e.BorderColors = color
	e.BorderColor = StrColor2RGBA(color)
	return e
}

// WithTid 指定使用的模板id
//
//	@param tid
//	@return *External
func (e *External) WithTid(tid string) *External {
	e.Tid = tid
	return e
}

// WithPath 指定原始图片,预览图片路径(带水印)
//
//	@param path
//	@return *External
func (e *External) WithPath(path string) *External {
	e.SavePath = getTmpPreviewPath(path)
	e.SourcePath = path
	return e
}

// WithSmallPreviewPath 指定原始图片路径,预览小图的路径(不带水印)
//
//	@param path
//	@return *External
func (e *External) WithSmallPreviewPath(path string) *External {
	e.SavePath = getSmallPreviewPath(path)
	e.SourcePath = path
	return e
}

// WithOnlyBottomFlag 指定标识
//
//	@param flag
//	@return *External
func (e *External) WithOnlyBottomFlag(flag bool) *External {
	e.OnlyBottom = flag
	return e
}

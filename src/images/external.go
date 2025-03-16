package images

import (
	"image/color"
)

// External 外部数据,由接口导入,注入到水印模板中控制水印图生成
type External struct {
	OnlyBottom         bool
	SetWordsColor      bool
	SetWords           bool
	BorderColor        color.RGBA
	Tid                string
	BorderColors       string
	SourcePath         string
	SavePath           string
	Words              string
	LensModel          string
	Model              string
	FirstWordsColors   string
	FirstWordsColor    color.RGBA
	SecondBorderColors string
	SecondBorderColor  color.RGBA
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

// WithWords 指定模板文字
//
//	@param words
//	@return *External
func (e *External) WithWords(words string) *External {
	e.SetWords = true
	e.Words = words
	return e
}

// WithModel 指定相机名称
//
//	@param model
//	@return *External
func (e *External) WithModel(model string) *External {
	e.SetWords = true
	e.Model = model
	return e
}

// WithLensModel 指定镜头名称
//
//	@param lensModel
//	@return *External
func (e *External) WithLensModel(lensModel string) *External {
	e.SetWords = true
	e.LensModel = lensModel
	return e
}

// WithFirstWordsColor 指定第一行文字颜色
//
//	@param firstWordsColor
//	@return *External
func (e *External) WithFirstWordsColor(firstWordsColor string) *External {
	e.SetWordsColor = true
	e.FirstWordsColors = firstWordsColor
	e.FirstWordsColor = StrColor2RGBA(firstWordsColor)
	return e
}

// WithSecondBorderColor 指定第二行文字颜色
//
//	@param secondBorderColor
//	@return *External
func (e *External) WithSecondBorderColor(secondBorderColor string) *External {
	e.SetWordsColor = true
	e.SecondBorderColors = secondBorderColor
	e.SecondBorderColor = StrColor2RGBA(secondBorderColor)
	return e
}

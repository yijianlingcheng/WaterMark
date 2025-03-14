package images

import (
	"image/color"
)

// External 外部数据
type External struct {
	OnlyBottom bool
	Color      color.RGBA
	Tid        string
	Colors     string
	SourcePath string
	SavePath   string
}

// newExternal
//
//	@return *External
func NewExternal() *External {
	return &External{}
}

// WithBoderColor
//
//	@param color
//	@return *External
func (e *External) WithBoderColor(color string) *External {
	e.Colors = color
	e.Color = StrColor2RGBA(color)
	return e
}

// WithDefaultBoderColor 默认边框颜色
//
//	@return *External
func (e *External) WithDefaultBoderColor() *External {
	color := "255,255,255,255"
	e.Colors = color
	e.Color = StrColor2RGBA(color)
	return e
}

// WithTid
//
//	@param tid
//	@return *External
func (e *External) WithTid(tid string) *External {
	e.Tid = tid
	return e
}

// WithTid
//
//	@param tid
//	@return *External
func (e *External) WithDefaultTid(tid string) *External {
	e.Tid = tid
	return e
}

// WithPath
//
//	@param path
//	@return *External
func (e *External) WithPath(path string) *External {
	e.SavePath = getTmpPreviewPath(path)
	e.SourcePath = path
	return e
}

// WithSmallPreviewPath
//
//	@param path
//	@return *External
func (e *External) WithSmallPreviewPath(path string) *External {
	e.SavePath = getSmallPreviewPath(path)
	e.SourcePath = path
	return e
}

// WithOnlyBottomFlag
//
//	@param flag
//	@return *External
func (e *External) WithOnlyBottomFlag(flag bool) *External {
	e.OnlyBottom = flag
	return e
}

package images

import (
	"image"
)

// separateStrategy
type separateStrategy interface {
	drawSeparate(w *WaterMark)
}

// bottomLeftSep
type bottomLeftSep struct {
	Strategy separateStrategy
}

// drawSeparate implements SeparateStrategy.
//
//	@param w
func (b *bottomLeftSep) drawSeparate(w *WaterMark) {

	// borderT 获取边框模板
	borderT := w.WT.BorderT
	// logoT 获取logo模板
	logoT := w.WT.LogoT
	// SepT 获取分隔符模板
	SepT := w.WT.SeparateT

	// 竖线画分隔符
	x := borderT.LeftWidth + logoT.Width + SepT.MarginLeft
	y := borderT.TopHeight + w.SourceHeight + SepT.MarginTop
	for i := x; i <= x+SepT.Width; i++ {
		start := image.Point{i, y}
		end := image.Point{i, y + SepT.Height}
		w.drawLine(w.Draw, start, end, SepT.FontColor)
	}
}

// bottomCenterSep
type bottomCenterSep struct {
	Strategy separateStrategy
}

// drawSeparate implements SeparateStrategy.
//
//	@param w
func (b *bottomCenterSep) drawSeparate(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WT.BorderT
	// logoT 获取logo模板
	logoT := w.WT.LogoT
	// SepT 获取分隔符模板
	SepT := w.WT.SeparateT

	// 竖线画分隔符
	x := w.Draw.Bounds().Dx() - borderT.RightWidth - logoT.MarginRight + logoT.Width + SepT.MarginLeft
	y := borderT.TopHeight + w.SourceHeight + SepT.MarginTop
	for i := x; i <= x+SepT.Width; i++ {
		start := image.Point{i, y}
		end := image.Point{i, y + SepT.Height}
		w.drawLine(w.Draw, start, end, SepT.FontColor)
	}
}

// bottomRightSep
type bottomRightSep struct {
	Strategy separateStrategy
}

// drawSeparate implements SeparateStrategy.
//
//	@param w
func (b *bottomRightSep) drawSeparate(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WT.BorderT
	// logoT 获取logo模板
	logoT := w.WT.LogoT
	// SepT 获取分隔符模板
	SepT := w.WT.SeparateT

	// 竖线画分隔符
	x := w.Draw.Bounds().Dx() - borderT.RightWidth - logoT.Width - logoT.MarginLeft - SepT.MarginRight
	y := borderT.TopHeight + w.SourceHeight + SepT.MarginTop
	for i := x; i <= x+SepT.Width; i++ {
		start := image.Point{i, y}
		end := image.Point{i, y + SepT.Height}
		w.drawLine(w.Draw, start, end, SepT.FontColor)
	}
}

// stackblurSep
type stackblurSep struct {
	Strategy separateStrategy
}

// drawSeparate drawLogo implements.
//
//	@param w
func (b *stackblurSep) drawSeparate(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WT.BorderT
	// logoT 获取logo模板
	logoT := w.WT.LogoT
	// SepT 获取分隔符模板
	SepT := w.WT.SeparateT

	// 竖线画分隔符
	x := borderT.LeftWidth + logoT.MarginLeft + logoT.Width + SepT.MarginLeft
	y := w.Draw.Bounds().Dy() - borderT.BottomHeight + SepT.MarginTop
	for i := x; i <= x+SepT.Width; i++ {
		start := image.Point{i, y}
		end := image.Point{i, y + SepT.Height}
		w.drawLine(w.Draw, start, end, SepT.FontColor)
	}
}

// SimpleSeparateFactory
type SimpleSeparateFactory struct {
}

// create
//
//	@param ext
//	@return SeparateStrategy
func (simple *SimpleSeparateFactory) create(ext string) separateStrategy {
	switch ext {
	case "BOTTOM_LOGO_LEFT":
		return &bottomLeftSep{}
	case "BOTTOM_LOGO_CENTER":
		return &bottomCenterSep{}
	case "BOTTOM_LOGO_RIGHT":
		return &bottomRightSep{}
	case "STACK_BLUR":
		return &stackblurSep{}
	}
	return nil
}

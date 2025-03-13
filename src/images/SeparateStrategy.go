package images

import (
	"image"
)

// SeparateStrategy
type SeparateStrategy interface {
	drawSeparate(w *WaterMark)
}

// BottomLeftSeparateStrategy
type BottomLeftSeparateStrategy struct {
	Strategy SeparateStrategy
}

// drawSeparate implements SeparateStrategy.
//
//	@param w
func (b *BottomLeftSeparateStrategy) drawSeparate(w *WaterMark) {

	// borderT 获取边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate
	// logoT 获取logo模板
	logoT := w.WaterMarkTemplate.LogoTemplate
	// SepT 获取分隔符模板
	SepT := w.WaterMarkTemplate.SeparateTamplate

	// 竖线画分隔符
	x := borderT.LeftWidth + logoT.Width + SepT.MarginLeft
	y := borderT.TopHeight + w.SourceHeight + SepT.MarginTop
	for i := x; i <= x+SepT.Width; i++ {
		start := image.Point{i, y}
		end := image.Point{i, y + SepT.Height}
		w.drawLine(w.Draw, start, end, SepT.FontColor)
	}
}

// BottomCenterSeparateStrategy
type BottomCenterSeparateStrategy struct {
	Strategy SeparateStrategy
}

// drawSeparate implements SeparateStrategy.
//
//	@param w
func (b *BottomCenterSeparateStrategy) drawSeparate(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate
	// logoT 获取logo模板
	logoT := w.WaterMarkTemplate.LogoTemplate
	// SepT 获取分隔符模板
	SepT := w.WaterMarkTemplate.SeparateTamplate

	// 竖线画分隔符
	x := w.Draw.Bounds().Dx() - borderT.RightWidth - logoT.MarginRight + logoT.Width + SepT.MarginLeft
	y := borderT.TopHeight + w.SourceHeight + SepT.MarginTop
	for i := x; i <= x+SepT.Width; i++ {
		start := image.Point{i, y}
		end := image.Point{i, y + SepT.Height}
		w.drawLine(w.Draw, start, end, SepT.FontColor)
	}
}

// BottomRightSeparateStrategy
type BottomRightSeparateStrategy struct {
	Strategy SeparateStrategy
}

// drawSeparate implements SeparateStrategy.
//
//	@param w
func (b *BottomRightSeparateStrategy) drawSeparate(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate
	// logoT 获取logo模板
	logoT := w.WaterMarkTemplate.LogoTemplate
	// SepT 获取分隔符模板
	SepT := w.WaterMarkTemplate.SeparateTamplate

	// 竖线画分隔符
	x := w.Draw.Bounds().Dx() - borderT.RightWidth - logoT.Width - logoT.MarginLeft - SepT.MarginRight
	y := borderT.TopHeight + w.SourceHeight + SepT.MarginTop
	for i := x; i <= x+SepT.Width; i++ {
		start := image.Point{i, y}
		end := image.Point{i, y + SepT.Height}
		w.drawLine(w.Draw, start, end, SepT.FontColor)
	}
}

// StackblurSeparateStrategy
type StackblurSeparateStrategy struct {
	Strategy SeparateStrategy
}

// drawSeparate drawLogo implements.
//
//	@param w
func (b *StackblurSeparateStrategy) drawSeparate(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate
	// logoT 获取logo模板
	logoT := w.WaterMarkTemplate.LogoTemplate
	// SepT 获取分隔符模板
	SepT := w.WaterMarkTemplate.SeparateTamplate

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
func (simple *SimpleSeparateFactory) create(ext string) SeparateStrategy {
	switch ext {
	case "BOTTOM_LOGO_LEFT":
		return &BottomLeftSeparateStrategy{}
	case "BOTTOM_LOGO_CENTER":
		return &BottomCenterSeparateStrategy{}
	case "BOTTOM_LOGO_RIGHT":
		return &BottomRightSeparateStrategy{}
	case "STACK_BLUR":
		return &StackblurSeparateStrategy{}
	}
	return nil
}

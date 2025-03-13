package images

import (
	"image"
	"image/draw"
)

// LogoStrategy
type LogoStrategy interface {
	drawLogo(w *WaterMark)
}

// BottomLeftLogoStrategy logo 在底部的左边
type BottomLeftLogoStrategy struct {
	Strategy LogoStrategy
}

// drawLogo implements LogoStrategy.
//
//	@param w
func (b *BottomLeftLogoStrategy) drawLogo(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate
	// logoT 获取logo模板
	logoT := w.WaterMarkTemplate.LogoTemplate
	// 读取LOGO
	w.loadLogo()

	watermarkX := borderT.LeftWidth
	watermarkY := w.Draw.Bounds().Max.Y - logoT.Height - logoT.MarginTop
	draw.Draw(w.Draw, image.Rectangle{Min: image.Point{X: watermarkX, Y: watermarkY}, Max: image.Point{X: watermarkX + logoT.Width, Y: watermarkY + logoT.Height}}, w.LogoImage, image.Point{0, 0}, draw.Over)
}

// BottomCenterLogoStrategy logo 在底部的中间
type BottomCenterLogoStrategy struct {
	Strategy LogoStrategy
}

// drawLogo implements LogoStrategy.
//
//	@param w
func (b *BottomCenterLogoStrategy) drawLogo(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate
	// logoT 获取logo模板
	logoT := w.WaterMarkTemplate.LogoTemplate
	// 读取LOGO
	w.loadLogo()

	watermarkX := w.Draw.Bounds().Dx() - borderT.RightWidth - logoT.MarginRight
	watermarkY := w.Draw.Bounds().Dy() - logoT.Height - logoT.MarginTop
	draw.Draw(w.Draw, image.Rectangle{Min: image.Point{X: watermarkX, Y: watermarkY}, Max: image.Point{X: watermarkX + logoT.Width, Y: watermarkY + logoT.Height}}, w.LogoImage, image.Point{0, 0}, draw.Over)
}

// BottomRightLogoStrategy logo 在底部的右边
type BottomRightLogoStrategy struct {
	Strategy LogoStrategy
}

// drawLogo implements LogoStrategy.
//
//	@param w
func (b *BottomRightLogoStrategy) drawLogo(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate
	// logoT 获取logo模板
	logoT := w.WaterMarkTemplate.LogoTemplate
	// 读取LOGO
	w.loadLogo()

	watermarkX := w.Draw.Bounds().Dx() - borderT.RightWidth - logoT.Width
	watermarkY := w.Draw.Bounds().Max.Y - logoT.Height - logoT.MarginTop
	draw.Draw(w.Draw, image.Rectangle{Min: image.Point{X: watermarkX, Y: watermarkY}, Max: image.Point{X: watermarkX + logoT.Width, Y: watermarkY + logoT.Height}}, w.LogoImage, image.Point{0, 0}, draw.Over)
}

// StackblurLogoStrategy
type StackblurLogoStrategy struct {
	Strategy LogoStrategy
}

// drawLogo implements LogoStrategy.
//
//	@param w
func (b *StackblurLogoStrategy) drawLogo(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate
	// logoT 获取logo模板
	logoT := w.WaterMarkTemplate.LogoTemplate
	// 读取LOGO
	w.loadLogo()

	watermarkX := borderT.LeftWidth + logoT.MarginLeft
	watermarkY := w.Draw.Bounds().Max.Y - borderT.BottomHeight + logoT.MarginTop

	draw.Draw(w.Draw, image.Rectangle{Min: image.Point{X: watermarkX, Y: watermarkY}, Max: image.Point{X: watermarkX + logoT.Width, Y: watermarkY + logoT.Height}}, w.TransLogoImage, image.Point{0, 0}, draw.Over)
}

// SimpleLogoFactory
type SimpleLogoFactory struct {
}

// create
//
//	@param t
//	@return LogoStrategy
func (simple *SimpleLogoFactory) create(t string) LogoStrategy {
	switch t {
	case "BOTTOM_LOGO_LEFT":
		return &BottomLeftLogoStrategy{}
	case "BOTTOM_LOGO_CENTER":
		return &BottomCenterLogoStrategy{}
	case "BOTTOM_LOGO_RIGHT":
		return &BottomRightLogoStrategy{}
	case "STACK_BLUR":
		return &StackblurLogoStrategy{}
	case "BOTTOM_LOGO_LEFT_AUTO":
		return &BottomLeftLogoAutoStrategy{}
	case "BOTTOM_LOGO_CENTER_AUTO":
		return &BottomLeftLogoAutoStrategy{}
	case "BOTTOM_LOGO_RIGHT_AUTO":
		return &BottomLeftLogoAutoStrategy{}
	case "STACK_BLUR_AUTO":
		return &BottomLeftLogoAutoStrategy{}
	}
	return nil
}

// BottomLeftLogoAutoStrategy
type BottomLeftLogoAutoStrategy struct {
	Strategy LogoStrategy
}

// drawLogo
//
//	@param w
func (b *BottomLeftLogoAutoStrategy) drawLogo(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate

	// logoT 获取logo模板
	w.WaterMarkTemplate.LogoTemplate = newLogoTemplate().WithWidth(borderT.BottomHeight).WithHeight(borderT.BottomHeight)
	logoT := w.WaterMarkTemplate.LogoTemplate
	// 读取LOGO
	w.loadLogo()

	watermarkX := borderT.LeftWidth
	watermarkY := w.Draw.Bounds().Max.Y - logoT.Height - logoT.MarginTop
	draw.Draw(w.Draw, image.Rectangle{Min: image.Point{X: watermarkX, Y: watermarkY}, Max: image.Point{X: watermarkX + logoT.Width, Y: watermarkY + logoT.Height}}, w.LogoImage, image.Point{0, 0}, draw.Over)
}

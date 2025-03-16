package images

import (
	"image"
	"image/draw"
)

// logoStrategy
type logoStrategy interface {
	drawLogo(w *WaterMark)
}

// bottomLeft logo在底部的左侧
type bottomLeft struct {
	Strategy logoStrategy
}

// drawLogo
//
//	@param w
func (b *bottomLeft) drawLogo(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WT.BorderT
	// logoT 获取logo模板
	logoT := w.WT.LogoT
	// 读取LOGO
	w.loadLogo()

	watermarkX := borderT.LeftWidth
	watermarkY := w.Draw.Bounds().Max.Y - logoT.Height - logoT.MarginTop
	draw.Draw(w.Draw, image.Rectangle{Min: image.Point{X: watermarkX, Y: watermarkY}, Max: image.Point{X: watermarkX + logoT.Width, Y: watermarkY + logoT.Height}}, w.LogoImage, image.Point{0, 0}, draw.Over)
}

// bottomCenter logo在底部的中部
type bottomCenter struct {
	Strategy logoStrategy
}

// drawLogo implements LogoStrategy.
//
//	@param w
func (b *bottomCenter) drawLogo(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WT.BorderT
	// logoT 获取logo模板
	logoT := w.WT.LogoT
	// 读取LOGO
	w.loadLogo()

	watermarkX := w.Draw.Bounds().Dx() - borderT.RightWidth - logoT.MarginRight
	watermarkY := w.Draw.Bounds().Dy() - logoT.Height - logoT.MarginTop
	draw.Draw(w.Draw, image.Rectangle{Min: image.Point{X: watermarkX, Y: watermarkY}, Max: image.Point{X: watermarkX + logoT.Width, Y: watermarkY + logoT.Height}}, w.LogoImage, image.Point{0, 0}, draw.Over)
}

// bottomRight logo在底部的右侧
type bottomRight struct {
	Strategy logoStrategy
}

// drawLogo
//
//	@param w
func (b *bottomRight) drawLogo(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WT.BorderT
	// logoT 获取logo模板
	logoT := w.WT.LogoT
	// 读取LOGO
	w.loadLogo()

	watermarkX := w.Draw.Bounds().Dx() - borderT.RightWidth - logoT.Width
	watermarkY := w.Draw.Bounds().Max.Y - logoT.Height - logoT.MarginTop
	draw.Draw(w.Draw, image.Rectangle{Min: image.Point{X: watermarkX, Y: watermarkY}, Max: image.Point{X: watermarkX + logoT.Width, Y: watermarkY + logoT.Height}}, w.LogoImage, image.Point{0, 0}, draw.Over)
}

// stackblurL 高斯模糊logo模板
type stackblurL struct {
	Strategy logoStrategy
}

// drawLogo
//
//	@param w
func (b *stackblurL) drawLogo(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WT.BorderT
	// logoT 获取logo模板
	logoT := w.WT.LogoT
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
//	@return logoStrategy
func (simple *SimpleLogoFactory) create(t string) logoStrategy {
	switch t {
	case "BOTTOM_LOGO_LEFT":
		return &bottomLeft{}
	case "BOTTOM_LOGO_CENTER":
		return &bottomCenter{}
	case "BOTTOM_LOGO_RIGHT":
		return &bottomRight{}
	case "STACK_BLUR":
		return &stackblurL{}
	case "BOTTOM_LOGO_LEFT_AUTO":
		return &bottomLeftAuto{}
	case "BOTTOM_LOGO_CENTER_AUTO":
		return &bottomLeftAuto{}
	case "BOTTOM_LOGO_RIGHT_AUTO":
		return &bottomLeftAuto{}
	case "STACK_BLUR_AUTO":
		return &bottomLeftAuto{}
	}
	return nil
}

// bottomLeftAuto logo在底部左侧,logo大小自动计算
type bottomLeftAuto struct {
	Strategy logoStrategy
}

// drawLogo
//
//	@param w
func (b *bottomLeftAuto) drawLogo(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WT.BorderT

	// logoT 获取logo模板
	w.WT.LogoT = newLogoTemplate().WithWidth(borderT.BottomHeight).WithHeight(borderT.BottomHeight)
	logoT := w.WT.LogoT
	// 读取LOGO
	w.loadLogo()

	watermarkX := borderT.LeftWidth
	watermarkY := w.Draw.Bounds().Max.Y - logoT.Height - logoT.MarginTop
	draw.Draw(w.Draw, image.Rectangle{Min: image.Point{X: watermarkX, Y: watermarkY}, Max: image.Point{X: watermarkX + logoT.Width, Y: watermarkY + logoT.Height}}, w.LogoImage, image.Point{0, 0}, draw.Over)
}

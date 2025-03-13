package images

import (
	"image"
	"image/draw"

	"github.com/disintegration/imaging"
)

// BorderStrategy
type BorderStrategy interface {
	drawBorder(w *WaterMark)
}

// NormalBorderStrategy  普通的边框样式
type NormalBorderStrategy struct {
	Strategy BorderStrategy
}

// drawBorder implements BorderStrategy.
//
//	@param w
func (b *NormalBorderStrategy) drawBorder(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate

	x := w.SourceWidth + borderT.getWidth()
	y := w.SourceHeight + borderT.getHeight()

	// des 转化为RGBA
	// borderRect 创建画布
	borderRect := image.Rect(0, 0, x, y)
	w.Draw = image.NewRGBA(borderRect)

	if borderT.IsRound {
		// 图片有圆角,需要使用png格式进行保存
		w.setPngFlag()

		// 创建圆角画布
		c := radius{p: image.Point{X: x, Y: y}, r: borderT.Radius}
		// draw 填充边框背景色
		draw.DrawMask(w.Draw, w.Draw.Bounds(), &image.Uniform{borderT.Color}, image.Point{}, &c, image.Point{}, draw.Over)
	} else {
		// draw 填充边框背景色
		draw.Draw(w.Draw, w.Draw.Bounds(), &image.Uniform{borderT.Color}, image.Point{0, 0}, draw.Src)
	}

	// draw 填充原图
	draw.Draw(w.Draw, w.SourceImage.Bounds().Add(image.Point{borderT.LeftWidth, borderT.TopHeight}), w.SourceImage, image.Point{0, 0}, draw.Over)
}

// StackblurBorderStrategy 高斯模糊边框
type StackblurBorderStrategy struct {
	Strategy BorderStrategy
}

// drawBorder implements BorderStrategy.
//
//	@param w
func (b *StackblurBorderStrategy) drawBorder(w *WaterMark) {
	// 先将原图进行高斯模糊处理
	sourceImg := w.stackblur()

	sourceWidth := sourceImg.Bounds().Dx()
	sourceHeight := sourceImg.Bounds().Dy()

	// des 转化为RGBA
	// borderRect 创建画布
	borderRect := image.Rect(0, 0, sourceWidth, sourceHeight)
	w.Draw = image.NewRGBA(borderRect)

	// borderT 获取边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate

	if borderT.IsRound {
		// 图片有圆角,需要使用png格式进行保存
		w.setPngFlag()
		// 圆角
		c := radius{p: image.Point{X: sourceWidth, Y: sourceHeight}, r: borderT.Radius}
		draw.DrawMask(w.Draw, w.Draw.Bounds(), sourceImg, image.Point{}, &c, image.Point{}, draw.Over)
	} else {
		// draw 填充高斯模糊图
		draw.Draw(w.Draw, sourceImg.Bounds(), sourceImg, image.Point{0, 0}, draw.Src)
	}
	//压缩原图
	newWidth := w.SourceWidth - borderT.getWidth()
	newHeight := w.SourceHeight - borderT.getHeight()
	newSourceImage := imaging.Resize(w.SourceImage, newWidth, newHeight, imaging.Lanczos)

	// draw 填充原图
	if borderT.IsRound {
		//圆角
		c1 := radius{p: image.Point{X: newWidth, Y: newHeight}, r: borderT.Radius}
		draw.DrawMask(w.Draw, newSourceImage.Bounds().Add(image.Point{borderT.LeftWidth, borderT.TopHeight}), newSourceImage, image.Point{0, 0}, &c1, image.Point{0, 0}, draw.Over)
	} else {
		draw.Draw(w.Draw, newSourceImage.Bounds().Add(image.Point{borderT.LeftWidth, borderT.TopHeight}), newSourceImage, image.Point{0, 0}, draw.Src)
	}
}

// SimpleBorderFactory
type SimpleBorderFactory struct {
}

// create
//
//	@param t
//	@return BorderStrategy
func (simple *SimpleBorderFactory) create(t string) BorderStrategy {
	switch t {
	case "BOTTOM_LOGO_LEFT", "BOTTOM_LOGO_CENTER", "BOTTOM_LOGO_RIGHT":
		return &NormalBorderStrategy{}
	case "STACK_BLUR":
		return &StackblurBorderStrategy{}
	}
	return nil
}

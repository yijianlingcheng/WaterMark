package images

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/imaging"
)

// borderStrategy
type borderStrategy interface {
	drawBorder(w *WaterMark)
}

// normal 普通边框样式,需要解析完整的边框模板信息
type normal struct {
	Strategy borderStrategy
}

// drawBorder 画边框
//
//	@param w
func (b *normal) drawBorder(w *WaterMark) {
	// borderT 获取边框模板
	borderT := w.WT.BorderT

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

// stackblurB 高斯模糊边框,需要解析全部的模板参数
type stackblurB struct {
	Strategy borderStrategy
}

// drawBorder 画边框
//
//	@param w
func (b *stackblurB) drawBorder(w *WaterMark) {
	// 先将原图进行高斯模糊处理
	sourceImg := w.stackblur()

	sourceWidth := sourceImg.Bounds().Dx()
	sourceHeight := sourceImg.Bounds().Dy()

	// des 转化为RGBA
	// borderRect 创建画布
	borderRect := image.Rect(0, 0, sourceWidth, sourceHeight)
	w.Draw = image.NewRGBA(borderRect)

	// borderT 获取边框模板
	borderT := w.WT.BorderT

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
//	@return borderStrategy
func (simple *SimpleBorderFactory) create(t string) borderStrategy {
	switch t {
	case "BOTTOM_LOGO_LEFT", "BOTTOM_LOGO_CENTER", "BOTTOM_LOGO_RIGHT": // 这个使用比较繁琐,需要自己去模板进行调整
		return &normal{}
	case "STACK_BLUR": // 这个使用比较繁琐,需要自己去模板进行调整
		return &stackblurB{}
	case "BOTTOM_LOGO_LEFT_AUTO", "BOTTOM_LOGO_CENTER_AUTO", "BOTTOM_LOGO_RIGHT_AUTO": // 普通自动模式.不需要指定边框的宽高
		return &normalAuto{}
	case "STACK_BLUR_AUTO":
		return &stackblurAuto{}
	}
	return nil
}

// normalAuto 自动边框样式,边框宽高自动计算生成
type normalAuto struct {
	Strategy borderStrategy
}

// drawBorder 画边框
//
//	@param w
func (b *normalAuto) drawBorder(w *WaterMark) {
	// 计算边框
	b.calculateLeftAutoBorderT(w)

	// 获取原图宽高
	x := w.SourceWidth
	y := w.SourceHeight

	borderT := w.WT.BorderT

	// des 转化为RGBA
	// borderRect 创建画布
	w.Draw = cacheLoadImageRGBA(w.SourceImgPath, 0, 0, x+borderT.getWidth(), y+borderT.getHeight())

	draw.Draw(w.Draw, w.Draw.Bounds(), &image.Uniform{borderT.Color}, image.Point{0, 0}, draw.Src)
	// draw 填充原图
	draw.Draw(w.Draw, w.SourceImage.Bounds().Add(image.Point{borderT.LeftWidth, borderT.TopHeight}), w.SourceImage, image.Point{0, 0}, draw.Over)
}

// calculateLeftAutoBorderT 计算边框的边距
//
//	@param w
func (b *normalAuto) calculateLeftAutoBorderT(w *WaterMark) {
	boderColor := color.RGBA{255, 255, 255, 255}
	if w.IsSetBorderColor {
		boderColor = w.WT.BorderT.Color
	}
	ratio := 0.05 // 比例
	boderWidth := int((float64(w.SourceWidth) * ratio) / 2)
	boderHeight := int((float64(w.SourceHeight) * ratio) / 2)

	// 转换宽高,让宽高相同
	if boderWidth > boderHeight {
		boderHeight = boderWidth
	} else {
		boderWidth = boderHeight
	}
	if w.WT.BorderT.OnlyBottom {
		w.WT.BorderT = newBorderTemplate().WithBottomHeight(boderHeight * 3).WithBoderColor(boderColor)
	} else {
		w.WT.BorderT = newBorderTemplate().WithWidth(boderWidth * 2).WithHeight(boderHeight * 4).WithBoderColor(boderColor)
	}

}

// stackblurAuto 高斯模糊自动边框
type stackblurAuto struct {
	Strategy borderStrategy
}

// drawBorder
//
//	@param w
func (b *stackblurAuto) drawBorder(w *WaterMark) {
	fmt.Println("stackblurAuto drawBorder")
}

package native

import (
	"image"
	"image/color"
	"image/draw"

	"WaterMark/pkg"
)

// loadImageRGBAWithColor 获取对象.
func loadImageRGBAWithColor(
	x1, y1, x2, y2 int,
	color color.RGBA,
) (*image.RGBA, pkg.EError) {
	// 创建画布
	imageFrame := loadImageRGBAWithColorAndDraw(x1, y1, x2, y2, color)

	return imageFrame, pkg.NoError
}

// loadImageRGBA 获取对象.
func loadImageRGBA(x1, y1, x2, y2 int) *image.RGBA {
	// 创建
	borderRect := image.Rect(x1, y1, x2, y2)
	image := image.NewRGBA(borderRect)

	return image
}

// loadImageRGBAWithColorAndDraw 获取已填充背景的空白图片对象.
func loadImageRGBAWithColorAndDraw(x1, y1, x2, y2 int, color color.RGBA) *image.RGBA {
	// 创建
	imageRGBAWithColor := loadImageRGBA(x1, y1, x2, y2)
	draw.Draw(imageRGBAWithColor, imageRGBAWithColor.Bounds(), &image.Uniform{color}, image.Point{0, 0}, draw.Src)

	return imageRGBAWithColor
}

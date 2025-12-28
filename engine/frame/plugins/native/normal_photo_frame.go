package native

import (
	"image"
	"image/draw"
	"sync"

	"WaterMark/layout"
	"WaterMark/pkg"
)

type photoFrame struct {
	basePhotoFrame
}

// 画主体.
func (fm *photoFrame) drawMainImage(wg *sync.WaitGroup) {
	defer wg.Done()

	// 生成照片主体
	if fm.opts.needSourceImage() {
		draw.Draw(
			fm.frameDraw,
			fm.srcImage.imgDecode.Bounds().Add(image.Point{fm.borImage.leftWidth, fm.borImage.topHeight}),
			fm.srcImage.imgDecode,
			image.Pt(0, 0),
			draw.Over,
		)
	}
}

// 画边框与文字.
func (fm *photoFrame) drawBorderImage(wg *sync.WaitGroup) pkg.EError {
	defer wg.Done()

	// 生成边框对象
	fm.borderDraw = loadImageRGBA(0, 0, fm.srcImage.width, fm.borImage.bottomHeight)

	draw.Draw(
		fm.borderDraw,
		fm.borderDraw.Bounds(),
		&image.Uniform{fm.borImage.bgColor},
		image.Point{0, 0},
		draw.Src,
	)
	// 相机的logo如果没有找到,则使用特定标识的logo进行代替
	logo, logoErr := layout.GetLogoImageByNameAndWidhtAndHeight(
		layout.GetLogoNameByMake(fm.opts.getMakeFromExif()),
		fm.borImage.logoLay.layout.width,
		fm.borImage.logoLay.layout.height,
	)
	if pkg.HasError(logoErr) {
		return logoErr
	}
	fm.borImage.logoLay.item = logo
	simpleBorderFactory := &SimpleBorderFactory{}
	simpleBorderFactory.createBorder(fm.opts.Params.Name).drawBorder(fm)

	return pkg.NoError
}

// 画出照片主体与边框
// 为了性能考虑采用协程组实现.
func (fm *photoFrame) drawFrame() {
	// 协程组
	var wg sync.WaitGroup
	// 2协程
	wg.Add(2)

	// 画主体
	go fm.drawMainImage(&wg)
	// 画边框
	go fm.drawBorderImage(&wg)

	wg.Wait()
}

// 拼接这两张图片
// 将照片与生成好的边框水印图片拼接在一起.
func (fm *photoFrame) drawMerge() draw.Image {
	draw.Draw(
		fm.frameDraw,
		image.Rect(
			fm.borImage.leftWidth,
			fm.borImage.topHeight+fm.srcImage.height,
			fm.borImage.leftWidth+fm.srcImage.width+fm.borImage.rightWidth,
			fm.borImage.topHeight+fm.srcImage.height+fm.borImage.bottomHeight,
		),
		fm.borderDraw,
		image.Pt(0, 0),
		draw.Src,
	)

	return fm.frameDraw
}

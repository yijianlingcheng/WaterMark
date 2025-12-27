package native

import (
	"image"
	"image/draw"
	"path/filepath"

	"github.com/disintegration/imaging"

	"WaterMark/internal"
	"WaterMark/message"
	"WaterMark/pkg"
)

type blurPhotoFrame struct {
	basePhotoFrame
}

// 初始化.
func (fm *blurPhotoFrame) initFrame(opts map[string]any) pkg.EError {
	fm.opts = newFrameOption(opts)

	sourceImagePath := fm.opts.getSourceImageFile()
	// 初始化各种对象
	fm.srcImage = newSourceImage(sourceImagePath)
	// 判断是否需要旋转图片
	fm.resetSourceImageXAndY()

	borImage, borderErr := getBorderImage(fm)
	if pkg.HasError(borderErr) {
		message.SendErrorMsg(borderErr.String())

		return borderErr
	}
	fm.borImage = borImage

	// 判断是否需要加载原图
	if !fm.checkBlurBackgroundImageExist() {
		sourceImage, loadSourceImageErr := fm.loadSourceImage(sourceImagePath)
		if pkg.HasError(loadSourceImageErr) {
			return loadSourceImageErr
		}
		fm.srcImage.SetImage(sourceImage)
	}

	// 初始化最终绘制对象
	fm.finImage = newFinalImage(fm.opts)

	var frame draw.Image
	var frameErr pkg.EError

	// 模糊模板特殊处理
	frame, frameErr = fm.createBlurDraw()
	if pkg.HasError(frameErr) {
		return frameErr
	}
	fm.frameDraw = frame

	return pkg.NoError
}

// 创建模糊模板画布.
func (fm *blurPhotoFrame) createBlurDraw() (*image.NRGBA, pkg.EError) {
	// 判断是否已经存在模糊背景图片
	if fm.checkBlurBackgroundImageExist() {
		blurImage, loadErr := internal.CacheLoadImageWithDecode(fm.getBlurBackgroundImageFilePath())
		if pkg.HasError(loadErr) {
			return nil, loadErr
		}

		return imaging.Clone(blurImage), pkg.NoError
	}
	sourceImage := fm.srcImage.imgDecode
	// 减小
	newImage := imaging.Resize(sourceImage,
		sourceImage.Bounds().Dx()/5,
		sourceImage.Bounds().Dy()/5, imaging.Lanczos)
	// 高斯模糊
	newImage = imaging.Blur(newImage, 50)
	// 放大
	newImage = imaging.Resize(newImage, fm.finImage.width, fm.finImage.height, imaging.Lanczos)

	// 保存模糊之后的背景图片
	autoSaveBlurBackgroundImage(fm.getBlurBackgroundImageFilePath(), imaging.Clone(newImage), 75)

	return newImage, pkg.NoError
}

// 获取模糊模板处理之后的的图片文件路径.
func (fm *blurPhotoFrame) getBlurBackgroundImageFilePath() string {
	return "./runtime/blur/" + filepath.Base(fm.srcImage.path)
}

// 检查模糊模板处理之后的的图片文件是否存在.
func (fm *blurPhotoFrame) checkBlurBackgroundImageExist() bool {
	return internal.PathExists(fm.getBlurBackgroundImageFilePath())
}

// 加载图片.
func (fm *blurPhotoFrame) loadSourceImage(path string) (image.Image, pkg.EError) {
	image, loadErr := internal.CacheLoadImageWithDecode(path)
	if pkg.HasError(loadErr) {
		return nil, loadErr
	}
	orientation := pkg.GetOrientation(pkg.AnyToString(fm.opts.getExif().Fields["Orientation"]))
	if orientation > 0 {
		image = pkg.ImageRotate(orientation, image)
	}

	return image, pkg.NoError
}

// 绘制模糊模板主体. 绘制圆角阴影使用imagemagick工具实现.
func (fm *blurPhotoFrame) drawBlurMainImage() {
	if !fm.opts.needSourceImage() {
		return
	}
	w := fm.opts.getSourceImageX()
	borderRadius := w * fm.opts.Params.BorderRadius / 1000
	// 判断是否需要画圆角
	if borderRadius == 0 {
		// 不需要画圆角
		return
	}
	// 画圆角
}

// 画模糊边框与文字.
func (fm *blurPhotoFrame) drawBlurBorderImage() pkg.EError {
	// 生成边框对象
	fm.borderDraw = loadImageRGBA(0, 0, fm.srcImage.width, fm.borImage.bottomHeight)

	simpleBorderFactory := &SimpleBorderFactory{}
	simpleBorderFactory.createBorder(fm.opts.Params.Name).drawBorder(fm)

	return pkg.NoError
}

// 画出照片主体与边框
// 为了性能考虑采用协程组实现.
func (fm *blurPhotoFrame) drawFrame() {
	// 画主体
	fm.drawBlurMainImage()
	// 画边框
	fm.drawBlurBorderImage()
}

// 获取最终合成的图片.
func (fm *blurPhotoFrame) drawBlurMerge() draw.Image {
	return fm.frameDraw
}

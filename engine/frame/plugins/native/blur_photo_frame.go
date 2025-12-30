package native

import (
	"fmt"
	"image"
	"image/draw"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"

	"WaterMark/internal"
	"WaterMark/internal/cmd"
	"WaterMark/message"
	"WaterMark/pkg"
)

type blurPhotoFrame struct {
	basePhotoFrame

	// 临时模糊结果图片存放路径.
	tmpBlurResultImagePath string
}

// 初始化.
func (fm *blurPhotoFrame) initFrame(opts map[string]any) pkg.EError {
	if len(opts) == 0 {
		return pkg.RequestParamNilError
	}
	fm.isBlur = true
	fm.opts = newFrameOption(opts)
	sourceImagePath := fm.opts.getSourceImageFile()
	// 初始化各种对象
	fm.srcImage = newSourceImage(sourceImagePath)
	// 判断是否需要旋转图片
	fm.resetSourceImageXAndY()
	borImage, borderErr := getBorderImage(fm)
	if pkg.HasError(borderErr) {
		return borderErr
	}
	fm.borImage = borImage
	// 判断是否需要加载原图
	if fm.opts.needSourceImage() && !checkBlurImageExist(fm.getBlurBackgroundImageFilePath()) {
		sourceImage, loadSourceImageErr := fm.loadSourceImage(sourceImagePath)
		if pkg.HasError(loadSourceImageErr) {
			return loadSourceImageErr
		}
		fm.srcImage.SetImage(sourceImage)
	} else {
		fm.srcImage.SetImageXAndY(fm.opts.getSourceImageX(), fm.opts.getSourceImageY())
	}
	// 初始化最终绘制对象
	fm.finImage = newFinalImage(fm.opts)

	var frame draw.Image
	var frameErr pkg.EError
	// 需要加载原图的时候创建模糊模板
	if fm.opts.needSourceImage() {
		frame, frameErr = fm.createBlurDraw()
	} else {
		frame, frameErr = fm.createTransparentDraw()
	}
	if pkg.HasError(frameErr) {
		return frameErr
	}
	fm.frameDraw = frame
	fm.setTmpBlurResultImagePath()

	return pkg.NoError
}

// 创建透明画布.
func (fm *blurPhotoFrame) createTransparentDraw() (draw.Image, pkg.EError) {
	return loadImageRGBA(0, 0, fm.finImage.width, fm.finImage.height),
		pkg.NoError
}

// 创建模糊模板画布.
func (fm *blurPhotoFrame) createBlurDraw() (*image.NRGBA, pkg.EError) {
	blurBackgroundImagePath := fm.getBlurBackgroundImageFilePath()
	// 判断是否已经存在模糊背景图片
	if checkBlurImageExist(blurBackgroundImagePath) {
		blurImage, loadErr := pkg.LoadImageWithDecode(blurBackgroundImagePath)
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
	newImage = imaging.Blur(newImage, 20)
	// 放大
	newImage = imaging.Resize(newImage, fm.finImage.width, fm.finImage.height, imaging.Lanczos)

	blurImage := imaging.Clone(newImage)
	// 保存模糊之后的背景图片
	go func(path string, img *image.NRGBA) {
		// 判断是否已经在写入列表中,如果在则不进行写入操作
		if checkBlurImageInWriteList(path) {
			return
		}
		// 添加到待写入列表
		addBlurImageToWriteList(path)
		saveJpgImage(path, img, 100)
		// 从待写入列表移除
		moveToBlurImageFileList(path)
	}(blurBackgroundImagePath, blurImage)

	return newImage, pkg.NoError
}

// 获取模糊模板处理之后的的图片文件路径.
func (fm *blurPhotoFrame) getBlurBackgroundImageFilePath() string {
	return internal.GetAppBlurFilePath(filepath.Base(fm.srcImage.path))
}

// 设置临时模糊结果图片存放路径.
func (fm *blurPhotoFrame) setTmpBlurResultImagePath() {
	if !fm.opts.needSourceImage() {
		return
	}
	fm.tmpBlurResultImagePath = internal.GetAppBlurFilePath("tmp_" + filepath.Base(fm.srcImage.path))
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

// 获取保存图片地址.
func (fm *blurPhotoFrame) getBlurSaveImageFile() string {
	if fm.opts.SaveImageFile != "" {
		return fm.opts.SaveImageFile
	}

	return fm.tmpBlurResultImagePath
}

// 绘制模糊模板主体. 绘制圆角阴影使用imagemagick工具实现.
func (fm *blurPhotoFrame) drawBlurMainImage() {
	if !fm.opts.needSourceImage() {
		return
	}
	// 等待模糊图片生成完成
	waitBlurImageInList(fm.getBlurBackgroundImageFilePath())

	w := fm.opts.getSourceImageX()
	h := fm.opts.getSourceImageY()
	borderRadius := max(w, h) * fm.opts.Params.BorderRadius / 1000

	_, imageErr := cmd.CommandRunWithArgs(5*time.Minute, fm.getMagickCmdArgs(borderRadius, fm.getBlurSaveImageFile()))
	if pkg.HasError(imageErr) {
		message.SendErrorMsg(imageErr.String())

		return
	}

	// 加载生成的图片
	image, loadErr := internal.CacheLoadImageWithDecode(fm.getBlurSaveImageFile())
	if pkg.HasError(loadErr) {
		message.SendErrorMsg(loadErr.String())

		return
	}
	// 画图
	fm.frameDraw = imaging.Clone(image)
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

// 获取magick命令参数.
func (fm *blurPhotoFrame) getMagickCmdArgs(borderRadius int, saveFilePath string) []string {
	shadowWith := fm.opts.getSourceImageX() + 4*borderRadius
	shadowHeight := fm.opts.getSourceImageY() + 4*borderRadius

	// 计算位置
	marginLeft := fm.opts.Params.MainMarginLeft
	marginTop := fm.opts.Params.MainMarginTop
	if borderRadius > 0 {
		marginLeft = (fm.finImage.width - shadowWith) / 2
		marginTop = fm.opts.Params.MainMarginTop - borderRadius*2
	}
	// 合并图片时候的位置参数
	position := fmt.Sprintf("%dx%d+%d+%d", shadowWith, shadowHeight, marginLeft, marginTop)
	// 没有圆角
	if borderRadius == 0 {
		return []string{
			internal.GetMagickBinPath(),
			fm.getBlurBackgroundImageFilePath(), "-compose", "over",
			fm.srcImage.path, "-auto-orient", "-geometry", position, "-composite",
			saveFilePath,
		}
	}
	// 阴影
	shadow := fmt.Sprintf("50x%d+0+0", borderRadius)
	// 圆角
	roundrectangle := fmt.Sprintf(
		"roundrectangle 0,0 %d,%d %d,%d",
		fm.opts.getSourceImageX(),
		fm.opts.getSourceImageY(),
		borderRadius,
		borderRadius,
	)

	return []string{
		internal.GetMagickBinPath(),
		fm.getBlurBackgroundImageFilePath(),
		"(", fm.srcImage.path, "-auto-orient", "(", "+clone", "-alpha", "extract", "-draw", roundrectangle,
		"-negate", ")", "-alpha", "off", "-compose", "CopyOpacity", "-composite", "(", "+clone", "-background",
		"grey", "-shadow", shadow, ")", "+swap", "-background", "none", "-compose", "over", "-layers", "merge",
		"+repage", ")", "-geometry", position, "-compose", "over", "-composite", saveFilePath,
	}
}

// 获取最终合成的图片.
func (fm *blurPhotoFrame) drawBlurMerge() draw.Image {
	return fm.frameDraw
}

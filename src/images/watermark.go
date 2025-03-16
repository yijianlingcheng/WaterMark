package images

import (
	"WaterMark/src/cmd"
	"WaterMark/src/exif"
	. "WaterMark/src/logs"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/esimov/stackblur-go"
	"github.com/fatih/structs"
)

// WaterMark 水印
type WaterMark struct {

	// IsSavePng 是否保存为png格式文件(模板中使用了圆角,图片边框拥有透明度,因此需要将图片保存为png)
	IsSavePng bool

	// IsSetBorderColor 是否外部设置了边框颜色
	IsSetBorderColor bool

	// Quality 图片质量,jpeg图片在保存的时候需要指定图片质量
	Quality int

	// SourceWidth 原始图片宽
	SourceWidth int

	// SourceHeight 原始图片高
	SourceHeight int

	// LogoImgPath logo图片地址
	LogoImgPath string

	// TransLogoImgPath 高斯模糊模板使用的logo
	TransLogoImgPath string

	// SourceImgPath 原始图片地址
	SourceImgPath string

	// SaveImgPath 保存图片地址
	SaveImgPath string

	// LogoImage 相机logo图片
	LogoImage image.Image

	// TransLogoImage 特殊的相机logo
	TransLogoImage image.Image

	// SourceImage 原始图片
	SourceImage image.Image

	// Draw 绘画
	Draw *image.RGBA

	// ExifInfo exif信息
	ExifInfo exif.Exif

	// ExifMap
	ExifMap map[string]any

	// WT 水印模板
	WT *WaterMarkTemplate
}

// setPngFlag 保存的图片需要是png格式
func (w *WaterMark) setPngFlag() {
	w.IsSavePng = true
}

// 初始化
func newWaterMark() *WaterMark {
	return &WaterMark{
		Quality: 100,
	}
}

// loadSource 加载资源
//
//	@param path
//	@param save
//	@param tplId
//	@return error
func (w *WaterMark) loadSource(path string, save string, tplId string) error {
	exifInfo, err := cmd.CacheLoadExifTool(path)
	if err != nil {
		return err
	}
	logoPath, translogoPath, err := getLogoPath(exifInfo.Make)
	if err != nil {
		return err
	}
	tpl, err := findTemplateById(tplId)
	if err != nil {
		return err
	}
	if tpl.BorderT == nil || tpl.LogoT == nil || tpl.SeparateT == nil || tpl.WordsT == nil {
		err := "模板自动实例化失败,部分对象为空.模板ID:" + tplId
		return errors.New(err)
	}
	w.setImgOptions(logoPath, translogoPath, path, save).setExif(exifInfo).loadTemplate(tpl)
	return nil
}

// setExif 设置exif信息
//
//	@param e
//	@return *WaterMark
func (w *WaterMark) setExif(e exif.Exif) *WaterMark {
	w.ExifInfo = e
	w.ExifMap = structs.Map(e)
	return w
}

// loadTemplate 加载模板
//
//	@param t
//	@return *WaterMark
func (w *WaterMark) loadTemplate(t WaterMarkTemplate) *WaterMark {
	w.WT = newEmptyWaterMarkTemplate()
	w.WT = &t
	return w
}

// setImgOptions 设置图片信息
//
//	@param logo1
//	@param logo2
//	@param source
//	@param save
//	@return *WaterMark
func (w *WaterMark) setImgOptions(logo1 string, logo2, source string, save string) *WaterMark {
	w.LogoImgPath = logo1
	w.TransLogoImgPath = logo2
	w.SourceImgPath = source
	w.SaveImgPath = save
	return w
}

// loadLogo 加载logo
//
//	@return error
func (w *WaterMark) loadLogo() error {

	//获取logo模板
	logoT := w.WT.LogoT

	// 是高斯模糊模板,使用特殊的透明LOGO
	if w.WT.Stackblur {
		// 根据logo图片地址,logo宽高获取图片,如果对应尺寸图片不存在,则重新生成一张并返回
		w.TransLogoImgPath = GetImageByWidthHeight(w.TransLogoImgPath, logoT.Width, logoT.Height)
		// 加载图片
		transLogoImg, err := cacheLoadImage(w.TransLogoImgPath)
		if err != nil {
			return err
		}
		w.TransLogoImage = transLogoImg
	} else {
		// 根据logo图片地址,logo宽高获取图片,如果对应尺寸图片不存在,则重新生成一张并返回
		w.LogoImgPath = GetImageByWidthHeight(w.LogoImgPath, logoT.Width, logoT.Height)
		// 加载图片
		logoImg, err := cacheLoadImage(w.LogoImgPath)
		if err != nil {
			return err
		}
		w.LogoImage = logoImg
	}
	return nil
}

// loadSourceImg 加载原始图片
//
//	@return error
func (w *WaterMark) loadSourceImg() error {
	sourceImg, err := cacheLoadImage(w.SourceImgPath)
	if err != nil {
		return err
	}
	if w.ExifInfo.OrientationNum > 0 {
		var newSourceImg *image.NRGBA
		switch w.ExifInfo.OrientationNum {
		case 90:
			newSourceImg = imaging.Rotate90(sourceImg) // 逆时针90
		case 180:
			newSourceImg = imaging.Rotate180(sourceImg) // 逆时针180
		case 270:
			newSourceImg = imaging.Rotate270(sourceImg) // 逆时针270度
		default:
			newSourceImg = imaging.Clone(sourceImg) // 无需旋转
		}
		w.SourceImage = newSourceImg
		w.SourceWidth = newSourceImg.Bounds().Dx()
		w.SourceHeight = newSourceImg.Bounds().Dy()
	} else {
		w.SourceImage = sourceImg
		w.SourceWidth = sourceImg.Bounds().Dx()
		w.SourceHeight = sourceImg.Bounds().Dy()
	}
	return nil
}

// exportExternal 导入外部数据
//
//	@param e
func (w *WaterMark) exportExternal(e *External) {
	// 设置边框标识
	w.setBorderOnlyBottom(e.OnlyBottom)
	// 设置边框颜色
	w.setBorderColor(e.BorderColor)
}

// setBorderOnlyBottom 设置模板
//
//	@param flag
func (w *WaterMark) setBorderOnlyBottom(flag bool) {
	w.WT.BorderT.OnlyBottom = flag
}

// setBorderColor 设置边框颜色
//
//	@param color
func (w *WaterMark) setBorderColor(color color.RGBA) {
	w.WT.BorderT.Color = color
	w.IsSetBorderColor = true
}

// beforeProcess 前置处理
func (w *WaterMark) beforeProcess() {
	// 只有底部边框的模式,border模板的top,left,bottom需要赋0
	if w.WT.BorderT.OnlyBottom {
		w.WT.BorderT.LeftWidth = 0
		w.WT.BorderT.RightWidth = 0
		w.WT.BorderT.TopHeight = 0
	}
}

// drawLogo2Image 将相机logo写入图片中
//
//	@return *WaterMark
func (w *WaterMark) drawLogo2Image() *WaterMark {
	//前置处理
	w.beforeProcess()

	// 画边框
	simpleBorderFactory := &SimpleBorderFactory{}
	borderStrategy := simpleBorderFactory.create(w.WT.Type)
	borderStrategy.drawBorder(w)

	// 填充logo,此处需要判断水印模板的类型
	// 画logo
	simpleLogoFactory := &SimpleLogoFactory{}
	logoStrategy := simpleLogoFactory.create(w.WT.Type)
	logoStrategy.drawLogo(w)

	// 画分隔符
	SepT := w.WT.SeparateT
	if SepT.Exist {
		simpleSeparateFactory := &SimpleSeparateFactory{}
		separateFactory := simpleSeparateFactory.create(w.WT.Type)
		separateFactory.drawSeparate(w)
	}
	return w
}

// drawLine 画线
//
//	@param img
//	@param start
//	@param end
//	@param c
//	@return *WaterMark
func (w *WaterMark) drawLine(img draw.Image, start, end image.Point, c color.Color) *WaterMark {
	dx := abs(end.X - start.X)
	dy := abs(end.Y - start.Y)
	sx, sy := 0, 0
	if start.X < end.X {
		sx = 1
	} else if start.X > end.X {
		sx = -1
	}
	if start.Y < end.Y {
		sy = 1
	} else if start.Y > end.Y {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(start.X, start.Y, c)
		if start.X == end.X && start.Y == end.Y {
			break
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			start.X += sx
		}
		if e2 < dx {
			err += dx
			start.Y += sy
		}
	}
	return w
}

// abs 返回整数的绝对值
//
//	@param x
//	@return int
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// getWords 获取文字信息
//
//	@param t
//	@return string
func (w *WaterMark) getWords(t string) string {
	tpl := structs.Map(w.WT.WordsT)
	if v, ok := tpl[t]; ok {
		str := v.(string)
		list := strings.Split(str, ",")
		r := []string{}
		for _, item := range list {
			r = append(r, w.ExifMap[item].(string))
		}
		return strings.Join(r, " ")
	}
	return ""
}

// drawFont2Image 将文字信息写入图片中
//
//	@return *WaterMark
func (w *WaterMark) drawFont2Image() *WaterMark {
	simpleWordFactory := &SimpleWordFactory{}
	wordFactory := simpleWordFactory.create(w.WT.Type)
	wordFactory.drawWords(w)
	return w
}

// stackblur 高斯模糊,先对图片进行高斯模糊处理,在生成小图片覆盖在原图上并添加水印
func (w *WaterMark) stackblur() *image.NRGBA {
	// return imaging.Blur(w.SourceImage, float64(w.WT.BlurRadius))
	// 使用stackblur.Process进行高斯模糊,比imaging.Blur消耗低
	t, _ := stackblur.Process(w.SourceImage, uint32(w.WT.BlurRadius))
	return t
}

// saveImg 保存图片
func (w *WaterMark) saveImg() {
	if w.IsSavePng {
		savePngImage(w.SaveImgPath, w.Draw)
	} else {
		saveJpegImage(w.SaveImgPath, w.Draw, w.Quality)
	}
	exif.CoverImgExifInfo(w.SaveImgPath, w.ExifInfo)
}

// exportData 导出模板信息
//
//	@return map
func (w *WaterMark) exportData() map[string]string {
	r := map[string]string{}
	r["BorderColors"] = Color2Str(w.WT.BorderT.Color)
	r["SaveImgPath"] = w.SaveImgPath
	r["SourceImgPath"] = w.SourceImgPath
	return r
}

// exportErrorData exportData 导出错误的模板信息
//
//	@param err
//	@return map
func (w *WaterMark) exportErrorData(err error) map[string]string {
	r := map[string]string{}
	r["BorderColors"] = ""
	r["SaveImgPath"] = ""
	r["SourceImgPath"] = ""
	r["error"] = err.Error()
	return r
}

// ProcessWaterMark 生成水印
//
//	@param tid 模板id
//	@param path 图片路径
//	@param save 目标图片路径
func processWaterMark(tid string, path string, save string) {
	waterMark := newWaterMark()
	// 加载资源
	if err := waterMark.loadSource(path, save, tid); err != nil {
		Errors.Println(err)
	}
	// 读取原始图片
	if err := waterMark.loadSourceImg(); err != nil {
		Errors.Println(err)
	}
	// 生成水印
	waterMark.drawLogo2Image().drawFont2Image()
	// 保存图片
	waterMark.saveImg()
}

// 水印预览图生成位置
var PreviewPath string

// 预览小图生成位置
var SmallPreviewPath string

// getTmpPreviewPath 获取预览的临时目录
//
//	@param path
//	@return string
func getTmpPreviewPath(path string) string {
	t := strings.Split(path, "/")
	return PreviewPath + t[len(t)-1]
}

// getSmallPreviewPath
//
//	@param path
//	@return string
func getSmallPreviewPath(path string) string {
	t := strings.Split(path, "/")
	return SmallPreviewPath + t[len(t)-1]
}

// CreatePreviewWaterMark 生成水印预览图片
//
//	@param e
//	@return map
func CreatePreviewWaterMark(e *External) map[string]string {
	waterMark := newWaterMark()
	// 加载资源
	if err := waterMark.loadSource(e.SourcePath, e.SavePath, e.Tid); err != nil {
		Errors.Println(err)
		return waterMark.exportErrorData(err)
	}
	// 读取原始图片
	if err := waterMark.loadSourceImg(); err != nil {
		Errors.Println(err)
		return waterMark.exportErrorData(err)
	}
	waterMark.exportExternal(e)
	// 生成水印
	waterMark.drawLogo2Image().drawFont2Image()
	// 保存图片
	waterMark.saveImg()

	return waterMark.exportData()
}

// CreateSmallPreview 生成预览小图
//
//	@param e
func CreateSmallPreview(e *External) {
	// 需要实现回写exif的角度字段,解决浏览器预览角度不对的问题
	img, _ := cacheLoadImage(e.SourcePath)
	newImg := imaging.Resize(img, img.Bounds().Dx()/10, img.Bounds().Dy()/10, imaging.Lanczos)
	saveJpegImage(e.SavePath, newImg, 100)
}

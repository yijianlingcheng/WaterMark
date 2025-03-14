package images

import (
	"WaterMark/src/log"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

// imgCompress 图片实际压缩实现
//
//	@param buf 打开的图片buff
//	@param w 图片width
//	@param h 图片height
//	@param q 图片质量，jpeg格式需要
//	@return []byte
//	@return error
func imgCompress(buf []byte, w uint, h uint, q int) ([]byte, error) {

	decodeBuf, layout, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	set := resize.Resize(w, h, decodeBuf, resize.Lanczos3)
	NewBuf := bytes.Buffer{}

	switch layout {
	case "jpeg", "jpg":
		err = jpeg.Encode(&NewBuf, set, &jpeg.Options{Quality: q})
	case "png":
		err = png.Encode(&NewBuf, set)
	default:
		return nil, errors.New("该图片格式不支持压缩")
	}

	if err != nil {
		return nil, err
	}

	if NewBuf.Len() < len(buf) {
		buf = NewBuf.Bytes()
	}

	return buf, nil
}

// ImgResize 图片调整大小
//
//	@param s 源图片路径
//	@param o 待保存的图片路径
//	@param w 待保存图片的width
//	@param h 待保存图片的height
//	@return error
func ImgResize(s string, o string, w uint, h uint) error {

	file, err := os.Open(s)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)

	compressedBuf, err := imgCompress(buf.Bytes(), w, h, 100)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(o)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	outputFile.Write(compressedBuf)

	return nil
}

// loadImage 将指定路径的图片读取并编码,以便后续程序使用,目前只支持JPEG,PNG两种格式
//
//	@param path 图片路径
//	@return img
//	@return err
func loadImage(path string) (img image.Image, err error) {

	r, err := os.Open(path)
	if err != nil {
		log.ErrorLogger.Println(path + "文件打开失败:" + err.Error())
		return nil, err
	}
	defer r.Close()

	buff := make([]byte, 512)
	_, err = r.Read(buff)
	if err != nil {
		log.ErrorLogger.Println(path + "文件读取失败:" + err.Error())
		return nil, err
	}
	filetype := http.DetectContentType(buff) //根据http库获取文件类型

	io, err := os.Open(path)
	if err != nil {
		log.ErrorLogger.Println(path + "文件打开失败:" + err.Error())
		return nil, err
	}
	defer io.Close()

	switch filetype {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(io)
		if err != nil {
			log.ErrorLogger.Println(path + "jpeg.Decode 失败:" + err.Error())
			return nil, err
		}
	case "image/png":
		img, err = png.Decode(io)
		if err != nil {
			log.ErrorLogger.Println(path + "png.Decode 失败:" + err.Error())
			return nil, err
		}
	default:
		log.ErrorLogger.Println(path + "文件不是支持的格式")
		return nil, errors.New("文件不是支持的格式")
	}
	return img, nil
}

// saveJpegImage 将内存中编码的图片保存为jpeg图片
//
//	@param path 保存图片的路径
//	@param m 编码的图像
//	@param q 图片质量
//	@return error
func saveJpegImage(path string, m image.Image, q int) error {
	io, err := os.Create(path)
	if err != nil {
		log.ErrorLogger.Println(path + "创建文件出错:" + err.Error())
		return err
	}
	defer io.Close()
	err = jpeg.Encode(io, m, &jpeg.Options{Quality: q})
	if err != nil {
		log.ErrorLogger.Println(path + "图片保存失败:" + err.Error())
		return err
	}
	return nil
}

// savePngImage 保存png图片
//
//	@param path
//	@param m
//	@return error
func savePngImage(path string, m image.Image) error {
	io, err := os.Create(path)
	if err != nil {
		log.ErrorLogger.Println(path + "创建文件出错:" + err.Error())
		return err
	}
	defer io.Close()
	err = png.Encode(io, m)
	if err != nil {
		log.ErrorLogger.Println(path + "图片保存失败:" + err.Error())
		return err
	}
	return nil
}

// strColor2RGBA 将字符串RGBA颜色转换成color.RGBA
//
//	@param s
//	@return color.RGBA
func StrColor2RGBA(s string) color.RGBA {
	list := strings.Split(s, ",")
	r0, _ := strconv.ParseUint(list[0], 10, 8)
	r1, _ := strconv.ParseUint(list[1], 10, 8)
	r2, _ := strconv.ParseUint(list[2], 10, 8)
	r3, _ := strconv.ParseUint(list[3], 10, 8)

	return color.RGBA{uint8(r0), uint8(r1), uint8(r2), uint8(r3)}
}

// Color2Str color.RGBA转成字符串
//
//	@param color
//	@return string
func Color2Str(color color.RGBA) string {
	r := fmt.Sprint(color.R)
	g := fmt.Sprint(color.G)
	b := fmt.Sprint(color.B)
	a := fmt.Sprint(color.A)
	return r + "," + g + "," + b + "," + a
}

// GetImageByWidthHeight 获取指定宽高的图片
//
//	@param p
//	@param w
//	@param h
//	@return string
func GetImageByWidthHeight(p string, w int, h int) string {

	if runtime.GOOS == "windows" {
		p = strings.ReplaceAll(p, "\\", "/")
	}
	t := strings.Split(p, ".")
	t[len(t)-2] = t[len(t)-2] + "_" + fmt.Sprintf("%d", w) + "_" + fmt.Sprintf("%d", h)

	p_ := strings.Join(t, ".")

	_, err := os.Stat(p_)
	if os.IsNotExist(err) {
		ImgResize(p, p_, uint(w), uint(h))
	}
	return p_
}

// ChangePngColor png图片文字颜色转换
//
//	@param o
//	@param s
//	@param oldColor
//	@param newColor
func ChangePngColor(o string, s string, oldColor color.RGBA, newColor color.RGBA) {
	// 加载原始图片
	input, _ := os.Open(o)
	defer input.Close()
	img, _ := png.Decode(input)

	// 颜色转换配置
	config := ColorConvertConfig{
		OldColor:       oldColor, // 要替换的颜色
		NewColor:       newColor, // 新颜色
		ColorTolerance: 30,       // 颜色匹配容差
		AntiAliasLevel: 2,        // 抗锯齿等级 (1-4)
	}

	// 执行转换与抗锯齿处理
	result := ConvertColorWithAA(img, config)

	// 保存结果
	output, _ := os.Create(s)
	png.Encode(output, result)
	output.Close()
}

// ColorConvertConfig 配置参数结构体
type ColorConvertConfig struct {
	OldColor       color.RGBA
	NewColor       color.RGBA
	ColorTolerance uint8
	AntiAliasLevel int
}

// ConvertColorWithAA 核心处理函数
//
//	@param src
//	@param cfg
//	@return *image.RGBA
func ConvertColorWithAA(src image.Image, cfg ColorConvertConfig) *image.RGBA {
	// Step 1: 超采样提升精度
	supersampled := imaging.Resize(src,
		src.Bounds().Dx()*cfg.AntiAliasLevel,
		src.Bounds().Dy()*cfg.AntiAliasLevel,
		imaging.Lanczos)

	// Step 2: 在超采样图像上执行颜色替换
	replaced := replaceColor(supersampled, cfg)

	// Step 3: 智能模糊抗锯齿
	blurred := imaging.Blur(replaced, float64(cfg.AntiAliasLevel)*0.8)

	// Step 4: 下采样恢复尺寸
	downsampled := imaging.Resize(blurred,
		src.Bounds().Dx(),
		src.Bounds().Dy(),
		imaging.MitchellNetravali)

	// Step 5: 边缘锐化补偿
	g := gift.New(
		gift.UnsharpMask(1.0, 1.0, 0), // 参数：sigma, amount, threshold
	)

	// 应用滤镜
	dst := image.NewRGBA(g.Bounds(downsampled.Bounds()))
	g.Draw(dst, downsampled)

	return dst
}

// replaceColor 颜色替换（带边缘过渡）
//
//	@param img
//	@param cfg
//	@return *image.RGBA
func replaceColor(img image.Image, cfg ColorConvertConfig) *image.RGBA {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			original := img.At(x, y)
			r, g, b, a := original.RGBA()

			// 计算颜色相似度
			similarity := colorSimilarity(original, cfg.OldColor, cfg.ColorTolerance)

			// 混合新旧颜色（保留透明度）
			newR := mixColor(uint8(r>>8), cfg.NewColor.R, similarity)
			newG := mixColor(uint8(g>>8), cfg.NewColor.G, similarity)
			newB := mixColor(uint8(b>>8), cfg.NewColor.B, similarity)

			dst.SetRGBA(x, y, color.RGBA{
				R: newR,
				G: newG,
				B: newB,
				A: uint8(a >> 8),
			})
		}
	}
	return dst
}

// mixColor 颜色混合函数（处理过渡边缘）
//
//	@param old
//	@param new
//	@param ratio
//	@return uint8
func mixColor(old, new uint8, ratio float64) uint8 {
	return uint8(float64(old)*(1-ratio) + float64(new)*ratio)
}

// colorSimilarity 改进型颜色相似度计算
//
//	@param c1
//	@param c2
//	@param tolerance
//	@return float64
func colorSimilarity(c1, c2 color.Color, tolerance uint8) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	// 转换为Lab色彩空间（更符合人眼感知）
	lab1 := rgbToLab(uint8(r1>>8), uint8(g1>>8), uint8(b1>>8))
	lab2 := rgbToLab(uint8(r2>>8), uint8(g2>>8), uint8(b2>>8))

	// 计算Delta E距离
	deltaL := lab1.L - lab2.L
	deltaA := lab1.A - lab2.A
	deltaB := lab1.B - lab2.B
	distance := math.Sqrt(deltaL*deltaL + deltaA*deltaA + deltaB*deltaB)

	// 将距离转换为相似度比例（0-1）
	maxDistance := 100.0 // Lab色彩空间理论最大距离约180
	similarity := 1.0 - math.Min(distance/maxDistance, 1.0)

	// 应用容差阈值
	if similarity < float64(tolerance)/255.0 {
		return 0.0
	}
	return math.Pow(similarity, 0.5) // 非线性响应曲线
}

// RGB转Lab色彩空间（简化版）
type Lab struct{ L, A, B float64 }

// rgbToLab
//
//	@param r
//	@param g
//	@param b
//	@return Lab
func rgbToLab(r, g, b uint8) Lab {
	// 此处应实现完整转换，为简化代码此处使用近似公式
	return Lab{
		L: 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b),
		A: float64(r - g),
		B: float64(g - b),
	}
}

// getDirFiles 获取指定目录下的文件内容
//
//	@param path
//	@return []string
func GetDirFiles(path string) []string {
	dir, err := os.Open(path)
	if err != nil {
		return []string{}
	}
	defer dir.Close()

	// 遍历文件夹中的文件和子文件夹
	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return []string{}
	}
	r := []string{}
	for _, file := range files {
		if !file.IsDir() {
			r = append(r, file.Name())
		}
	}
	return r
}

// TestProcessWaterMark 测试生成水印
func TestProcessWaterMark() {
	dir := "./test"
	list := GetDirFiles(dir)
	tplId := "1" // 模板id
	for _, file := range list {
		path := dir + "/" + file
		save := strings.ReplaceAll(path, dir, "./tmp/watermark")
		ProcessWaterMark(tplId, path, save)
	}
}

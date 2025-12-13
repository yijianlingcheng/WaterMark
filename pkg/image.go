package pkg

import (
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

// 加载图片.
func LoadImageWithDecode(path string) (image.Image, EError) {
	rio, err := os.Open(path)
	if err != nil {
		errmsg := path + ":文件打开失败:" + err.Error()

		return nil, NewErrors(FILE_NOT_OPEN_ERROR, errmsg)
	}
	filetype, eErr := GetFileType(rio)
	if HasError(eErr) {
		return nil, eErr
	}
	rio.Close()

	io, err := os.Open(path)
	if err != nil {
		errmsg := path + ":文件打开失败:" + err.Error()

		return nil, NewErrors(FILE_NOT_OPEN_ERROR, errmsg)
	}
	defer io.Close()

	var img image.Image
	switch filetype {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(io)
	case "image/png":
		img, err = png.Decode(io)
	default:
		errmsg := path + ":文件不是支持的格式"

		return nil, NewErrors(IMAGE_NO_SUPPORT_ERROR, errmsg)
	}
	// 判断是否decode成功
	if err != nil {
		errmsg := path + "image.Decode 失败:" + err.Error()

		return nil, NewErrors(IMAGE_DECODE_ERROR, errmsg)
	}

	return img, NoError
}

// 获取文件类型.
func GetFileType(io *os.File) (string, EError) {
	buff := make([]byte, 512)
	_, err := io.Read(buff)
	if err != nil {
		errmsg := "文件读取失败:" + err.Error()

		return "", NewErrors(FILE_NOT_READ_ERROR, errmsg)
	}

	return http.DetectContentType(buff), NoError // 根据http库获取文件类型
}

// 对指定图片文件生成指定宽高的图片.
func GenerateImageByWidthHeight(img image.Image, w, h int) image.Image {
	return resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
}

// 保存jpeg图片.
func SaveJpeg(path string, img image.Image, q int) EError {
	io, err := os.Create(path)
	if err != nil {
		return NewErrors(FILE_NOT_OPEN_ERROR, path+":文件创建失败:"+err.Error())
	}
	defer io.Close()

	encodeErr := jpeg.Encode(io, img, &jpeg.Options{
		Quality: q,
	})
	if encodeErr != nil {
		return ImageJpegSaveError
	}

	return NoError
}

// imageToRGBA 将图像转换为RGBA格式.
func ImageToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	// 使用并行处理提高转换速度
	numGoroutines := runtime.NumCPU()
	wg := sync.WaitGroup{}

	height := bounds.Dy()
	rowsPerGoroutine := height / numGoroutines

	for i := range numGoroutines {
		wg.Add(1)
		go func(startRow, endRow int) {
			defer wg.Done()
			for y := startRow; y < endRow; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					rgba.Set(x, y, img.At(x, y))
				}
			}
		}(bounds.Min.Y+i*rowsPerGoroutine,
			bounds.Min.Y+min((i+1)*rowsPerGoroutine, height))
	}

	wg.Wait()

	return rgba
}

// 并行顺时针旋转90度.
func Rotate90(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// 创建新的图像，宽高互换
	newImg := image.NewRGBA(image.Rect(0, 0, height, width))

	// 使用并行处理
	numGoroutines := runtime.NumCPU()
	wg := sync.WaitGroup{}

	rowsPerGoroutine := height / numGoroutines
	// 协程数量限制
	procs := runtime.GOMAXPROCS(0)
	ch := make(chan struct{}, max(5, procs))
	for i := range numGoroutines {
		ch <- struct{}{}
		wg.Add(1)
		go func(startRow, endRow int) {
			defer wg.Done()
			for y := startRow; y < endRow; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					// 计算新坐标：(x, y) -> (height-1-y, x)
					newX := height - 1 - (y - bounds.Min.Y)
					newY := x - bounds.Min.X
					newImg.Set(newX, newY, img.RGBAAt(x, y))
				}
			}
			<-ch
		}(bounds.Min.Y+i*rowsPerGoroutine,
			bounds.Min.Y+min((i+1)*rowsPerGoroutine, height))
	}

	wg.Wait()
	img = nil

	return newImg
}

// 并行旋转180度.
func Rotate180(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// 创建新的图像
	newImg := image.NewRGBA(bounds)

	// 使用并行处理
	numGoroutines := runtime.NumCPU()
	wg := sync.WaitGroup{}

	rowsPerGoroutine := height / numGoroutines

	// 协程数量限制
	procs := runtime.GOMAXPROCS(0)
	ch := make(chan struct{}, max(5, procs))
	for i := range numGoroutines {
		ch <- struct{}{}
		wg.Add(1)
		go func(startRow, endRow int) {
			defer wg.Done()
			for y := startRow; y < endRow; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					// 计算新坐标：(x, y) -> (width-1-x, height-1-y)
					newX := width - 1 - (x - bounds.Min.X)
					newY := height - 1 - (y - bounds.Min.Y)
					newImg.Set(newX, newY, img.RGBAAt(x, y))
				}
			}
			<-ch
		}(bounds.Min.Y+i*rowsPerGoroutine,
			bounds.Min.Y+min((i+1)*rowsPerGoroutine, height))
	}

	wg.Wait()
	img = nil

	return newImg
}

// 并行顺时针旋转270度（或逆时针旋转90度）.
func Rotate270(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// 创建新的图像，宽高互换
	newImg := image.NewRGBA(image.Rect(0, 0, height, width))

	// 使用并行处理
	numGoroutines := runtime.NumCPU()
	wg := sync.WaitGroup{}

	rowsPerGoroutine := height / numGoroutines
	// 协程数量限制
	procs := runtime.GOMAXPROCS(0)
	ch := make(chan struct{}, max(5, procs))
	for i := range numGoroutines {
		ch <- struct{}{}
		wg.Add(1)
		go func(startRow, endRow int) {
			defer wg.Done()
			for y := startRow; y < endRow; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					// 计算新坐标：(x, y) -> (y, width-1-x)
					newX := y - bounds.Min.Y
					newY := width - 1 - (x - bounds.Min.X)
					newImg.Set(newX, newY, img.RGBAAt(x, y))
				}
			}
			<-ch
		}(bounds.Min.Y+i*rowsPerGoroutine,
			bounds.Min.Y+min((i+1)*rowsPerGoroutine, height))
	}

	wg.Wait()
	img = nil

	return newImg
}

// 图片旋转.
func ImageRotate(orientation int, image image.Image) image.Image {
	switch orientation {
	case 90:
		image = imaging.Rotate90(image)
	case 180:
		image = imaging.Rotate180(image)
	case 270:
		image = imaging.Rotate270(image)
	}

	return image
}

package pkg

import (
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"

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
//
//nolint:gosec
func GenerateImageByWidthHeight(img image.Image, w, h int) image.Image {
	if img == nil {
		return nil
	}

	return resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
}

// 图片旋转.
func ImageRotate(orientation int, image image.Image) image.Image {
	if image == nil {
		return nil
	}
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

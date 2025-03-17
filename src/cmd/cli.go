package cmd

import (
	"WaterMark/src/exif"
	"WaterMark/src/paths"
	"runtime"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// InitToolPath 初始化工具目录
func InitToolPath() {
	if runtime.GOOS == "windows" {
		exiftoolBinary = paths.GetPwdPath(exiftoolBinary)
	}
}

// RunExifTool 运行exiftool工具,执行的时候通过hideCmdWindow隐藏cmd窗口
//
//	@param imgPath 需要获取exif信息的图片地址
//	@return exif.Exif
//	@return error
func RunExifTool(imgPath string) (exif.Exif, error) {
	et := NewExifTool(imgPath)
	return et.getExif()
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

// convertByte2Str 转换编码
//
//	@param b 需要转换的字符数组
//	@param c 指定的编码
//	@return string
func convertByte2Str(b []byte, c Charset) string {
	var s string
	switch c {
	case GB18030:
		d, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(b)
		s = string(d)
	case UTF8:
		fallthrough
	default:
		s = string(b)
	}
	return s
}

package cmd

import (
	"WaterMark/src/exif"
	"os"
	"runtime"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

var pwd string //当前路径

// init,设置pwd
func init() {
	pwd_, _ := os.Getwd()
	pwd = strings.ReplaceAll(pwd_, "\\", "/")
	if runtime.GOOS == "windows" {
		exiftoolBinary = pwd + exiftoolBinary
	}
}

// GetPwdPath 返回内部路径
//
//	@param p 文件路径
//	@return string
func GetPwdPath(p string) string {
	return pwd + p
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

package cmd

import (
	"WaterMark/src/exif"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// extraTags exiftool工具的标签,-TAG标签指定exiftool返回的字段内容
var extraTags = []string{"TAG", "Make", "Model", "CreateDate", "LensModel", "ExposureTime", "FNumber", "ISO", "FocalLength", "XResolution", "YResolution", "MechanicalShutterCount", "ShutterCount"}

// extractArgs 指定exiftool工具返回json格式数据
var extractArgs = []string{"-j"}

// Exiftool 结构体
type Exiftool struct {

	// exiftoolBinPath exiftool的可执行文件路径
	exiftoolBinPath string

	// filepath 需要获取exif信息的源文件路径
	filepath string

	// runTags 运行时的标签参数
	runTags []string
}

// NewExifTool 传入源文件路径,获取一个ExifTool结构体
//
//	@param p 源文件路径
//	@return *Exiftool
func NewExifTool(p string) *Exiftool {
	et := &Exiftool{
		exiftoolBinPath: exiftoolBinary,
		filepath:        p,
	}
	args := slices.Clone(extractArgs)
	if len(extraTags) > 0 {
		for _, v := range extraTags {
			args = append(args, "-"+v)
		}
	}
	args = append(args, p)
	et.runTags = args
	return et
}

// getExif 获取源文件的exif信息
//
//	@return exif.Exif
func (et *Exiftool) getExif() (exif.Exif, error) {

	args := et.exiftoolBinPath + " " + strings.Join(et.runTags, " ")
	r, err := cmdRun(args)
	if len(err) > 0 {
		return exif.Exif{}, errors.New(err)
	}
	return analysisResults(r), nil
}

// analysisResults 分析exiftool的执行结果,并将其序列化到Exif中
//
//	@param s cmd命令的执行结果
//	@return exif.Exif
func analysisResults(s string) exif.Exif {
	var maps []exif.Exif
	json.Unmarshal([]byte(s), &maps) // exiftool -j 返回的内容是json数组,因此取index=0作为结果
	item := maps[0]
	item.FNumberStr = strconv.FormatFloat(item.FNumber, 'f', 1, 64) //将光圈大小保留一位小数并转换为字符串
	item.ISOStr = fmt.Sprintf("%d", item.ISO)                       //将ISO转换为字符串
	return item
}

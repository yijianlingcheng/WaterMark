package api

import (
	"WaterMark/src/cmd"
	"WaterMark/src/exif"
)

// GetShutter 获取快门次数
//
//	@param imgPath 需要查看快门次数的图片路径
//	@return exif.Exif
func GetShutter(imgPath string, fullPathFlag bool) exif.Exif {
	if fullPathFlag {
		m, _ := cmd.RunExifTool(imgPath)
		return exif.Getshutter(m)
	}
	m, _ := cmd.RunExifTool(cmd.GetPwdPath(imgPath))
	return exif.Getshutter(m)
}

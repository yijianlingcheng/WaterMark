package exift

import (
	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/internal"
	"WaterMark/pkg"
)

// 初始化exiftool工具.
func InitExiftool() (*exiftool.Exiftool, pkg.EError) {
	et, err := exiftool.NewExiftool(exiftool.SetExiftoolBinaryPath(internal.GetExiftoolPath()))
	if err != nil {
		internal.Log.Error("初始化exiftool工具失败: " + err.Error())

		return nil, pkg.ExiftoolInitError
	}

	return et, pkg.NoError
}

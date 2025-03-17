package paths

import (
	"path"
	"runtime"
)

// Config
var rootPath string

// initRootPath 初始化项目根路径
func InitRootPath() {
	_, _filename_, _, _ := runtime.Caller(0)
	rootPath = path.Dir(path.Dir(path.Dir(_filename_)))
}

// getRooPath 获取项目根路径
//
//	@return string
func GetRooPath() string {
	return rootPath
}

// GetPwdPath 返回绝对路径
//
//	@param p 文件路径
//	@return string
func GetPwdPath(p string) string {
	return rootPath + p
}

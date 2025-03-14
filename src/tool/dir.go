package tool

import "os"

// Exists 判断路径是否存在
//	@param path
//	@return bool
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	return !os.IsNotExist(err)
}

// CreateDir 创建目录
//	@param path
func CreateDir(path string) {
	os.MkdirAll(path, os.ModePerm)
}

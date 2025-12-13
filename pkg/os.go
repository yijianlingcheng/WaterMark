package pkg

import (
	"os"
	"runtime"
)

const (

	// darwin系统.
	Darwin = "darwin"

	// windows系统.
	Window = "windows"
)

// 检查当前运行环境是否为window.
func IsWindows() bool {
	return runtime.GOOS == Window
}

// 获取指定文件夹下面的全部文件(不支持获取文件夹中下级文件夹中的文件).
func GetDirFiles(directory string) ([]string, EError) {
	list := make([]string, 0, 100)
	files, err := os.ReadDir(directory) // 读取目录中的文件
	if err != nil {
		return list, NewErrors(FILE_NOT_READ_ERROR, directory+":读取文件夹失败:"+err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		list = append(list, file.Name())
	}

	return list, NoError
}

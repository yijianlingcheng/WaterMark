//nolint:gosec
package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// 获取字符串md5.
func GetStrMD5(str string) string {
	h := md5.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		// 计算字符串md5都能报错,直接panic算了
		panic(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

// 获取文件md5.
func GetFileMD5(filePath string) (string, EError) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", NewErrors(FILE_NOT_OPEN_ERROR, filePath+":failed to open file: "+err.Error())
	}
	defer file.Close()

	// 创建一个md5哈希对象
	hasher := md5.New()

	// 将文件内容读入哈希对象
	if _, err = io.Copy(hasher, file); err != nil {
		return "", NewErrors(FILE_NOT_READ_ERROR, filePath+": failed to read file: "+err.Error())
	}

	// 计算MD5值并转换为十六进制字符串
	md5Sum := hasher.Sum(nil)

	return hex.EncodeToString(md5Sum), NoError
}

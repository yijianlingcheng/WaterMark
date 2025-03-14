package tool

import (
	"WaterMark/src/exif"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

// ExifToJson Exif 转string
//
//	@param param
//	@return string
func ExifToJson(param exif.Exif) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

// StrMD5 MD5 计算md5
//
//	@param str
//	@return string
func StrMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

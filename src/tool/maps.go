package tool

import (
	"WaterMark/src/exif"
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

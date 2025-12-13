package native

import (
	"strings"

	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/layout"
	"WaterMark/pkg"
)

// 将模板中的字符串转换为实际展示的字符串.
func changeText2ExifContent(exif exiftool.FileMetadata, str string) string {
	// 切割字符串
	strList := strings.Split(str, ",")
	exifList := make([]string, 0)

	for i := range strList {
		// 内部字符串使用#进行分割
		subs := strings.SplitSeq(strList[i], "#")
		for sub := range subs {
			// 默认使用原有字符串
			value := sub
			// 从exif中取值
			if v, ok := exif.Fields[sub]; ok {
				value = pkg.AnyToString(v)
			}
			value = changeExifShowStr(sub, value, exif)
			exifList = append(exifList, value)
		}
	}

	return strings.Join(exifList, "")
}

// 转换展示的exif信息.
func changeExifShowStr(sub, str string, exif exiftool.FileMetadata) string {
	// 经纬度特殊判断
	if strings.Contains(sub, GPS_OR_DATETIME) {
		gpsOrTime := pkg.AnyToString(exif.Fields[DATE_TIME_ORIGINAL])
		if gpsOrTime == "" {
			gpsOrTime = GPS_OR_DATETIME
		}
		str = layout.GetGPSOrDefault(
			pkg.AnyToString(exif.Fields[GPS_POSITION]),
			gpsOrTime,
		)

		return str
	}

	// 经纬度特殊判断
	if strings.Contains(sub, GPS_POSITION) {
		// 说明没有获取到GPS 信息
		if sub == str {
			return str
		}
		str = layout.GpsFormat(str)

		return str
	}

	// 优化展示
	if str != "" && strings.Contains(sub, FOCAL_LENGTH) {
		str = changeFocalLength(str)

		return str
	}
	if str != "" && strings.Contains(sub, DATE_TIME_ORIGINAL) {
		str = changeDataTime(str)

		return str
	}

	return str
}

// 优化焦段展示,去除多余展示的.0.
func changeFocalLength(str string) string {
	// 去除空格
	str = strings.Replace(str, " ", "", 1)
	t := strings.Split(str, ".")
	if len(t) == 1 {
		return str
	}

	return t[0] + "mm"
}

// 优化拍摄时间展示,去除多余展示的:00.
func changeDataTime(str string) string {
	t := str[len(str)-3:]
	if t == ":00" {
		return str[:len(str)-3]
	}

	return str
}

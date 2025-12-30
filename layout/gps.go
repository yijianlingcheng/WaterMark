package layout

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 获取GPS格式化之后的字符串.
func GpsFormat(gps string) string {
	if gps == "" {
		return gps
	}
	str := strings.Split(gps, ", ")

	if len(str) == 1 {
		return parseGPSInfo(str[0]) + " "
	}

	return parseGPSInfo(str[0]) + " " + parseGPSInfo(str[1])
}

// 获取GPS格式化之后的字符串,如果gps为空则展示默认值.
func GetGPSOrDefault(gps, other string) string {
	s := GpsFormat(gps)
	if s != "" && s != " " {
		return s
	}

	return other
}

// 解析GPS信息.
func parseGPSInfo(str string) string {
	if str == "" {
		return ""
	}
	// 正则表达式匹配格式：数字 deg 数字 ' 数字 " 方向
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*deg\s*(\d+(?:\.\d+)?)\s*'\s*(\d+(?:\.\d+)?)\s*"\s*([NSEW])`)
	matches := re.FindStringSubmatch(str)

	if matches == nil {
		return ""
	}
	// 解析各部分
	degrees, _ := strconv.ParseFloat(matches[1], 64)
	minutes, _ := strconv.ParseFloat(matches[2], 64)
	direction := matches[4]

	// 格式化输出
	result := fmt.Sprintf("%d°%02d′%s", int(degrees), int(minutes), direction)

	return result
}

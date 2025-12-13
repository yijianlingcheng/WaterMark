package pkg

import (
	"regexp"
	"strconv"
	"strings"
)

// 获取照片需要旋转的角度.
func GetOrientation(o string) int {
	if strings.Contains(o, "Horizontal") || strings.Contains(o, "normal") {
		return 0
	}
	if !strings.Contains(o, "Rotate") {
		return 0
	}
	if strings.Contains(o, "CCW") {
		// 逆时针旋转
		r := extractNumbers(o)
		if len(r) > 0 {
			return r[0]
		}

		return 0
	}
	if strings.Contains(o, "CW") {
		// 顺时针旋转
		r := extractNumbers(o)
		if len(r) > 0 {
			return 360 - r[0] // 转换成逆时针旋转角度
		}
	}

	return 0
}

// 使用正则表达式匹配所有数字.
func extractNumbers(s string) []int {
	// 使用正则表达式匹配所有数字
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(s, -1)

	// 将匹配到的字符串转换为整数
	numbers := make([]int, 0)
	for _, match := range matches {
		num, err := strconv.Atoi(match)
		if err == nil {
			numbers = append(numbers, num)
		}
	}

	return numbers
}

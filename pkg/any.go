package pkg

import (
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
)

// any类型转string.
func AnyToString(value any) string {
	var result string
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		result = v
	case int, int8, int16, int32, int64:
		result = strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		result = strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
	case float32:
		result = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		result = strconv.FormatFloat(v, 'f', -1, 64)
	case []byte:
		result = string(v)
	default:
		jsonValue, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		result = string(jsonValue)
	}

	return result
}

// 判断字符串是否存在某个字符串数组中.
func In(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}

	return false
}

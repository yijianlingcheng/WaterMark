package controller

import (
	"errors"
	"sync"
)

var imageInfoCache sync.Map

// imageInfoCacheSet
//
//	@param key
//	@param value
func imageInfoCacheSet(key string, value map[string]string) {
	imageInfoCache.Store(key, value)
}

// imageInfoCacheGet
//
//	@param key
//	@return map[string]string
//	@return error
func imageInfoCacheGet(key string) (map[string]string, error) {
	v, ok := imageInfoCache.Load(key)
	if ok {
		return v.(map[string]string), nil
	}
	return map[string]string{}, errors.New("获取失败")
}

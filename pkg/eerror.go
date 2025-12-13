package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
)

// 自定义错误描述.
type EError struct {
	Error error `json:"errmsg"`
	Code  int   `json:"code"`
}

var (
	// 没有发生错误.
	NoError = EError{
		Code:  0,
		Error: errors.New("no error"),
	}

	InternalError = EError{
		Code:  INTERNAL_ERROR,
		Error: errors.New("发生内部错误"),
	}

	// exiftool工具不存在.
	ExiftoolNotExistError = EError{
		Code:  EXIFTOOL_NOTEXIST_ERROR,
		Error: errors.New("exiftool工具不存在,请检查是否安装"),
	}

	// exiftool工具init失败.
	ExiftoolInitError = EError{
		Code:  EXIFTOOL_INIT_ERROR,
		Error: errors.New("exiftool工具init失败,请检查环境"),
	}

	// exiftool工具获取exif信息失败.
	ExiftoolImageError = EError{
		Code:  EXIFTOOL_IMAGE_EXIF_ERROR,
		Error: errors.New("exiftool工具获取图片exif信息失败"),
	}

	// 从缓存中获取的exif数据类型断言失败.
	ExiftoolCacheTypeError = EError{
		Code:  EXIFTOOL_IMAGE_EXIF_CACHE_ERROR,
		Error: errors.New("exif cache缓存数据类型断言失败"),
	}

	// 从缓存的图片解码数据类型断言失败.
	ImageDecodeCacheTypeError = EError{
		Code:  IMAGE_DECODE_CACHE_ERROR,
		Error: errors.New("缓存的图片解码数据类型断言失败"),
	}

	// 从缓存中获取的RGBA数据类型断言失败.
	ImageRGBACacheTypeError = EError{
		Code:  IMAGE_RGBA_CACHE_ERROR,
		Error: errors.New("缓存的RGBA数据类型断言失败"),
	}

	// 缓存中获取的字体对象类型断言失败.
	ImageTextCacheTypeError = EError{
		Code:  IMAGE_TEXT_FONT_CACHE_ERROR,
		Error: errors.New("缓存中获取的字体对象类型断言失败"),
	}

	// logo图片文件没有找到.
	ImageLogoNotFindError = EError{
		Code:  IMAGE_LOGO_NOT_FIND_ERROR,
		Error: errors.New("logo图片文件没有找到"),
	}

	// logo图片重置尺寸错误.
	ImageLogoResizeError = EError{
		Code:  IMAGE_LOGO_RESIZE_ERROR,
		Error: errors.New("logo图片重置尺寸错误"),
	}

	// jpeg图片保存失败.
	ImageJpegSaveError = EError{
		Code:  IMAGE_JPEG_SAVE_ERROR,
		Error: errors.New("jpeg图片保存失败"),
	}

	// 布局类型查找失败.
	LayoutNotFindError = EError{
		Code:  LAYOUT_TYPE_NOT_FIND_ERROR,
		Error: errors.New("布局类型查找失败"),
	}
)

// 是否正常.
func IsOk(e EError) bool {
	return e.Code == NO_ERROR
}

// 是否发生错误.
func HasError(e EError) bool {
	return e.Code > NO_ERROR
}

// 打印方法.
func (e EError) String() string {
	return fmt.Sprintf("{\"Code\":%d,\"Error\":%q}", e.Code, e.Error.Error())
}

// 自定义json格式.
func (e EError) MarshalJSON() ([]byte, error) {
	type Alias EError

	return json.Marshal(&struct {
		Alias
		Error string `json:"errmsg"`
	}{
		Error: e.Error.Error(),
		Alias: Alias(e),
	})
}

// 返回一个指定的错误.
func NewErrors(code int, msg string) EError {
	return EError{
		Code:  code,
		Error: errors.New(msg),
	}
}

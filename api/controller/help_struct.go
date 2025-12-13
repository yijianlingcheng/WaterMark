package controller

import "WaterMark/engine/frame"

type (
	ErrorInfo struct {
		Errmsg string `json:"errmsg"`
		Code   int    `json:"code"`
	}

	ExifInfoSuccess struct {
		File      string `json:"file"`
		Equipment string `json:"器材"`
		CMode     string `json:"模式"`
		Params    string `json:"参数"`
		Focal     string `json:"焦距"`
		Color     string `json:"色彩"`
		Time      string `json:"时间"`
		Shutter   string `json:"快门次数"`

		Orientation   string `json:"Orientation"`
		ImageWidth    string `json:"ImageWidth"`
		ImageHeight   string `json:"ImageHeight"`
		Make          string `json:"Make"`
		Model         string `json:"Model"`
		LensModel     string `json:"LensModel"`
		FocalLength   string `json:"FocalLength"`
		FNumber       string `json:"FNumber"`
		ExposureTime  string `json:"ExposureTime"`
		ISO           string `json:"ISO"`
		FileName      string `json:"FileName"`
		ImageSize     string `json:"ImageSize"`
		ImageDataSize string `json:"ImageDataSize"`

		Errmsg string `json:"errmsg"`
		Code   int    `json:"code"`
	}

	NoError struct {
		Errmsg string `json:"errmsg"`
		Code   int    `json:"code"`
	}

	ImportInfo struct {
		Errmsg string   `json:"errmsg"`
		Files  []string `json:"files"`
		Code   int      `json:"code"`
	}

	TemplatesInfo struct {
		List   map[string]string `json:"list"`
		Errmsg string            `json:"errmsg"`
		Code   int               `json:"code"`
	}

	ExifAndBorderInfo struct {
		Errmsg string          `json:"errmsg"`
		Exif   ExifInfoSuccess `json:"exif"`
		Text   []string        `json:"text"`
		Size   frame.PhotoSize `json:"size"`
		Code   int             `json:"code"`
	}

	Message struct {
		Errmsg string   `json:"errmsg"`
		List   []string `json:"list"`
		Code   int      `json:"code"`
	}
)

package controller

import "WaterMark/engine/frame"

type (
	// ErrorInfo 用于统一返回错误信息。
	// Errmsg: 错误描述，供前端或调用方展示。
	// Code: 状态码，非 0 表示出现错误，具体值可对应不同错误类型。
	ErrorInfo struct {
		Errmsg string `json:"errmsg"`
		Code   int    `json:"code"`
	}

	// ExifInfoSuccess 表示从图像文件中提取的 EXIF 与文件元数据。
	// File: 文件路径或名称。
	// Equipment: 器材（拍摄设备/相机，中文字段显示为“器材”）。
	// CMode: 拍摄模式（中文字段“模式”）。
	// Params: 拍摄参数汇总（中文字段“参数”）。
	// Focal: 焦距（中文字段“焦距”）。
	// Color: 色彩信息或色彩空间（中文字段“色彩”）。
	// Time: 拍摄时间（中文字段“时间”）。
	// Shutter: 快门次数或快门信息（中文字段“快门次数”）。
	// Orientation: 图像方向（Orientation EXIF 标签）。
	// ImageWidth: 图像宽度（像素）。
	// ImageHeight: 图像高度（像素）。
	// Make: 相机制造商（Make EXIF 标签）。
	// Model: 相机型号（Model EXIF 标签）。
	// LensModel: 镜头型号（LensModel EXIF 标签）。
	// FocalLength: 焦距字符串或数值（FocalLength EXIF 标签）。
	// FNumber: 光圈值（FNumber EXIF 标签）。
	// ExposureTime: 曝光时间（ExposureTime EXIF 标签）。
	// ISO: 感光度（ISO EXIF 标签）。
	// FileName: 文件名。
	// ImageSize: 图像尺寸描述（例如 "4000x3000"）。
	// ImageDataSize: 图像数据大小（字节或可读格式）。
	// Errmsg: 错误信息，当提取失败或部分字段不可用时提供错误描述。
	// Code: 状态码，表示操作结果（例如 0 表示成功，非零表示错误）。
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

	// NoError 用于表示操作成功，无错误信息。
	// Errmsg: 错误描述，成功时通常为空字符串。
	// Code: 状态码，0 表示成功。
	NoError struct {
		Errmsg string `json:"errmsg"`
		Code   int    `json:"code"`
	}

	// ImportInfo 表示导入文件操作的结果。
	// Errmsg: 错误信息，导入失败时提供具体原因。
	// Files: 成功导入的文件路径列表。
	// Code: 状态码，0 表示全部导入成功，非零表示存在失败。
	ImportInfo struct {
		Errmsg string   `json:"errmsg"`
		Files  []string `json:"files"`
		Code   int      `json:"code"`
	}

	// TemplatesInfo 用于返回模板列表及相关状态。
	// List: 模板名称到模板内容的映射。
	// Errmsg: 错误信息，获取模板列表失败时提供原因。
	// Code: 状态码，0 表示获取成功。
	TemplatesInfo struct {
		List   map[string]string `json:"list"`
		Errmsg string            `json:"errmsg"`
		Code   int               `json:"code"`
	}

	// ExifAndBorderInfo 封装了带边框处理所需的全部信息。
	// Errmsg: 错误信息，处理失败时提供原因。
	// Exif: 提取到的 EXIF 信息。
	// Text: 需要渲染到边框上的文字列表。
	// Size: 目标照片尺寸（含边框）。
	// Code: 状态码，0 表示处理成功。
	ExifAndBorderInfo struct {
		Errmsg string          `json:"errmsg"`
		Exif   ExifInfoSuccess `json:"exif"`
		Text   []string        `json:"text"`
		Size   frame.PhotoSize `json:"size"`
		Code   int             `json:"code"`
	}

	// Message 用于返回通用列表数据及状态。
	// Errmsg: 错误信息，获取列表失败时提供原因。
	// List: 字符串列表数据。
	// Code: 状态码，0 表示获取成功。
	Message struct {
		Errmsg string   `json:"errmsg"`
		List   []string `json:"list"`
		Code   int      `json:"code"`
	}
)

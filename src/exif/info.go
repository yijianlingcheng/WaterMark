package exif

// 程序使用到的exif信息
// "Make":                   "相机生产商",
// "Model":                  "相机型号",
// "CreateDate":             "拍照时间",
// "LensModel":              "镜头信息",
// "ExposureTime":           "曝光时间",
// "FNumber":                "光圈大小",
// "ISO":                    "ISO",
// "FocalLength":            "镜头焦距",
// "XResolution":            "x dpi,水平分辨率",
// "YResolution":            "y dpi,垂直分辨率",
// "MechanicalShutterCount": "机器快门数",
// "ShutterCount":           "快门数",

// Exif 简要的exif信息,这只是程序需要使用到的字段
type Exif struct {

	// ISO ISO 大小
	ISO int

	// XResolution 水平分辨率
	XResolution int

	// YResolution 垂直分辨率
	YResolution int

	// ShutterCount 快门数
	ShutterCount int

	// MechanicalShutterCount 机器快门数
	MechanicalShutterCount int

	// FNumber 光圈大小
	FNumber float64

	// FNumberStr 光圈大小,字符串类型,为了方便使用格式化为字符串
	FNumberStr string

	// FocalLength 镜头焦距
	FocalLength string

	// Make 相机制造商
	Make string

	// Model 相机信息
	Model string

	// CreateDate 拍照时间
	CreateDate string

	// LensModel 镜头信息
	LensModel string

	// ExposureTime 快门时间
	ExposureTime string

	// ISOStr ISO字符串
	ISOStr string
}

// Getshutter 获取快门信息
//	@param m
//	@return Exif
func Getshutter(m Exif) Exif {
	// MechanicalShutterCount 或者 ShutterCount为0代表exif没有对应的快门数字段
	return m
}

// GetWaterMarkInfo 获取水印需要的相关信息
//	@param m
//	@return Exif
func GetWaterMarkInfo(m Exif) Exif {
	return m
}

// CoverImgResolution 复写图片的XResolution,YResolution
//	@param path 图片路径
//	@param x 水平分辨率
//	@param y 垂直分辨率
func CoverImgResolution(path string, x int, y int) {

}

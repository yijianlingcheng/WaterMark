package native

const (
	// 分隔符默认颜色.
	SEPARATOR_COLOR = "203,203,201,255"

	// 默认颜色.
	COLOR = "255,255,255,255"

	// 使用gps或者时间,gps信息不存在则使用时间.
	GPS_OR_DATETIME = "GPS_OR_DATETIME"

	// 原始时间.
	DATE_TIME_ORIGINAL = "DateTimeOriginal"

	// gps信息.
	GPS_POSITION = "GPSPosition"

	// 焦段.
	FOCAL_LENGTH = "FocalLength"

	// 类型:边框.
	PHOTO_TYPE_BORDER = "border"
)

// 文字内容列表.
var textWordsList = [4]string{
	"text_one_content",
	"text_two_content",
	"text_three_content",
	"text_four_content",
}

package pkg

const (
	// 没错.
	NO_ERROR = 0

	// 文件不存在.
	FILE_NOT_EXIST_ERROR = 1000001

	// 文件打开失败.
	FILE_NOT_OPEN_ERROR = 1000002

	// 文件读取失败.
	FILE_NOT_READ_ERROR = 1000003

	// exiftool工具不存在.
	EXIFTOOL_NOTEXIST_ERROR = 2000001

	// exiftool工具初始化失败.
	EXIFTOOL_INIT_ERROR = 2000002

	// exiftool工具获取exif失败.
	EXIFTOOL_IMAGE_EXIF_ERROR = 2000003

	// 从缓存中获取的exif 缓存类型断言失败.
	EXIFTOOL_IMAGE_EXIF_CACHE_ERROR = 2000004

	// csv文件创建失败.
	CSV_CREATE_ERROR = 3000001

	// csv文件写入头信息失败.
	CSV_WRITE_HEADER_ERROR = 3000002

	// csv文件写入数据失败.
	CSV_WRITE_DATA_ERROR = 3000003

	// 图片解码失败.
	IMAGE_DECODE_ERROR = 4000001

	// 图片格式不支持.
	IMAGE_NO_SUPPORT_ERROR = 4000002

	// 从缓存中获取的图片解码信息类型断言失败.
	IMAGE_DECODE_CACHE_ERROR = 4000003

	// 从缓存中获取的RGBA对象类型断言失败.
	IMAGE_RGBA_CACHE_ERROR = 4000004

	// 从缓存中获取的字体对象类型断言失败.
	IMAGE_TEXT_FONT_CACHE_ERROR = 4000005

	// 画笔绘制失败.
	IMAGE_TEXT_DRAW_TXT_ERROR = 4000006

	// logo文件查找失败.
	IMAGE_LOGO_NOT_FIND_ERROR = 4000007

	// logo图片重置尺寸错误.
	IMAGE_LOGO_RESIZE_ERROR = 4000008

	// jpeg图片保存失败.
	IMAGE_JPEG_SAVE_ERROR = 4000009

	// cmd 执行命令失败.
	CMD_COMMAND_RUN_ERROR = 5000001

	// 布局类型查找失败.
	LAYOUT_TYPE_NOT_FIND_ERROR = 6000001

	// 内部错误.
	INTERNAL_ERROR = 9000001

	// 请求参数错误.
	REQUEST_PARAM_ERROR = 9000002
)

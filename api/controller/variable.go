package controller

var (
	// 请求参数名：file 字段（用于上传或指定源文件）.
	paramQueryFile = "file"
	// 请求参数名：save 字段（用于指定导出/保存路径）.
	paramQuerySave = "save"
	// 请求参数名：layout 字段（用于指定布局标识）.
	paramQueryLayout = "layout"
	// 请求参数名：preview_layout 字段（预览时使用的布局参数）.
	paramQueryPreviewLayout = "preview_layout"

	// 错误/提示信息：file 参数为空.
	paramFileIsEmpty = "file参数为空"
	// 错误/提示信息：请求的 file 文件不存在.
	paramFileIsNotExist = "file请求的文件不存在"
	// 错误/提示信息：layout 参数为空.
	paramLayoutIsEmpty = "layout参数为空"
	// 错误/提示信息：未选择保存路径.
	paramSaveIsEmpty = "未选择保存的路径"
	// 错误/提示信息：导出存放图片的路径不存在.
	paramSaveIsNotExist = "导出存放图片的路径不存在"

	// 导出进度通道（缓冲 100）：用于异步汇报导出进度.
	export_Progress_Chan = make(chan string, 100)

	// 默认缩放比例（用于导出时未指定比例时的默认值）.
	defaultRatio = 5
	// 最大缩放比例（用于导出时指定比例时的最大值）.
	maxRatio = 8
	// 最大图片尺寸（用于导出时指定比例时的最大值）.
	maxImageSize = 8000
)

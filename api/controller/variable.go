package controller

var (
	// file字段.
	paramQueryFile = "file"
	// save字段.
	paramQuerySave = "save"
	// layout字段.
	paramQueryLayout = "layout"
	// 预览时修改的布局参数.
	paramQueryPrevireLayout = "preview_layout"

	paramFileIsEmpty = "file参数为空"

	paramFileIsNotExist = "file请求的文件不存在"

	paramLayoutIsEmpty = "layout参数为空"

	paramSaveIsEmpty = "未选择保存的路径"

	paramSaveIsNotExist = "导出存放图片的路径不存在"

	// 导出进度条.
	export_Progress_Chan = make(chan string, 100)
)

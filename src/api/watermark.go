package api

import (
	"WaterMark/src/cmd"
	"WaterMark/src/images"
	"WaterMark/src/tool"
	"strings"
	"sync"
)

// getTplListType 获取模板列表
//
//	@return map
func getTplListType() map[string]string {
	r := map[string]string{}
	for _, v := range images.GetTemplates() {
		r[v.ID] = v.Type
	}
	return r
}

// getImageWaterMarkPreview 获取水印图片预览
//
//	@param tid
//	@param path
//	@param color
//	@param onlyBottomBorder
//	@return map
func getImageWaterMarkPreview(tid string, path string, color string, onlyBottomBorder bool) map[string]string {
	e := images.NewExternal()
	e.WithBoderColor(color).WithOnlyBottomFlag(onlyBottomBorder).WithPath(path).WithTid(tid)
	r := images.CreatePreviewWaterMark(e)
	r["SaveImgPath"] = cmd.GetPwdPath(strings.TrimLeft(r["SaveImgPath"], "."))
	return r
}

// addPreviewTask 添加水印图片任务
//
//	@param imageStrs
func addPreviewTask(imageStrs string) {
	tpls := images.GetTemplates()
	imageList := strings.Split(imageStrs, ",")
	g := images.NewWorkerlimit(images.MAX_WORKER_NUM)
	var wg = sync.WaitGroup{}
	for i := range imageList {
		wg.Add(1)
		imgPath := tool.ReplaceDir(imageList[i])
		goFunc := func() {
			e := images.NewExternal().WithDefaultBoderColor().WithPath(imgPath).WithTid(tpls[0].ID)
			images.CreatePreviewWaterMark(e)
			wg.Done()
		}
		g.Worker(goFunc)
	}
	wg.Wait()
}

// addPreviewTask 添加水印图片任务
//
//	@param imageStrs
func addImageResizeTask(imageStrs string) {
	imageList := strings.Split(imageStrs, ",")
	g := images.NewWorkerlimit(images.MAX_WORKER_NUM)
	var wg = sync.WaitGroup{}
	for i := range imageList {
		wg.Add(1)
		imgPath := tool.ReplaceDir(imageList[i])
		goFunc := func() {
			e := images.NewExternal().WithSmallPreviewPath(imgPath)
			images.CreateSmallPreview(e)
			wg.Done()
		}
		g.Worker(goFunc)
	}
	wg.Wait()
}

// getSmallPreviewPath 获取小图预览地址
//
//	@param imgPath
//	@return string
func getSmallPreviewPath(imgPath string) string {
	e := images.NewExternal().WithSmallPreviewPath(imgPath)
	return e.SavePath
}

// downloadFile 下载文件
//
//	@param sourcePath
//	@param previewPath
//	@return error
func downloadFile(sourcePath, previewPath string) error {
	sourcePath = tool.ReplaceDir(sourcePath)
	previewPath = tool.ReplaceDir(previewPath)

	t := strings.Split(sourcePath, ".")
	t[len(t)-2] = t[len(t)-2] + "_watermark"

	newSourcePath := strings.Join(t, ".")
	return tool.CopyFile(previewPath, newSourcePath, 4*1024)
}

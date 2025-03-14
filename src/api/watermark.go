package api

import (
	"WaterMark/src/cmd"
	"WaterMark/src/images"
	"runtime"
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

	r := images.GetPreviewWaterMark(e)
	r["SaveImgPath"] = cmd.GetPwdPath(strings.TrimLeft(r["SaveImgPath"], "."))
	return r
}

// addPreviewTask 添加水印图片任务
//
//	@param imageStrs
func addPreviewTask(imageStrs string) {
	tpls := images.GetTemplates()
	imageList := strings.Split(imageStrs, ",")

	limit := len(imageList)
	g := images.NewWorkerlimit(images.MAX_WORKER_NUM)
	var wg = sync.WaitGroup{}
	for i := range limit {
		wg.Add(1)
		imgPath := imageList[i]
		if runtime.GOOS == "windows" {
			imgPath = strings.ReplaceAll(imgPath, "\\", "/")
		}
		goFunc := func() {
			e := images.NewExternal()
			e.WithDefaultBoderColor().WithPath(imgPath).WithTid(tpls[0].ID)
			images.GetPreviewWaterMark(e)
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

	limit := len(imageList)
	g := images.NewWorkerlimit(images.MAX_WORKER_NUM)
	var wg = sync.WaitGroup{}
	for i := range limit {
		wg.Add(1)
		imgPath := imageList[i]
		if runtime.GOOS == "windows" {
			imgPath = strings.ReplaceAll(imgPath, "\\", "/")
		}
		goFunc := func() {
			e := images.NewExternal().WithSmallPreviewPath(imgPath)
			images.CeateSmallPreview(e)
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

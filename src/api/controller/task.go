package controller

import (
	"WaterMark/src/images"
	"WaterMark/src/tool"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// AddPreviewTask 添加水印图片任务
//
//	@param ctx
func AddPreviewTask(ctx *gin.Context) {
	c := Container(ctx)

	imageStrs := c.PostForm("images")
	tpls := images.GetTemplates()
	imageList := strings.Split(imageStrs, ",")

	// 添加任务,异步执行直接返回成功
	go previewWaterMarkTask(imageList, tpls[0].ID)
	c.JSON(gin.H{"message": "添加水印图片任务完成"})
}

// previewWaterMarkTask 图片水印预览生成任务
//
//	@param imageList 源图片列表
//	@param tid 模板id
func previewWaterMarkTask(imageList []string, tid string) {
	g := images.NewWorkerlimit(images.MAX_WORKER_NUM)
	var wg = sync.WaitGroup{}
	for i := range imageList {
		wg.Add(1)
		imgPath := tool.ReplaceDir(imageList[i])
		goFunc := func() {
			e := images.NewExternal().WithDefaultBoderColor().WithPath(imgPath).WithTid(tid)
			r := images.CreatePreviewWaterMark(e)
			// 新增缓存
			imageInfoCacheSet(imgPath, r)
			wg.Done()
		}
		g.Worker(goFunc)
	}
	wg.Wait()
}

// AddImageResizeTask 添加压缩图片任务(预览使用)
//
//	@param ctx
func AddImageResizeTask(ctx *gin.Context) {
	c := Container(ctx)

	imageStrs := c.PostForm("images")
	list := strings.Split(imageStrs, ",")

	if len(list) > 0 {
		imageResizeTask(list)
	}

	c.JSON(gin.H{"message": "指定图片的预览小图生成完成"})
}

// ImageResizeTask 图片压缩任务
//
//	@param list 需要压缩的图片
func imageResizeTask(list []string) {
	g := images.NewWorkerlimit(images.MAX_WORKER_NUM)
	var wg = sync.WaitGroup{}
	for i := range list {
		wg.Add(1)
		imgPath := tool.ReplaceDir(list[i])
		goFunc := func() {
			e := images.NewExternal().WithSmallPreviewPath(imgPath)
			images.CreateSmallPreview(e)
			wg.Done()
		}
		g.Worker(goFunc)
	}
	wg.Wait()
}

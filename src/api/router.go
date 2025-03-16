package api

import (
	"WaterMark/src/api/controller"

	"github.com/gin-gonic/gin"
)

// loadRouters 加载路由
//
//	@param router
func loadRouters(router *gin.Engine) {

	// 获取快门次数
	router.POST("/server/getShutterByFile", controller.GetShutterByFile)

	// 获取图片预览
	router.GET("/server/getImagePreview", controller.GetImagePreview)

	// 获取图片模板列表类型
	router.GET("/server/getTplListType", controller.GetTplListType)

	// 获取图片水印预览
	router.POST("/server/getImageWaterMarkPreview", controller.GetImageWaterMarkPreview)

	// 添加图片水印预览异步任务
	router.POST("/server/addPreviewTask", controller.AddPreviewTask)

	// 添加图片压缩任务
	router.POST("/server/addImageResizeTask", controller.AddImageResizeTask)

	// 根据原图返回预览小图
	router.GET("/server/imagePreviewSmall", controller.ImagePreviewSmall)

	// 下载文件
	router.POST("/server/downloadFile", controller.DownloadFile)

	// 转换文件路径
	router.POST("/server/changeImagePath", controller.ChangeImagePath)
}

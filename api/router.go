package api

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"WaterMark/api/controller"
	_ "WaterMark/docs"
)

func loadRouters(router *gin.Engine) {
	// Swagger UI路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动接口
	message := router.Group("message")
	// 获取启动的info消息
	message.GET("info", controller.InfoMessage)

	// 查看菜单接口
	view := router.Group("view")
	// 展示照片
	view.GET("showImage", controller.ShowImage)

	// 获取照片的exif信息
	view.POST("getImagesExifInfo", controller.GetPhotosExifInfo)
	// 导出照片exif信息到指定文件
	view.POST("exifInfoExportv2", controller.ExifInfoExportBySaveFile)

	// 边框菜单接口
	frame := router.Group("frame")
	// 导入图片资源
	frame.POST("importPhotoFiles", controller.ImportPhotoFiles)
	// 照片生成边框
	frame.POST("showPhotoFrame", controller.ShowPhotoFrame)
	// 创建导出任务
	frame.POST("createExportTask", controller.CreateExportTask)
	// 获取导出进度
	frame.GET("getExportProgress", controller.GetExportProgress)
	// 获取exif信息与图片信息
	frame.POST("getExifAndBorderInfo", controller.GetPhotoExifAndBorderInfo)
	// 重新加载logo文件夹下图片
	frame.POST("reloadLogoImages", controller.ReloadLogoImages)
	// 重新加载边框模板文件
	frame.POST("reloadFrameTemplate", controller.ReloadFrameTemplate)
	// 获取边框模板信息
	frame.GET("getFrameTemplateInfo", controller.GetFrameTemplateInfo)
}

// 注册不需要记录日志的API接口
// 部分接口返回的数据实在太多,为了保证日志文件不超大,因此这里面的接口不记录日志.
func getNoLogApis() []string {
	return []string{
		"/view/showImage",
		"/swagger/",
		"/frame/showPhotoFrame",
	}
}

package api

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"WaterMark/internal"
)

// 启动api server.
func ServerStart() {
	// 启动gin
	// 设置release 模式
	gin.SetMode(gin.ReleaseMode)
	// 关闭控制台颜色
	gin.DisableConsoleColor()
	// release模式判断
	if internal.ISRelease() {
		// 关闭gin框架的默认输出
		gin.DefaultWriter = io.Discard
	}
	// router
	router := gin.Default()
	// 注册中间件
	for _, w := range getMiddlewareList() {
		router.Use(w)
	}
	// 加载路由
	loadRouters(router)
	// 监听地址并运行
	err := router.Run(viper.GetString("server.address"))
	if err != nil {
		internal.Log.Panic("server start error:" + err.Error())
	}
}

package controller

import (
	"WaterMark/src/logs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server Controller控制器
type Server struct {
	ctx *gin.Context
}

// Container server
//
//	@param ctx
//	@return *Controller
func Container(ctx *gin.Context) *Server {
	server := &Server{
		ctx: ctx,
	}
	if ctx.Request.Method != "GET" {
		ctx.Request.ParseMultipartForm(32 << 20) // 32 << 20 是ParseMultipartForm的默认参数
		logs.API.Println(ctx.Request.Method, ctx.Request.URL, "请求参数", ctx.Request.Form)
	}
	return server
}

// json 返回json数据
//
//	@param obj
func (c *Server) JSON(obj any) {
	logs.API.Println(c.ctx.Request.Method, c.ctx.Request.URL, "接口返回", obj)
	c.ctx.JSON(http.StatusOK, obj)
}

// PostForm
//
//	@param key
//	@return string
func (c *Server) PostForm(key string) string {
	return c.ctx.PostForm(key)
}

// DefaultPostForm
//
//	@param key
//	@param defaultv
//	@return string
func (c *Server) DefaultPostForm(key string, defaultv string) string {
	return c.ctx.DefaultPostForm(key, defaultv)
}

// Query
//
//	@param key
//	@return string
func (c *Server) Query(key string) string {
	return c.ctx.Query(key)
}

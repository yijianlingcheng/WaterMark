package controller

import (
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
	return server
}

// json 返回json数据
//
//	@param obj
func (c *Server) JSON(obj any) {
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

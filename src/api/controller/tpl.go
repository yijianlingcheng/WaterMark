package controller

import (
	"WaterMark/src/images"

	"github.com/gin-gonic/gin"
)

// GetTplListType 获取模板列表
//
//	@param ctx
func GetTplListType(ctx *gin.Context) {
	c := Container(ctx)
	r := map[string]string{}
	for _, v := range images.GetTemplates() {
		r[v.ID] = v.Type
	}
	c.JSON(r)
}

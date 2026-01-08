package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// @Summary 生成二维码
// @Description 根据指定参数，生成二维码
// @Tags Qrcode
// @Produce image/png
// @Param content formData string true "二维码内容"
// @Router /qrcode/generate [post]
// @Failure 400 {object} ErrorInfo "错误信息".
func GenerateQRCode(ctx *gin.Context) {
	content := ctx.PostForm(paramQueryContent)
	if content == "" {
		ctx.JSON(400, requestParamError(paramContentIsEmpty))
		return
	}
	var pngs []byte
	pngs, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 返回二维码图片
	ctx.Writer.Header().Set("Content-Type", "image/png")
	ctx.Writer.Write(pngs)
}

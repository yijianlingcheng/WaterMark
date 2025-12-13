package controller

import (
	"time"

	"github.com/gin-gonic/gin"

	"WaterMark/message"
)

// 消息key.
const MESSAGE_KEY = "message"

// @Summary 获取启动页输出
// @Description 获取启动页输出
// @Tags message
// @Produce json
// @Router /message/info [get]
// @Success 200 {object} Message "成功信息".
func InfoMessage(ctx *gin.Context) {
	// 设置头
	setSSEHeader(ctx)

	infos := make([]string, 0)

msg:
	for {
		select {
		case info, ok := <-message.Info_Messge_Chan:
			if ok {
				infos = append(infos, info)
				if message.HasSendSuccess(info) {
					break msg
				}
			}
		default:
			time.Sleep(10 * time.Microsecond)
		}
	}

	ctx.JSON(200, Message{
		List: infos,
	})
}

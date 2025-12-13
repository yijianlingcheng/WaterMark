package message

import (
	"log"
	"sync"
	"time"

	"WaterMark/pkg"
)

var (
	// 错误信息传递管道.需要弹窗提示错误的消息通过此管道传递.
	Error_Messge_Chan = make(chan string, 100)

	// info 消息传递管道.需要弹窗提示的消息通过此管道进行传递.
	Info_Messge_Chan = make(chan string, 100)

	// 启动成功消息.
	start_success_message = "启动完成"

	// 关闭标致.
	info_channel_close_flag = false

	// 互斥锁.
	msgMtx sync.Mutex
)

// 发送启动成功消息.
func SendStartSuccess() {
	SendInfoMsg(start_success_message)
}

// 检查是否发送成功消息.
func HasSendSuccess(msg string) bool {
	return msg == start_success_message
}

// 发送错误消息.
func SendErrorMsg(errStr string) {
	Error_Messge_Chan <- errStr
}

// 发送错误或者成功消息.
func SendErrorOrInfo(err pkg.EError, success string) {
	if pkg.HasError(err) {
		SendErrorMsg(err.String())
		CloseInfoChannel()

		return
	}
	SendInfoMsg(success)
}

// 关闭info通道标识.
func CloseInfoChannel() {
	msgMtx.Lock()
	info_channel_close_flag = true
	msgMtx.Unlock()
}

// 发送信息.
func SendInfoMsg(info string) {
	msgMtx.Lock()
	defer msgMtx.Unlock()

	// 关闭标致,不写入info数据
	if info_channel_close_flag {
		return
	}

	Info_Messge_Chan <- info
}

// 关闭管道.
func Close() {
	close(Info_Messge_Chan)
	close(Error_Messge_Chan)
}

// 处理channel.
func ApiDebug() {
	for {
		select {
		case info, ok := <-Info_Messge_Chan:
			if ok {
				log.Println(info)
			}
		case err, ok := <-Error_Messge_Chan:
			if ok {
				panic("Error:" + err)
			}
		default:
			time.Sleep(10 * time.Microsecond)
		}
	}
}

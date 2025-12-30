package internal

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"WaterMark/pkg"
)

var Log = logrus.New()

// 初始化日志工具.
func initLogConfig() pkg.EError {
	// 在日志中记录对应的调用方法
	Log.SetReportCaller(true)

	// JSON话日志格式
	Log.SetFormatter(&logrus.JSONFormatter{})

	// 设置日志存放位置,存放在app.log文件中,这个不需要定制化存放路径,程序中直接写死即可
	file, _ := os.OpenFile(GetLogPath("app.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	Log.Out = file

	if ISRelease() {
		Log.SetOutput(file)
		// 设置info级别模式
		Log.SetLevel(logrus.InfoLevel)
	} else {
		std := os.Stdout
		Log.SetOutput(io.MultiWriter(file, std))
		// 设置debug级别模式
		Log.SetLevel(logrus.TraceLevel)
	}

	return pkg.NoError
}

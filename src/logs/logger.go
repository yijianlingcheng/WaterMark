package logs

import (
	"WaterMark/src/paths"
	"WaterMark/src/tool"
	"log"
	"os"
)

var (
	Errors *log.Logger
	Info   *log.Logger
	API    *log.Logger
)

const (
	LOG_DIR   = "/logs"
	RUN_LOG   = LOG_DIR + "/run.log"
	API_LOG   = LOG_DIR + "/api.log"
	FILE_FlAG = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	FILE_MODE = 0666
	LOG_FLAG  = log.Ldate | log.Ltime | log.Lshortfile
)

// init
func InitLog() {
	// 创建目录
	dir := paths.GetPwdPath(LOG_DIR)
	if !tool.Exists(dir) {
		tool.CreateDir(dir)
	}

	// 运行日志
	runLog := paths.GetPwdPath(RUN_LOG)
	file, err := os.OpenFile(runLog, FILE_FlAG, FILE_MODE)
	if err == nil {
		Errors = log.New(file, "ERROR: ", LOG_FLAG)
		Info = log.New(file, "INFO: ", LOG_FLAG)
	}

	// api日志
	apiLog := paths.GetPwdPath(API_LOG)
	file, err = os.OpenFile(apiLog, FILE_FlAG, FILE_MODE)
	if err == nil {
		API = log.New(file, "API:", log.Ldate|log.Ltime)
	}
}

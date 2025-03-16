package logs

import (
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
	LOG_DIR   = "./logs"
	RUN_LOG   = LOG_DIR + "./run.log"
	API_LOG   = LOG_DIR + "./api.log"
	FILE_FlAG = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	FILE_MODE = 0666
	LOG_FLAG  = log.Ldate | log.Ltime | log.Lshortfile
)

// init
func init() {
	// 创建目录
	if !tool.Exists(LOG_DIR) {
		tool.CreateDir(LOG_DIR)
	}
	// 其他日志
	file, err := os.OpenFile(RUN_LOG, FILE_FlAG, FILE_MODE)
	if err == nil {
		Errors = log.New(file, "ERROR: ", LOG_FLAG)
		Info = log.New(file, "INFO: ", LOG_FLAG)
	}

	// api
	file, err = os.OpenFile(API_LOG, FILE_FlAG, FILE_MODE)
	if err == nil {
		API = log.New(file, "API:", log.Ldate|log.Ltime)
	}
}

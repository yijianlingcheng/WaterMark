package cmd

import (
	. "WaterMark/src/logs"
	"bytes"
	"os/exec"
	"syscall"
)

// exiftoolBinary exiftool可执行文件路径
var exiftoolBinary = "/tools/exiftool.exe"

// hideCmdWindow 隐藏win 下cmd的执行界面
//
//	@param cmd
func hideCmdWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true, // 隐藏调用命令的执行框
	}
}

// winRun 运行window cmd命令
//
//	@param args 具体的参数
//	@return string cmd标准输出结果
//	@return string cmd错误输出结果
func cmdRun(args string) (string, string) {
	cmd := exec.Command("cmd.exe", "/C", args)

	hideCmdWindow(cmd)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()

	if err != nil {
		errStr = err.Error() + errStr + pwd + ":" + args
		data := []byte(errStr)
		Errors.Println(convertByte2Str(data, GB18030))
	}
	return outStr, errStr
}

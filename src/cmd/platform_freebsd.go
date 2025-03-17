package cmd

import (
	"WaterMark/src/logs"
	"bytes"
	"os/exec"
)

// exiftoolBinary exiftool可执行文件路径
var exiftoolBinary = "exiftool"

// Freebsd 运行 cmd命令
//
//	@param args 具体的参数
//	@return string cmd标准输出结果
//	@return string cmd错误输出结果
func cmdRun(args string) (string, string) {
	cmd := exec.Command(args)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	if err != nil {
		errStr = err.Error() + errStr + ":" + args
		data := []byte(errStr)
		logs.Errors.Println(convertByte2Str(data, GB18030))
	}
	return outStr, errStr
}

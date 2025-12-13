package cmd

import (
	"bytes"
	"context"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"

	"WaterMark/pkg"
)

const (
	UTF8    charset = "UTF-8"
	GB18030 charset = "GB18030"
)

type charset string

// 隐藏window下的cmd命令框.
func hideWindowCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true, // 隐藏调用命令的执行框
	}
}

// commandRun 运行window cmd命令.
func CommandRun(timeout time.Duration, args string) (string, pkg.EError) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "cmd.exe", "/C", args)

	hideWindowCmd(cmd)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()

	cmdErr := pkg.NoError
	if err != nil {
		errStr = changeToUTF8String(errStr, GB18030)
		errStr = err.Error() + errStr + ":" + args
		cmdErr = pkg.NewErrors(pkg.CMD_COMMAND_RUN_ERROR, errStr)
	}

	return outStr, cmdErr
}

// 将字节切片转换为指定编码的字符串.
func changeToUTF8String(bytes string, charset charset) string {
	var str string
	switch charset {
	case GB18030:
		decoder := simplifiedchinese.GB18030.NewDecoder()
		var err error
		str, err = decoder.String(bytes)
		if err != nil {
			return ""
		}
	case UTF8:
		str = bytes
	default:
		str = bytes
	}

	return str
}

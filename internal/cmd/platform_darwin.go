package cmd

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"

	"WaterMark/pkg"
)

// CommandRun darwin 运行 cmd命令
func CommandRun(timeout time.Duration, args string) (string, pkg.EError) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", args)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()

	cmdErr := pkg.NoError
	if err != nil {
		errStr = err.Error() + errStr + ":" + args
		cmdErr = pkg.NewErrors(pkg.CMD_COMMAND_RUN_ERROR, errStr)

	}

	return outStr, cmdErr
}

// CommandRun darwin 运行 cmd命令
func CommandRunWithArgs(timeout time.Duration, args []string) (string, pkg.EError) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", args[0], args[1])

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()

	cmdErr := pkg.NoError
	if err != nil {
		errStr = err.Error() + errStr + ":" + strings.Join(args, " ")
		cmdErr = pkg.NewErrors(pkg.CMD_COMMAND_RUN_ERROR, errStr)

	}

	return outStr, cmdErr
}

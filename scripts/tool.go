package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	dev     = "dev"
	api     = "api-dev"
	release = "release"

	getRootPathFunc = getRootPath
	getFilePathFunc = getFilePath
	printErrorFunc  = printError
)

func main() {
	appMode := flag.String("appMode", "dev", "App Mode")
	flag.Parse()

	changeAppModeToDev(*appMode)
	replaceAppVersion()
}

func getRootPath() string {
	rootPath, _ := os.Getwd()

	return rootPath
}

func getFilePath(p string) string {
	return getRootPathFunc() + "/" + p
}

// 转换主main.go中的AppMode为开发模式.
func changeAppModeToDev(mode string) {
	log.Println("----------------Set main.go AppMode " + mode + "---------------")

	if mode != dev && mode != api && mode != release {
		printErrorFunc(errors.New("appMode value error"))
	}

	filePath := getFilePathFunc("main.go")
	codeBytes, err := os.ReadFile(filePath)
	if err != nil {
		printErrorFunc(err)
	}

	codeStr := replaceMode(string(codeBytes), mode)
	err = os.WriteFile(filePath, []byte(codeStr), 0o600)
	if err != nil {
		printErrorFunc(err)
	}
	log.Println("----------------change main.go success---------------")
}

// 打印错误信息.
//
//nolint:revive
func printError(err error) {
	log.Println("---------------------error---------------------------", err)
	os.Exit(1)
}

// 替换mode.
func replaceMode(codeStr, mode string) string {
	modes := [3]string{
		"internal.SetAppMode(internal.APP_RELEASE)",
		"internal.SetAppMode(internal.APP_DEV)",
		"internal.SetAppMode(internal.APP_API_DEV)",
	}
	newMode := modes[0]
	switch mode {
	case dev:
		newMode = modes[1]
	case api:
		newMode = modes[2]
	}
	newCodeStr := codeStr
	for i := range modes {
		newCodeStr = strings.ReplaceAll(
			newCodeStr,
			modes[i],
			newMode,
		)
	}

	return newCodeStr
}

// 获取APP当前版本号.
func getAppVersion() string {
	filePath := getFilePathFunc("version")
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		printErrorFunc(errors.New(filePath + ":文件不存在"))
	}
	version := strings.ReplaceAll(string(fileContent), "APP_VERSION:", "")

	return strings.Trim(version, "\r\n")
}

// 替换版本说明文件中的APP版本号.
func replaceAppVersion() {
	log.Println("----------------replace code version ---------------")

	filePath := getFilePathFunc("frontend/src/A/aboutVersionView.html")
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		printErrorFunc(errors.New(filePath + ":文件不存在"))
	}

	re := regexp.MustCompile(`(?i)(<p>\s*版本\s*:\s*)([^<]*?)(\s*</p>)`)
	newContent := re.ReplaceAllStringFunc(string(fileContent), func(m string) string {
		sub := re.FindStringSubmatch(m)
		// sub[0] = 整个匹配, sub[1] = 前缀, sub[2] = 旧版本文本, sub[3] = 后缀
		if len(sub) < 4 {
			return m // 安全退回：若未按预期捕获则不修改
		}

		return sub[1] + getAppVersion() + sub[3]
	})
	successOut := "-------------replace code version success-------------"
	if newContent == string(fileContent) {
		log.Println(successOut)

		return
	}
	err = os.WriteFile(filePath, []byte(newContent), 0o600)
	if err != nil {
		printErrorFunc(errors.New(filePath + ":文件内容替换失败"))
	}

	log.Println(successOut)
}

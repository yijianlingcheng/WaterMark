package internal

import (
	"os"
	"regexp"
	"strings"
	"sync"

	"WaterMark/message"
)

// 确保只执行一次.
var once sync.Once

// 初始化整个程序需要使用的配置与资源.
func InitAppConfigsAndRes() {
	once.Do(func() {
		// 初始化项目根路径
		initRootPath()

		// 初始化程序需要使用到的各种文件夹
		message.SendInfoMsg("检查程序执行需要的目录")
		createAppDS(getAppAllRuntimeDS())
		message.SendInfoMsg("检查程序执行需要的目录通过")

		// 初始化配置
		message.SendInfoMsg("加载应用配置")
		err := initAppConfig()
		message.SendErrorOrInfo(err, "加载应用配置完成")

		// 初始化日志记录的配置
		message.SendInfoMsg("初始化日志记录的配置")
		initLogConfig()
		message.SendErrorOrInfo(err, "初始化日志记录的配置完成")

		// 初始化用户自定义配置
		message.SendInfoMsg("初始化用户自定义配置")
		err = initCustomizeConfig()
		message.SendErrorOrInfo(err, "初始化用户自定义配置完成")

		// 检查exiftool工具是否安装
		message.SendInfoMsg("检查exiftool工具是否安装")
		err = checkInstallExif()
		message.SendErrorOrInfo(err, "exiftool工具已安装")

		message.SendInfoMsg("检查ImageMagick工具是否安装")
		err = checkInstallImageMagick()
		message.SendErrorOrInfo(err, "ImageMagick工具已安装")

		// 释放字体文件
		message.SendInfoMsg("检查字体文件是否释放")
		err = checkFontFile()
		message.SendErrorOrInfo(err, "字体文件已存放")

		// 释放logo文件
		message.SendInfoMsg("检查logo是否释放")
		err = checkLogoFile()
		message.SendErrorOrInfo(err, "logo已存在")
		// 是否开启api调试
		if ISApiDebug() {
			go message.ApiDebug()

			return
		}
	})
}

// 获取APP当前版本号.
func GetAppVersion() string {
	filePath := GetRootPath() + "/version"
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		Log.Panic(filePath + ":文件不存在")
	}
	version := strings.ReplaceAll(string(fileContent), "APP_VERSION:", "")

	return strings.Trim(version, "\r\n")
}

// 替换版本说明文件中的APP版本号.
func ReplaceAppVersion() {
	if ISRelease() {
		return
	}
	filePath := GetRootPath() + "/frontend/src/A/aboutVersionView.html"
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		Log.Panic(filePath + ":文件不存在")
	}

	re := regexp.MustCompile(`(?i)(<p>\s*版本\s*:\s*)([^<]*?)(\s*</p>)`)
	newContent := re.ReplaceAllStringFunc(string(fileContent), func(m string) string {
		sub := re.FindStringSubmatch(m)
		// sub[0] = 整个匹配, sub[1] = 前缀, sub[2] = 旧版本文本, sub[3] = 后缀
		if len(sub) < 4 {
			return m // 安全退回：若未按预期捕获则不修改
		}

		return sub[1] + GetAppVersion() + sub[3]
	})
	if newContent == string(fileContent) {
		return
	}
	err = os.WriteFile(filePath, []byte(newContent), 0o600)
	if err != nil {
		Log.Panic(filePath + ":文件内容替换失败")
	}
}

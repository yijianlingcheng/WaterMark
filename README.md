# WaterMark
一个为照片添加边框与水印功能的程序

### 简要说明
 1. 使用Go开发(https://golang.google.cn/)
 2. 桌面程序使用wails构建(https://wails.io/zh-Hans/docs/introduction)
 3. 照片exif信息获取使用exiftool工具(https://exiftool.org/)
 4. 照片模糊模板使用ImageMagick添加阴影边框效果，[ImageMagick](https://imagemagick.org/)
 5. 后端接口服务采用gin框架(https://gin-gonic.com/zh-cn/)
 6. Go exiftool库fork https://github.com/barasher/go-exiftool 并进行了部分修改
 7. 文字水印使用Alibaba-PuHuiTi-Bold.ttf,Alibaba-PuHuiTi-Light.ttf字体(https://alibabafont.taobao.com/)
 8. 目前MacOS,Win10,Win11
 9. 源码请访问github(https://github.com/yijianlingcheng/WaterMark)

### Windows exiftool
 1. Windows系统下,程序已经内置打包exiftool工具,运行时会自动解压到指定的路径
### MacOS exiftool
 1. MacOS安装完成exiftool后需将exiftool加入环境变量

### Windows ImageMagick
 1. Windows系统下,程序已经内置打包ImageMagick工具,运行时会自动解压到指定的路径
### MacOS ImageMagick
 1. MacOS安装完成ImageMagick后需将ImageMagick加入环境变量


### 项目开发与调试 
 1. 先安装Go并配置环境(Go1.18+)
 2. 参照https://wails.io/zh-Hans/docs/gettingstarted/installation 安装wails
 3. 执行wails doctor检查依赖
 4. clone 代码完成之后 cd 进入代码目录
 5. API接口使用`swag` 自动生成api文档,如果新增或者修改了接口,需要执行`swag init`重新生成文档
 6. 代码统一使用`LF`作为换行符, IDE需要安装插件识别项目中的`.editorconfig`文件
 7. IDE需要识别项目中的`.gitignore`文件,避免提交不必要的文件
 8. 代码使用`golangci-lint`作为代码检查工具,IDE需要安装插件识别项目中的`.golangci.yaml`文件
 9. 修改代码之后需要执行`golangci-lint run`检查代码质量
 10. Win系统下执行tools.bat脚本进行调试开发
 11. Mac系统下执行tools.sh脚本进行调试开发
 12. 开发新版本时需要修改version文件中`APP_VERSION`对应的版本号

#### 说明
 1. APP_DEV app 开发模式,此模式用于 ` wails dev `
 2. APP_API_DEV api 开发模式,此模式用于 `go run main.go ` 启动api接口服务与必要的工具
 3. APP_RELEASE 打包模式,此模式用于 ` wails build `
 4. 模式修改位于` main.go` 文件的第一行:`SetAppMode()`
 5. 部分图片、配置文件、字体文件通过`go-bindata`打包进代码中
 6. 打包命令如下:
    1. 杂项 `go-bindata -pkg=assetmixfs   ./configs ./fonts ./logos`
    2. exiftool工具`go-bindata -pkg=assetexiffs ./exiftool`
    3. imagemagick工具`go-bindata -pkg=assetmagickfs ./magick`


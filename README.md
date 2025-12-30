# WaterMark
一个为照片添加边框与水印功能的程序

### 简要说明
 1. 使用[Go](https://golang.google.cn/)开发
 2. 桌面程序使用[wails构建](https://wails.io/zh-Hans/docs/introduction)
 3. 照片exif信息获取使用[exiftool工具](https://exiftool.org/)
 4. 照片模糊模板使用ImageMagick添加阴影边框效果,[ImageMagick](https://imagemagick.org/)
 5. 后端接口服务采用[gin框架](https://gin-gonic.com/zh-cn/)
 6. Go exiftool库fork ([go-exiftool](https://github.com/barasher/go-exiftool)) 并进行了部分修改
 7. 修改之后的[go-exiftool](https://github.com/yijianlingcheng/go-exiftool)地址
 8. 文字水印使用[Alibaba-PuHuiTi-Bold.ttf,Alibaba-PuHuiTi-Light.ttf字体](https://alibabafont.taobao.com/)
 9. 目前MacOS,Win10,Win11
 10. 源码请访问[github代码](https://github.com/yijianlingcheng/WaterMark)
   
### 下载与使用
 1. Win10与Win11下使用本程序
    1. 推荐官方下载,访问[下载地址](https://github.com/yijianlingcheng/WaterMark/releases),下载最新的release版本的exe文件
    2. 如果无法打开上面的官方下载地址,可以尝试复制下载地址在迅雷中打开,使用迅雷进行下载
    3. 可以在分享的百度网盘进行exe文件下载(网盘下载请注意分别是否是官方发布的程序,避免计算机中毒)
 2. 源码下载
    1. 访问[github代码](https://github.com/yijianlingcheng/WaterMark)下载源码
    2. 环境搭建参考后续项目开发与调试的内容

### 使用说明
 1. 使用说明
 2. Mac下使用说明
 3. 常见问题

### Windows exiftool
 1. Windows系统下,程序已经内置打包exiftool工具,运行时会自动解压到指定的路径
### MacOS exiftool
 1. MacOS安装完成exiftool后需将exiftool加入环境变量

### Windows ImageMagick
 1. Windows系统下,程序已经内置打包ImageMagick工具,运行时会自动解压到指定的路径
### MacOS ImageMagick
 1. MacOS安装完成ImageMagick后需将ImageMagick加入环境变量


### 项目开发与调试
#### 环境搭建 
 1. 先安装Go并配置环境(Go1.18+)
 2. 参照[wails官方文档](https://wails.io/zh-Hans/docs/gettingstarted/installation) 安装wails
 3. 安装wails完成之后,执行`wails doctor`检查依赖
 4. 安装swag工具
    1. 执行安装命令`go install github.com/swaggo/swag/cmd/swag@latest`
 5. 安装golangci-lint工具
    1. 执行安装命令`go install github.com/golangci/golangci-lint/cmd/golangci-lint`
 6. Mac下需要安装exiftool与ImageMagick并加入环境变量
    1. 安装exiftool工具
       1. 从官网([exiftool官方](https://exiftool.org/)) 下载对应MacOS exiftool.pkg并安装
       2. 也可以使用`brew install exiftool`通过brew命令安装
       3. 安装完成后执行`exiftool`命令检查是否加入环境变量
    2. 安装ImageMagick工具
       1. 打开[ImageMagick官网下载](https://imagemagick.org/script/download.php),选择源码下载并安装
       2. 也可以执行`brew install imagemagick`进行安装
       3. 安装完成之后执行`magick -version`命令检查是否安装成功,是否添加到了环境变量
       4. 如果安装的imagemagick版本为6x版本,对应的命令为`convert -version`
 7. Win下exiftool与ImageMagick
    1. Win下面程序将exiftool工具与ImageMagick打包进了可执行程序
    2. exiftool文件在项目的exiftool目录中
    3. 附带的exiftool工具版本为13.21
    4. ImageMagick文件在项目的magick目录中
    5. 附带的magick工具版本为 7.1.2-11 Q8 x64
 8. 代码使用`golangci-lint`作为代码检查工具,IDE需要安装插件识别项目中的`.golangci.yaml`文件
 9.  代码统一使用`LF`作为换行符, IDE需要安装插件识别项目中的`.editorconfig`文件

#### 开发与调试
 1.  clone 代码完成之后 cd 进入代码目录
 2.  开发新版本时需要修改version文件中`APP_VERSION`对应的版本号
 3.  程序提供了简易开发工具简化开发与调试流程
 4.  开发工具位于`scripts/`目录下
 5.  Window系统运行`tools.bat`
 6.  Mac系统运行`tools.sh`
 7.  tools工具提供以下几个命令
     1.  `tools.bat/tools.sh check`,执行此命令自动运行`golangci-lint`检查程序代码质量
     2.  `tools.bat/tools.sh api`,执行此命令自动运行`swag init`生成接口文档,自动执行`go run .\scripts\tool.go -appMode=api-dev`修改程序运行方式(参考appMode模式说明),自动将version文件中的版本号同步至版本展示页面,自动执行`go run .\main.go`
     3.  `tools.bat/tools.sh dev`,执行此命令自动运行`swag init`生成接口文档,自动执行`go run .\scripts\tool.go -appMode=dev`修改程序运行方式(参考appMode模式说明),自动将version文件中的版本号同步至版本展示页面,自动执行`wails dev`
     4.  `tools.bat/tools.sh build`,执行此命令自动运行`swag init`生成接口文档,自动执行`golangci-lint`检查程序代码质量,自动执行`go run .\scripts\tool.go -appMode=release`修改程序运行方式(参考appMode模式说明),自动将version文件中的版本号同步至版本展示页面,自动执行`wails build -clean -o "APP_VERSION"`进行打包
 8. 不使用`tools.bat/tools.sh`工具,进行开发与调试请参考以下步骤
    1. 独立调试接口
       1. 手动修改main.go中的SetAppMode为APP_API_DEV
       2. 手动执行`go run .\main.go`
    2. 整体调试
       1. 手动修改main.go中的SetAppMode为APP_DEV
       2. 手动执行`wails dev`
    3. 打包发布
       1. 手动修改main.go中的SetAppMode为APP_RELEASE
       2. 手动执行`swag init`生成文档
       3. 手动执行`golangci-lint`检查程序代码质量
       4. 修改代码保证`golangci-lint`结果为0
       5. 手动修改`./frontend/src/A/aboutVersionView.html`中的版本号
       6. 手动执行命令`wails build -clean -o "APP_VERSION"`

### 其他说明
 1. appMode模式,分别为api-dev,dev,release
 2. api-dev模式只启动exiftool工具与gin服务,不启动窗口程序.用于api接口调式
 3. dev模式,正常开发模式,启动exiftool工具,gin服务与窗口程序,用于调试整个程序
 4. release模式,打包模式,构建release包文件
 5. 相机logo图片、配置文件、字体文件、工具通过`go-bindata`打包进代码中
 6. 打包命令如下:
    1. 杂项 `go-bindata -pkg=assetmixfs   ./configs ./fonts ./logos`
    2. exiftool工具`go-bindata -pkg=assetexiffs ./exiftool`
    3. imagemagick工具`go-bindata -pkg=assetmagickfs ./magick`


### 简要说明
 1. 使用Go开发(https://golang.google.cn/)
 2. 桌面程序使用wails构建(https://wails.io/zh-Hans/docs/introduction)
 3. 照片exif信息获取使用exiftool工具(https://exiftool.org/)
 4. 后端接口服务采用gin框架(https://gin-gonic.com/zh-cn/)
 5. 页面布局使用Bootstrap v4(https://v4.bootcss.com/docs/getting-started/introduction/)
 6. 图片水印生成使用Go库，不需要额外安装扩展
 7. 文字水印使用Alibaba-PuHuiTi-Bold.ttf,Alibaba-PuHuiTi-Light.ttf字体(https://alibabafont.taobao.com/)
 8. 目前MacOS，Win10，Win11
 9. 源码请访问github(https://github.com/yijianlingcheng/WaterMark)

### windows exiftool
 1. windows系统下,exiftool工具需要放到tools目录下。需要自行去官网(https://exiftool.org/)下载
 2. 下载完成后进行解压,需要将exiftool-xx.xx_64下的文件全部放入tools目录
 3. 将tools目录下的exiftool(-k).exe重命名为exiftool.exe
 4. 一个正确的示例: 
    1. path 是项目路径,那么exiftool的路径应该如下
       1. path/tools/exiftool.exe
       2. path/tools/exiftool_files


### 项目路径
 1. 项目不支持中文路径。
 2. 将代码clone到本地之后,一个推荐的做法
    1. Windows下,路径为D:\WaterMark 或者E:\WaterMark

### 项目开发
 1. 先安装Go并配置环境(Go1.18+)
 2. 参照https://wails.io/zh-Hans/docs/gettingstarted/installation 安装wails
 3. 执行wails doctor检查依赖
 4. clone 代码完成之后 cd 进入代码目录
 5. 安装exiftool工具 -- window请参照上文将整个exiftool文件夹拷贝到代码目录下;MacOS安装完成exiftool后需将exiftool加入环境变量
 6. 执行 wails dev 进行调试开发
 7. 执行 wails build 进行构建
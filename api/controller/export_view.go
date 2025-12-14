package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yijianlingcheng/go-exiftool"

	"WaterMark/engine/frame"
	"WaterMark/internal"
	"WaterMark/layout"
	"WaterMark/pkg"
)

// 检查基本的参数.
func checkCreateExportTask(ctx *gin.Context) pkg.EError {
	save := ctx.PostForm(paramQuerySave)
	if save == "" {
		return pkg.NewErrors(pkg.REQUEST_PARAM_ERROR, paramSaveIsEmpty)
	}
	if !internal.PathExists(save) {
		return pkg.NewErrors(pkg.REQUEST_PARAM_ERROR, paramSaveIsNotExist)
	}
	file := ctx.PostForm(paramQueryFile)
	if file == "" {
		return pkg.NewErrors(pkg.REQUEST_PARAM_ERROR, paramFileIsEmpty)
	}
	layoutStr := ctx.PostForm(paramQueryLayout)
	if layoutStr == "" {
		return pkg.NewErrors(pkg.REQUEST_PARAM_ERROR, paramLayoutIsEmpty)
	}

	return pkg.NoError
}

// @Summary 创建导出任务
// @Description 对指定照片创建导出任务,异步执行导出
// @Tags Frame
// @Produce json
// @Param save formData string true "导出文件存放的路径"
// @Param file formData string true "照片路径;多个文件,分割"
// @Param layout formData string true "布局信息,JSON字符串:必须包含frame_name字段"
// @Param preview_layout formData string true "布局信息,边框预览时调整保存的参数"
// @Router /frame/createExportTask [post]
// @Success 200 {object} NoError "成功信息".
// @Failure 400 {object} ErrorInfo "错误信息".
func CreateExportTask(ctx *gin.Context) {
	checkErr := checkCreateExportTask(ctx)
	if pkg.HasError(checkErr) {
		ctx.JSON(400, checkErr)

		return
	}
	save := ctx.PostForm(paramQuerySave)
	save = strings.ReplaceAll(save, "\\", "/")
	file := ctx.PostForm(paramQueryFile)
	layoutStr := ctx.PostForm(paramQueryLayout)
	layoutTpl, buildErr := buildFramePrams(layoutStr)
	if pkg.HasError(buildErr) {
		ctx.JSON(400, buildErr)

		return
	}
	previewLayoutParams := ctx.PostForm(paramQueryPrevireLayout)
	var previewLayoutMap map[string]string
	err := json.Unmarshal([]byte(previewLayoutParams), &previewLayoutMap)
	if err != nil {
		ctx.JSON(400, requestParamError("预览修改参数格式错误"))

		return
	}
	// 开异步执行
	go exportFrame(save, file, &layoutTpl, previewLayoutMap)

	ctx.JSON(200, NoError{
		Code:   0,
		Errmsg: "success",
	})
}

// @Summary 获取导出进度,输出SSE消息
// @Description 获取导出进度
// @Tags Frame
// @Produce json
// @Router /frame/getExportProgress [get]
// @Success 200 {object} string "SSE消息".
func GetExportProgress(ctx *gin.Context) {
	// 设置头
	setSSEHeader(ctx)

	infos := make([]string, 0)

loop:
	for {
		select {
		case info, ok := <-export_Progress_Chan:
			if ok {
				infos = append(infos, info)
			}
		default:
			time.Sleep(10 * time.Microsecond)

			break loop
		}
	}

	w := ctx.Writer
	flusher, _ := w.(http.Flusher)
	for i := range infos {
		fmt.Fprintf(w, "data: %s\n\n", infos[i])
		flusher.Flush()
	}
}

// 导出执行函数.
func exportFrame(save, file string, layoutTpl *layout.FrameLayout, previewLayoutMap map[string]string) {
	prex := time.Now().Format("2006-01-02-15_04_05")
	task := make(chan struct{}, 6)
	var wg sync.WaitGroup
	errors := make(map[string]pkg.EError, 0)
	for file := range strings.SplitSeq(file, ",") {
		task <- struct{}{}
		wg.Add(1)

		go func(path string) {
			defer wg.Done()

			exifInfo, tpl, checkErr := checkExportFrameTask(path, previewLayoutMap, layoutTpl)
			if pkg.HasError(checkErr) {
				errors[file] = checkErr

				sendExportProgress(path)
				<-task

				return
			}

			exportFrameTask(save, path, prex, exifInfo, tpl)

			time.Sleep(100 * time.Microsecond)

			sendExportProgress(path)
			<-task
		}(file)
	}
	wg.Wait()

	writeExportError(save, errors)
}

// 生成导出失败文件.
func writeExportError(save string, errors map[string]pkg.EError) {
	if len(errors) == 0 {
		return
	}
	if !internal.PathExists(save) {
		return
	}
	savePath := save + "/导出失败.csv"
	csv := pkg.CreateCSV(filepath.Base(savePath), save+"/", true)
	csvData := make([][]string, 0, len(errors))
	for key, e := range errors {
		csvData = append(csvData, []string{key, e.Error.Error()})
	}
	csv.AddData(csvData)
	err := csv.Generate()
	if pkg.HasError(err) {
		internal.Log.Error("生成" + savePath + "失败:" + err.String())
	}
}

// 检查执行导出时的参数.
func checkExportFrameTask(
	path string,
	previewLayoutMap map[string]string,
	layoutTpl *layout.FrameLayout,
) (exiftool.FileMetadata, *layout.FrameLayout, pkg.EError) {
	if !internal.PathExists(path) {
		return exiftool.FileMetadata{}, &layout.FrameLayout{}, pkg.NewErrors(pkg.FILE_NOT_EXIST_ERROR, path+":文件不存在")
	}
	exifInfo, exifErr := getExifAndCheckPhotoLogoExist(path)
	if pkg.HasError(exifErr) {
		return exiftool.FileMetadata{}, &layout.FrameLayout{}, exifErr
	}

	var tpl *layout.FrameLayout
	// 如果预览时修改了参数,按照预览参数为准
	if previewStr, ok := previewLayoutMap[path]; ok {
		previewLayout, buildErr := buildFramePrams(previewStr)
		if pkg.HasError(buildErr) {
			return exiftool.FileMetadata{}, &layout.FrameLayout{}, buildErr
		}
		tpl = &previewLayout
	} else {
		tpl = layoutTpl
	}

	return exifInfo, tpl, pkg.NoError
}

// 执行导出.
func exportFrameTask(save, path, prex string, exifInfo exiftool.FileMetadata, tpl *layout.FrameLayout) {
	plug := frame.GetPlugin()
	plug.CreateFrameImageRGBA(
		map[string]any{
			"sourceImageFile": path,
			"photoType":       "photo",
			"exif":            exifInfo,
			"Params":          tpl,
			"saveImageFile":   save + "/" + prex + "_" + filepath.Base(path),
		},
	)
}

// 发送导出进度.
func sendExportProgress(str string) {
	export_Progress_Chan <- str
}

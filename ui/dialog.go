package ui

import (
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"WaterMark/internal"
	"WaterMark/pkg"
)

// SelectImageFile
//
// 选择单个图片文件,不支持选择raw格式图片
// 如果已选择,则返回对应的图片地址
// 如果没有选择,返回空字符串
//
//	@return string
func (a *App) SelectImageFile() string {
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  "",
		Title:            "请选择图片",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Images (*.jpg;*.JPG;*.jpeg;*.JPEG;)",
				Pattern:     "*.jpg;*.JPG;*.jpeg;*.JPEG;",
			},
		},
	})
	if err != nil {
		internal.Log.Error("SelectImageFile error:" + err.Error())

		return ""
	}

	return strings.ReplaceAll(result, "\\", "/")
}

// SelectImageFileSupportRaw
//
// 选择单个图片文件,支持选择raw格式图片
// 如果已选择,则返回对应的图片地址
// 如果没有选择,返回空字符串
//
//	@return string
func (a *App) SelectImageFileSupportRaw() string {
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  "",
		Title:            "请选择图片",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Images (*.jpg;*.JPG;*.jpeg;*.JPEG;*NEF)",
				Pattern:     "*.jpg;*.JPG;*.jpeg;*.JPEG;*NEF",
			},
		},
	})
	if err != nil {
		internal.Log.Error("SelectImageFileSupportRaw error:" + err.Error())

		return ""
	}

	return strings.ReplaceAll(result, "\\", "/")
}

// SelectMultipleImageFile
//
// 选择多个图片文件
// 如果已选择,则返回对应的图片地址,多个图片地址中间使用,隔开
// 如果没有选择,返回空字符串
//
//	@return string
func (a *App) SelectMultipleImageFile() string {
	result, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  "",
		Title:            "请选择图片",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Images (*.jpg;*.JPG;*.jpeg;*.JPEG;)",
				Pattern:     "*.jpg;*.JPG;*.jpeg;*.JPEG;",
			},
		},
	})
	if err != nil {
		internal.Log.Error("SelectMultipleImageFile error:" + err.Error())

		return ""
	}

	return strings.ReplaceAll(strings.Join(result, ","), "\\", "/")
}

// ShowSaveFileDialog
//
// 保存文件对话框,name 为默认保存的文件名称
//
// @return string.
func (a *App) ShowSaveFileDialog(name string) string {
	result, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  name,
		Title:            "请选择文件保存的路径",
	})
	if err != nil {
		internal.Log.Error("ShowSaveFileDialog error:" + err.Error())

		return ""
	}

	return result
}

// SelectDirectory
//
// 导选择路径
//
// @return string.
func (a *App) SelectDirectory(title string) string {
	result, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            title,
		DefaultDirectory: filepath.Dir(a.lastDir),
	})
	if err != nil {
		internal.Log.Error("SelectDirectory error:" + err.Error())

		return ""
	}
	dir := strings.ReplaceAll(result, "\\", "/")
	a.lastDir = dir

	return dir
}

// ShowExportPhotoTips
//
// 展示导出确认弹窗
//
// @return string.
func (a *App) ShowExportPhotoTips(savePath string) string {
	result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   "导出确认",
		Message: "是否执行导出操作,并将导出文件存放在:\"" + savePath + "\"文件夹中?",
	})
	if err != nil {
		internal.Log.Error("SelectDirectory error:" + err.Error())

		return ""
	}

	return result
}

// GetDirectoryJpgFiles
//
// 获取指定路径下的jpg图片文件名称列表
//
// @return string.
func (a *App) GetDirectoryJpgFiles(path string) string {
	if !internal.PathExists(path) {
		return ""
	}
	files, err := pkg.GetDirFiles(path)
	if pkg.HasError(err) {
		return ""
	}
	arr := make([]string, 0)
	for i := range files {
		ext := filepath.Ext(files[i])
		if !strings.Contains(ext, ".jpg") && !strings.Contains(ext, ".JPG") && !strings.Contains(ext, ".jpeg") &&
			!strings.Contains(ext, ".JPEG") {
			continue
		}
		arr = append(arr, path+"/"+files[i])
	}

	return strings.Join(arr, ",")
}

package gui

import (
	"WaterMark/src/log"
	"context"
	goruntime "runtime"
	"strings"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) Startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	// 设置主题
	if goruntime.GOOS == "windows" {
		wailsruntime.WindowSetDarkTheme(a.ctx)
	}
}

// domReady is called after front-end resources have been loaded
func (a App) DomReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
}

// SelectDirectory 选择文件夹,如果没有选择,返回空字符串
//
//	@return string 选择的文件夹路径
func (a *App) SelectDirectory() string {
	result, err := wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  "",
		Title:            "请选择文件夹",
	})
	if err != nil {
		log.ErrorLogger.Println("SelectDirectory error:" + err.Error())
		return ""
	}
	return result
}

// SelectImageFile 选择图片文件,如果没有选择,返回空字符串
//
//	@return string
func (a *App) SelectImageFile() string {
	result, err := wailsruntime.OpenFileDialog(a.ctx, wailsruntime.OpenDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  "",
		Title:            "请选择图片",
		Filters: []wailsruntime.FileFilter{
			{
				DisplayName: "Images (*.jpg;*.JPG;*.jpeg;*.JPEG;)",
				Pattern:     "*.jpg;*.JPG;*.jpeg;*.JPEG;",
			},
		},
	})
	if err != nil {
		log.ErrorLogger.Println("SelectImageFile error:" + err.Error())
		return ""
	}
	return result
}

// SelectImageFile 选择图片文件,如果没有选择,返回空字符串
//
//	@return string
func (a *App) SelectMultipleImageFile() string {
	result, err := wailsruntime.OpenMultipleFilesDialog(a.ctx, wailsruntime.OpenDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  "",
		Title:            "请选择图片",
		Filters: []wailsruntime.FileFilter{
			{
				DisplayName: "Images (*.jpg;*.JPG;*.jpeg;*.JPEG;*NEF)",
				Pattern:     "*.jpg;*.JPG;*.jpeg;*.JPEG;*NEF",
			},
		},
	})
	if err != nil {
		log.ErrorLogger.Println("SelectImageFile error:" + err.Error())
		return ""
	}
	return strings.Join(result, ",")
}

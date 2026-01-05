package ui

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"WaterMark/engine"
	"WaterMark/internal"
	"WaterMark/message"
)

// App struct.
//
//nolint:containedctx
type App struct {
	ctx           context.Context
	data          map[string]string
	previewParams map[string]string
	lastDir       string
	mtx           sync.Mutex
}

// 预览桶.
var bucketPreview = "preview"

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{
		data:          make(map[string]string),
		previewParams: make(map[string]string),
	}
}

// startup is called at application startup.
func (a *App) Startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	// 初始化引擎
	engine.InitAllTools()
	// 绑定错误消息
	go a.BindErrorMessage()
}

// domReady is called after front-end resources have been loaded.
func (a *App) DomReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) BeforeClose(ctx context.Context) bool {
	// 关闭启动的工具
	engine.QuitAllTools()
	message.Close()
	internal.CleanDir()

	return false
}

// shutdown is called at application termination.
func (a *App) Shutdown(ctx context.Context) {
	internal.Log.Info("UI程序退出")
}

// 获取后端接口的地址.
func (a *App) GetApiServerHost() string {
	return "http://" + viper.GetString("server.address")
}

// 绑定消息处理逻辑.
func (a *App) BindErrorMessage() {
	for {
		select {
		case err, ok := <-message.Error_Messge_Chan:
			if ok {
				_, e := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
					Type:    runtime.ErrorDialog,
					Title:   "发生了一个错误",
					Message: err,
				})
				if e != nil {
					internal.Log.Panic(err)
				}
				internal.Log.Error(e.Error())
			}
		default:
			time.Sleep(10 * time.Microsecond)
		}
	}
}

// 清理暂存空间.
func (a *App) TemporaryClean() string {
	a.mtx.Lock()
	a.data = make(map[string]string)
	a.previewParams = make(map[string]string)
	a.mtx.Unlock()

	return "clean success"
}

// 暂存key对应的数据.
func (a *App) TemporaryStorage(key, value, bucket string) string {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	// 根据不同的桶名称选择不同的存放数据的map
	switch bucket {
	case bucketPreview:
		a.previewParams[key] = value
	default:
		a.data[key] = value
	}

	return value
}

// 获取暂存的数据.
func (a *App) GetTemporaryStorage(key, bucket string) string {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if bucket == bucketPreview {
		if v, ok := a.previewParams[key]; ok {
			return v
		}

		return ""
	}

	if v, ok := a.data[key]; ok {
		return v
	}

	return ""
}

// 获取暂存的全部数据.
func (a *App) GetTemporaryAll(bucket string) string {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if bucket == bucketPreview {
		jsonByte, err := json.Marshal(a.previewParams)
		if err == nil {
			return string(jsonByte)
		}
		internal.Log.Error(err.Error())

		return ""
	}
	jsonByte, err := json.Marshal(a.data)
	if err == nil {
		return string(jsonByte)
	}
	internal.Log.Error(err.Error())

	return ""
}

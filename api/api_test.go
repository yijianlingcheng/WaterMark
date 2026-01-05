package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"WaterMark/engine"
	"WaterMark/internal"
	"WaterMark/message"
	"WaterMark/pkg"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	cleanUp()
	os.Exit(code)
}

func setup() {
	internal.SetAppMode(internal.APP_API_DEV)
	internal.InitAppConfigsAndRes()
	go ServerStart()
	engine.InitAllTools()
}

func cleanUp() {
	engine.QuitAllTools()
	message.Close()
}

// TestGetPhotosExifInfo 测试获取照片EXIF信息接口.
func TestGetPhotosExifInfo(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name     string
		method   string
		reqType  string
		reqBody  map[string]string
		wantCode int
		wantBody string
	}{
		{
			name:     "获取照片EXIF信息失败-请求方法错误1",
			method:   http.MethodGet,
			reqBody:  map[string]string{"file": "test.jpg"},
			wantCode: http.StatusNotFound,
			wantBody: `404 page not found`,
		}, {
			name:     "获取照片EXIF信息失败-请求方法错误2",
			method:   http.MethodPost,
			reqType:  "application/json",
			reqBody:  map[string]string{"file": "test.jpg"},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"file参数为空","file":""}`,
		}, {
			name:     "获取照片EXIF信息失败-文件不存在",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": "test.jpg"},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":1000001,"errmsg":"file请求的文件不存在","file":"test.jpg"}`,
		}, {
			name:     "获取照片EXIF信息失败-参数错误",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"xxxx": "xxxx"},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"file参数为空","file":""}`,
		}, {
			name:     "获取照片EXIF信息-文件不存在-路径错误",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": "./tests/DSC_4352.JPG"},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":1000001,"errmsg":"file请求的文件不存在","file":"./tests/DSC_4352.JPG"}`,
		}, {
			name:     "获取照片EXIF信息-成功1",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/DSC_4352.JPG"},
			wantCode: http.StatusOK,
			wantBody: `{"file":"` + internal.GetRootPath() + "/tests/DSC_4352.JPG" + `","器材":"NIKON Z 8 NIKKOR Z 24-120mm f/4 S","模式":"曝光模式:Shutter speed priority AE 测光模式:Multi-segment 曝光补偿:0","参数":"光圈:4 快门:1/3200 ISO:100","焦距":"24.0 mm (35 mm equivalent: 24.0 mm) 视角:69.9 deg (0.49 m)","色彩":"白平衡:Auto 色彩空间:sRGB","时间":"2026:01:04 16:10:18","快门次数":"4169","Orientation":"Horizontal (normal)","ImageWidth":"8256","ImageHeight":"5504","Make":"NIKON CORPORATION","Model":"NIKON Z 8","LensModel":"NIKKOR Z 24-120mm f/4 S","FocalLength":"24.0 mm","FNumber":"4","ExposureTime":"1/3200","ISO":"100","FileName":"DSC_4352.JPG","ImageSize":"8256x5504","ImageDataSize":"9862956","errmsg":"","code":0}`,
		}, {
			name:     "获取照片EXIF信息-成功2",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/DSC_4352_1.JPG"},
			wantCode: http.StatusOK,
			wantBody: `{"file":"` + internal.GetRootPath() + "/tests/DSC_4352_1.JPG" + `","器材":" ","模式":"曝光模式: 测光模式: 曝光补偿:0","参数":"光圈: 快门: ISO:","焦距":" 视角:","色彩":"白平衡: 色彩空间:","时间":"","快门次数":"","Orientation":"","ImageWidth":"8256","ImageHeight":"5504","Make":"","Model":"","LensModel":"","FocalLength":"","FNumber":"","ExposureTime":"","ISO":"","FileName":"DSC_4352_1.JPG","ImageSize":"8256x5504","ImageDataSize":"","errmsg":"","code":0}`,
		},
	}
	requestURL := fmt.Sprintf("http://%s%s", viper.GetString("server.address"), "/view/getImagesExifInfo?")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个 http 请求
			if tc.method == http.MethodGet {
				params := []string{}
				for k, v := range tc.reqBody {
					params = append(params, fmt.Sprintf("%s=%s", k, v))
				}
				resp, err := http.Get(requestURL + strings.Join(params, "&"))
				// 断言没有错误
				assert.NoError(t, err)
				// 设置请求头
				defer resp.Body.Close()
				// 检查响应码
				assert.Equal(t, tc.wantCode, resp.StatusCode)
				respBody, err := io.ReadAll(resp.Body)
				// 断言没有错误
				assert.NoError(t, err)
				// 检查内容
				assert.Equal(t, tc.wantBody, string(respBody))
			}
			if tc.method == http.MethodPost {
				values := make(url.Values)
				for k, v := range tc.reqBody {
					values.Set(k, v)
				}
				resp, err := http.Post(requestURL, tc.reqType, strings.NewReader(values.Encode()))
				// 断言没有错误
				assert.NoError(t, err)
				// 设置请求头
				defer resp.Body.Close()
				// 检查响应码
				assert.Equal(t, tc.wantCode, resp.StatusCode)
				respBody, err := io.ReadAll(resp.Body)
				// 断言没有错误
				assert.NoError(t, err)
				// 检查内容
				assert.Equal(t, tc.wantBody, string(respBody))
			}
		})
	}
}

// 测试展示照片内容
func TestShowImage(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name     string
		method   string
		reqType  string
		reqBody  map[string]string
		wantCode int
		wantBody string
	}{
		{
			name:     "展示照片内容-请求方法错误",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/DSC_4352.JPG"},
			wantCode: http.StatusNotFound,
			wantBody: `404 page not found`,
		},
		{
			name:     "展示照片内容-参数错误1",
			method:   http.MethodGet,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"xxx": ""},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":1000001,"errmsg":"file请求的文件不存在","file":""}`,
		},
		{
			name:     "展示照片内容-参数错误2",
			method:   http.MethodGet,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": ""},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":1000001,"errmsg":"file请求的文件不存在","file":""}`,
		},
	}
	requestURL := fmt.Sprintf("http://%s%s", viper.GetString("server.address"), "/view/showImage?")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个 http 请求
			if tc.method == http.MethodGet {
				params := []string{}
				for k, v := range tc.reqBody {
					params = append(params, fmt.Sprintf("%s=%s", k, v))
				}
				resp, err := http.Get(requestURL + strings.Join(params, "&"))
				// 断言没有错误
				assert.NoError(t, err)
				// 设置请求头
				defer resp.Body.Close()
				// 检查响应码
				assert.Equal(t, tc.wantCode, resp.StatusCode)
				respBody, err := io.ReadAll(resp.Body)
				// 断言没有错误
				assert.NoError(t, err)
				// 检查内容
				assert.Equal(t, tc.wantBody, string(respBody))
			}
			if tc.method == http.MethodPost {
				values := make(url.Values)
				for k, v := range tc.reqBody {
					values.Set(k, v)
				}
				resp, err := http.Post(requestURL, tc.reqType, strings.NewReader(values.Encode()))
				// 断言没有错误
				assert.NoError(t, err)
				// 设置请求头
				defer resp.Body.Close()
				// 检查响应码
				assert.Equal(t, tc.wantCode, resp.StatusCode)
				respBody, err := io.ReadAll(resp.Body)
				// 断言没有错误
				assert.NoError(t, err)
				// 检查内容
				assert.Equal(t, tc.wantBody, string(respBody))
			}
		})
	}
}

// 测试导出照片exif信息
func TestExifInfoExportBySaveFile(t *testing.T) {
	savePath := internal.GetRuntimePath("export.csv")
	// 定义测试用例
	testCases := []struct {
		name     string
		method   string
		reqType  string
		reqBody  map[string]string
		wantCode int
		wantBody string
	}{
		{
			name:     "导出照片exif信息-请求方法错误1",
			method:   http.MethodGet,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"xxx": ""},
			wantCode: http.StatusNotFound,
			wantBody: `404 page not found`,
		},
		{
			name:     "导出照片exif信息-请求方法错误2",
			method:   http.MethodPost,
			reqType:  "application/json",
			reqBody:  map[string]string{"xxx": ""},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":1000001,"errmsg":"file请求的文件不存在","file":""}`,
		},
		{
			name:     "导出照片exif信息-请求方法错误3",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"xxx": ""},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":1000001,"errmsg":"file请求的文件不存在","file":""}`,
		},
		{
			name:     "导出照片exif信息-请求方法错误4",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": ""},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":1000001,"errmsg":"file请求的文件不存在","file":""}`,
		},
		{
			name:     "导出照片exif信息-请求方法错误5",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/DSC_4352.JPG"},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"未选择保存的路径","file":""}`,
		},
		{
			name:     "导出照片exif信息-请求方法错误6",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/DSC_4352.JPG", "save": ""},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"未选择保存的路径","file":""}`,
		},
		{
			name:    "导出照片exif信息-请求方法错误7",
			method:  http.MethodPost,
			reqType: "application/x-www-form-urlencoded",
			reqBody: map[string]string{
				"file": internal.GetRootPath() + "/tests/DSC_4352",
				"save": savePath,
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":1000001,"errmsg":"file请求的文件不存在","file":"E:/WaterMark/tests/DSC_4352"}`,
		},
		{
			name:    "导出照片exif信息-导出成功",
			method:  http.MethodPost,
			reqType: "application/x-www-form-urlencoded",
			reqBody: map[string]string{
				"file": internal.GetRootPath() + "/tests/DSC_4352.JPG",
				"save": savePath,
			},
			wantCode: http.StatusOK,
			wantBody: `{"code":0,"errmsg":"no error"}`,
		},
	}
	requestURL := fmt.Sprintf("http://%s%s", viper.GetString("server.address"), "/view/exifInfoExportv2")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个 http 请求
			if tc.method == http.MethodGet {
				params := []string{}
				for k, v := range tc.reqBody {
					params = append(params, fmt.Sprintf("%s=%s", k, v))
				}
				resp, err := http.Get(requestURL + strings.Join(params, "&"))
				// 断言没有错误
				assert.NoError(t, err)
				// 设置请求头
				defer resp.Body.Close()
				// 检查响应码
				assert.Equal(t, tc.wantCode, resp.StatusCode)
				respBody, err := io.ReadAll(resp.Body)
				// 断言没有错误
				assert.NoError(t, err)
				// 检查内容
				assert.Equal(t, tc.wantBody, string(respBody))
			}
			if tc.method == http.MethodPost {
				values := make(url.Values)
				for k, v := range tc.reqBody {
					values.Set(k, v)
				}
				resp, err := http.Post(requestURL, tc.reqType, strings.NewReader(values.Encode()))
				// 断言没有错误
				assert.NoError(t, err)
				// 设置请求头
				defer resp.Body.Close()
				// 检查响应码
				assert.Equal(t, tc.wantCode, resp.StatusCode)
				respBody, err := io.ReadAll(resp.Body)
				// 断言没有错误
				assert.NoError(t, err)
				// 检查内容
				assert.Equal(t, tc.wantBody, string(respBody))
			}
		})
	}
	if internal.PathExists(savePath) {
		os.Remove(savePath)
	}
}

// 测试显示照片边框-验证参数
func TestShowPhotoFrameParams(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name     string
		method   string
		reqType  string
		reqBody  map[string]string
		wantCode int
		wantBody string
	}{
		{
			name:     "照片边框-请求方法错误1",
			method:   http.MethodGet,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"xxx": ""},
			wantCode: http.StatusNotFound,
			wantBody: `404 page not found`,
		},
		{
			name:     "照片边框-请求方法错误2",
			method:   http.MethodGet,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": ""},
			wantCode: http.StatusNotFound,
			wantBody: `404 page not found`,
		},
		{
			name:     "照片边框-请求方法错误3",
			method:   http.MethodPost,
			reqType:  "application/json",
			reqBody:  map[string]string{"file": ""},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"file参数为空"}`,
		},
		{
			name:     "照片边框-请求方法错误4",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": ""},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"file参数为空"}`,
		},
		{
			name:     "照片边框-请求方法错误5",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/1.jpg"},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"type参数为空"}`,
		},
		{
			name:     "照片边框-请求方法错误6",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/1.jpg", "type": ""},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"type参数为空"}`,
		},
		{
			name:     "照片边框-请求方法错误7",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/1.jpg", "type": "xxx"},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"type参数类型错误"}`,
		},
		{
			name:     "照片边框-请求方法错误8",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/1.jpg", "type": "null"},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"type参数类型错误"}`,
		},
		{
			name:     "照片边框-请求方法错误9",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/1.jpg", "type": "border"},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"layout参数为空"}`,
		},
		{
			name:     "照片边框-请求方法错误10",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/1.jpg", "type": "photo"},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"layout参数为空"}`,
		},
		{
			name:     "照片边框-请求方法错误11",
			method:   http.MethodPost,
			reqType:  "application/x-www-form-urlencoded",
			reqBody:  map[string]string{"file": internal.GetRootPath() + "/tests/1.jpg", "type": "photo", "layout": ""},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"layout参数为空"}`,
		},
		{
			name:    "照片边框-请求方法错误12",
			method:  http.MethodPost,
			reqType: "application/x-www-form-urlencoded",
			reqBody: map[string]string{
				"file":   internal.GetRootPath() + "/tests/1.jpg",
				"type":   "photo",
				"layout": "xxx",
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":9000002,"errmsg":"xxx:布局信息格式错误,json解析失败"}`,
		},
		{
			name:    "照片边框-请求方法错误13",
			method:  http.MethodPost,
			reqType: "application/x-www-form-urlencoded",
			reqBody: map[string]string{
				"file":   internal.GetRootPath() + "/tests/1.jpg",
				"type":   "photo",
				"layout": "{}",
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":6000001,"errmsg":"布局类型查找失败"}`,
		},
		{
			name:    "照片边框-请求方法错误14",
			method:  http.MethodPost,
			reqType: "application/x-www-form-urlencoded",
			reqBody: map[string]string{
				"file":   internal.GetRootPath() + "/tests/1.jpg",
				"type":   "photo",
				"layout": "{\"xxxx\":\"yyyy\"}",
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":6000001,"errmsg":"布局类型查找失败"}`,
		},
		{
			name:    "照片边框-请求方法错误15",
			method:  http.MethodPost,
			reqType: "application/x-www-form-urlencoded",
			reqBody: map[string]string{
				"file":   internal.GetRootPath() + "/tests/1.jpg",
				"type":   "photo",
				"layout": "{\"frame_name\":\"yyyy\"}",
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":6000001,"errmsg":"布局类型查找失败"}`,
		},
		{
			name:    "照片边框-请求方法错误16",
			method:  http.MethodPost,
			reqType: "application/x-www-form-urlencoded",
			reqBody: map[string]string{
				"file":   internal.GetRootPath() + "/tests/1.jpg",
				"type":   "photo",
				"layout": "{\"frame_name\":\"\"}",
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":6000001,"errmsg":"布局类型查找失败"}`,
		},
		{
			name:    "照片边框-请求方法错误17",
			method:  http.MethodPost,
			reqType: "application/x-www-form-urlencoded",
			reqBody: map[string]string{
				"file":   internal.GetRootPath() + "/tests/1.jpg",
				"type":   "photo",
				"layout": "{\"frame_name\":\"null\"}",
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"code":6000001,"errmsg":"布局类型查找失败"}`,
		},
	}
	requestURL := fmt.Sprintf("http://%s%s", viper.GetString("server.address"), "/frame/showPhotoFrame?")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个 http 请求
			if tc.method == http.MethodGet {
				params := []string{}
				for k, v := range tc.reqBody {
					params = append(params, fmt.Sprintf("%s=%s", k, v))
				}
				resp, err := http.Get(requestURL + strings.Join(params, "&"))
				// 断言没有错误
				assert.NoError(t, err)
				// 设置请求头
				defer resp.Body.Close()
				// 检查响应码
				assert.Equal(t, tc.wantCode, resp.StatusCode)
				respBody, err := io.ReadAll(resp.Body)
				// 断言没有错误
				assert.NoError(t, err)
				// 检查内容
				assert.Equal(t, tc.wantBody, string(respBody))
			}
			if tc.method == http.MethodPost {
				values := make(url.Values)
				for k, v := range tc.reqBody {
					values.Set(k, v)
				}
				resp, err := http.Post(requestURL, tc.reqType, strings.NewReader(values.Encode()))
				// 断言没有错误
				assert.NoError(t, err)
				// 设置请求头
				defer resp.Body.Close()
				// 检查响应码
				assert.Equal(t, tc.wantCode, resp.StatusCode)
				respBody, err := io.ReadAll(resp.Body)
				// 断言没有错误
				assert.NoError(t, err)
				// 检查内容
				assert.Equal(t, tc.wantBody, string(respBody))
			}
		})
	}
}

// 测试照片边框功能-验证功能
func TestShowPhotoFrameFeature(t *testing.T) {
	requestURL := fmt.Sprintf("http://%s%s", viper.GetString("server.address"), "/frame/getFrameTemplateInfo?")
	req, err := http.Get(requestURL)
	if req == nil {
		t.Errorf("TestShowPhotoFrameFeaturehttp.Get requestURL failed, err: %v", err)
	}
	if req.StatusCode != http.StatusOK {
		t.Errorf("TestShowPhotoFrameFeaturehttp.Get requestURL failed, status code: %d", req.StatusCode)
	}
	defer req.Body.Close()
	tplBody, err := io.ReadAll(req.Body)
	if err != nil {
		t.Errorf("TestShowPhotoFrameFeaturehttp.Get requestURL failed, err: %v", err)
	}
	tplMaps := make(map[string]any)
	err = json.Unmarshal(tplBody, &tplMaps)
	if err != nil {
		t.Errorf("TestShowPhotoFrameFeaturehttp.Get requestURL failed, err: %v", err)
	}
	if int(tplMaps["code"].(float64)) != 0 {
		t.Errorf("TestShowPhotoFrameFeaturehttp.Get requestURL failed, code: %d", int(tplMaps["code"].(float64)))
	}
	if len(tplMaps["list"].(map[string]any)) == 0 {
		t.Errorf(
			"TestShowPhotoFrameFeaturehttp.Get requestURL failed, list len: %d",
			len(tplMaps["list"].(map[string]any)),
		)
	}

	requestURL = fmt.Sprintf("http://%s%s", viper.GetString("server.address"), "/frame/showPhotoFrame?")
	listAny := tplMaps["list"].(map[string]any)

	for key := range listAny {
		for i, name := range []string{"/tests/1.jpg", "/tests/2.jpg"} {
			t.Run(key+fmt.Sprintf("%d", i+1), func(t *testing.T) {
				savePath := internal.GetRootPath() + "/tests/" + key + "_" + fmt.Sprintf("%d", i+1) + ".jpg"
				originFramePath := internal.GetRootPath() + "/tests/" + key + "/" + filepath.Base(name)
				originPath := internal.GetRootPath() + name
				// 创建一个 http 请求
				values := make(url.Values)
				values.Set("file", originPath)
				values.Set("type", "photo")
				values.Set("layout", fmt.Sprintf("{\"frame_name\":\"%s\"}", key))
				values.Set("save", savePath)
				resp, err := http.Post(
					requestURL,
					"application/x-www-form-urlencoded",
					strings.NewReader(values.Encode()),
				)
				// 断言没有错误
				assert.NoError(t, err)
				// 设置请求头
				defer resp.Body.Close()
				// 检查响应码
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				_, bodyErr := io.ReadAll(resp.Body)
				// 断言没有错误
				assert.NoError(t, bodyErr)
				// 检查内容
				if !assert.FileExists(t, savePath) {
					t.Errorf("TestShowPhotoFrameFeaturehttp.Post savePath failed, savePath: %s", savePath)
				}
				newFileMd5, md5Error := pkg.GetFileMD5(savePath)
				if pkg.HasError(md5Error) || newFileMd5 == "" {
					t.Errorf("TestShowPhotoFrameFeaturehttp.Post savePath failed, md5Error: %v", md5Error)
				}

				originFileMd5, originMd5Error := pkg.GetFileMD5(originFramePath)
				if pkg.HasError(originMd5Error) || originFileMd5 == "" {
					t.Errorf("TestShowPhotoFrameFeaturehttp.Post savePath failed, md5Error: %v", originFileMd5)
				}
				if newFileMd5 != originFileMd5 {
					t.Errorf(
						"TestShowPhotoFrameFeaturehttp.Post savePath failed, newFileMd5: %s, originFileMd5: %s",
						newFileMd5,
						originFileMd5,
					)
				}
				os.Remove(savePath)
			})
		}
	}
}

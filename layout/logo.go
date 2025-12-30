package layout

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"
	"sync"

	"WaterMark/internal"
	"WaterMark/pkg"
)

type (
	// 相机logo列表.
	logosMaps struct {
		logoMap map[string]*Logo
	}

	// 相机logo元素.
	Logo struct {
		LogoImage image.Image
		Name      string
		Ext       string
		LogoPath  string
		Width     int
		Height    int
		IsLoad    bool
	}
)

var (
	// logos 对象.
	logos *logosMaps

	// 修改logosMaps需要使用到的锁,防止并发读写导致map panic.
	writeLogosMapMtx sync.Mutex

	// 未知logo,使用特殊标识代替.
	UNSUPPORT_LOGO = "unsupported_logo"
)

// 根据exif信息中的model字段获取logo名称.
func GetLogoNameByMake(name string) string {
	// 转换成全小写比较
	name = strings.ToLower(name)
	if _, ok := logos.logoMap[name]; ok {
		return name
	}
	for key := range logos.logoMap {
		if strings.Contains(name, key) {
			return key
		}
	}

	return UNSUPPORT_LOGO
}

// 检查logo名称是否不支持.
func CheckLogoIsUnSupport(logoName string) bool {
	return logoName == UNSUPPORT_LOGO
}

// 根据logo名称获取logo信息.
func GetLogoImageByName(name string) (*Logo, pkg.EError) {
	if logoItem, ok := logos.logoMap[name]; ok {
		return logoItem, pkg.NoError
	}

	return &Logo{}, pkg.ImageLogoNotFindError
}

// 根据logo的名称与高度获取等比例压缩之后的宽高.
func GetLogoXAndYByNameAndHeight(name string, height int) map[string]int {
	r := make(map[string]int, 2)

	r["height"] = height
	r["width"] = 0

	logo, findErr := GetLogoImageByName(name)
	if pkg.HasError(findErr) {
		return r
	}
	r["width"] = logo.Width * height / logo.Height

	return r
}

// 根据logo名称,宽高获取logo信息.
func GetLogoImageByNameAndWidhtAndHeight(name string, width, height int) (*Logo, pkg.EError) {
	// 构造新的logo文件名称
	newLogoName := fmt.Sprintf("%d_%d_%s", width, height, name)

	logoItem, err := GetLogoImageByName(newLogoName)
	if !pkg.HasError(err) {
		return logoItem, err
	}

	originLogo, findErr := GetLogoImageByName(name)
	if pkg.HasError(findErr) {
		return &Logo{}, pkg.ImageLogoNotFindError
	}

	// 如果上面没找到,说明需要根据尺寸重新生成一个logo并加载进来
	// 加锁防止并发写
	writeLogosMapMtx.Lock()
	defer writeLogosMapMtx.Unlock()

	// 生成指定的新图片
	newLogoImage := pkg.GenerateImageByWidthHeight(originLogo.LogoImage, width, height)
	// 写入logos中
	logos.logoMap[newLogoName] = newLogoWithImage(newLogoName, newLogoImage)

	return logos.logoMap[newLogoName], pkg.NoError
}

// 加载指定的logo图片,并返回一个logo结构体.
func newLogo(name, fullPath string) (*Logo, pkg.EError) {
	name = strings.ToLower(name)
	ext := filepath.Ext(fullPath)
	imgaeDecode, loadErr := pkg.LoadImageWithDecode(fullPath)
	if pkg.HasError(loadErr) {
		return &Logo{}, loadErr
	}

	return &Logo{
		IsLoad:    true,
		Width:     imgaeDecode.Bounds().Dx(),
		Height:    imgaeDecode.Bounds().Dy(),
		Name:      name,
		Ext:       ext,
		LogoPath:  fullPath,
		LogoImage: imgaeDecode,
	}, loadErr
}

// 返回一个logo结构体.
func newLogoWithImage(name string, logoImage image.Image) *Logo {
	return &Logo{
		IsLoad:    true,
		Width:     logoImage.Bounds().Dx(),
		Height:    logoImage.Bounds().Dy(),
		Name:      name,
		Ext:       ".png",
		LogoPath:  "",
		LogoImage: logoImage,
	}
}

// 机logo初始化.
func LogosImagesInit() pkg.EError {
	logos = &logosMaps{
		logoMap: make(map[string]*Logo),
	}
	// 自动注册在结构体中
	dir := internal.GetLogosPath("")
	logoFiles, err := pkg.GetDirFiles(dir)
	if pkg.HasError(err) {
		return err
	}

	// 加锁写入数据
	writeLogosMapMtx.Lock()
	defer writeLogosMapMtx.Unlock()

	var logoErr pkg.EError
	for _, f := range logoFiles {
		// 获取全路径,文件名称,文件扩展名
		fullPath := internal.GetLogosPath(f)
		baseName := filepath.Base(fullPath)
		extName := filepath.Ext(fullPath)
		logoName := strings.ReplaceAll(baseName, extName, "")

		// logo写入map
		logos.logoMap[logoName], logoErr = newLogo(logoName, fullPath)
		if pkg.HasError(logoErr) {
			internal.Log.Error()

			return logoErr
		}
	}

	return pkg.NoError
}

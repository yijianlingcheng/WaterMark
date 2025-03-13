package images

import (
	"WaterMark/src/cmd"
	"errors"
	"strings"

	"github.com/spf13/viper"
)

// LogoList app.yaml文件中的logos配置项
type LogoList struct {
	Logos []Logo
}

// Logo app.yaml文件中的logo配置项
type Logo struct {
	Id        string
	Path      string
	TransPath string
}

// logoList
var logoList LogoList

// LoadLogoList 加载logo列表,将app.yaml配置中的logos序列化到结构体中,方便程序进行获取
func LoadLogoList() {
	viper.Unmarshal(&logoList)
}

// findLogo 查找logo,如果是不知道的logo类型,返回空字符串;查找成功之后返回对应的图片路径
//
//	@param str 相机生产商
//	@return string
//	@return string
func findLogo(str string) (string, string) {
	for _, v := range logoList.Logos {
		if strings.Contains(str, v.Id) {
			return v.Path, v.TransPath
		}
	}
	return "", ""
}

// getLogoPath 获取logo路径,此处返回的绝对地址
//
//	@param str
//	@return string
//	@return string
func getLogoPath(str string) (string, string, error) {
	r1, r2 := findLogo(str)
	if r1 == "" {
		return r1, r2, errors.New(str + "不支持的Logo,请前往app.yaml文件添加")
	}
	return cmd.GetPwdPath(r1), cmd.GetPwdPath(r2), nil
}

// LogoTemplate logo模板
type LogoTemplate struct {

	// LogoWidth logo宽度
	Width int

	// LogoHeight logo高度
	Height int

	// LogoMarginTop logo距离原始图片高度
	MarginTop int

	// LogoMarginLeft logo距离原始图片左边宽度
	MarginLeft int

	// LogoMarginRight logo距离原始图片右边宽度
	MarginRight int
}

// // setDefaultLogoTemplate 默认的Logo模板
// //
// //	@param w
// func setDefaultLogoTemplate(w *LogoTemplate) {
// 	w.Width = 300
// 	w.Height = 300
// 	w.MarginTop = 0
// 	w.MarginLeft = 0
// 	w.MarginRight = 0
// }

// // newLogoTemplate 构造一个Logo模板
// //
// //	@param opts 可选参数,map结构
// //	@return *LogoTemplate
// func newLogoTemplate(opts ...map[string]string) *LogoTemplate {
// 	l := LogoTemplate{}
// 	setDefaultLogoTemplate(&l)
// 	if len(opts) > 0 {
// 		val := reflect.ValueOf(&l)
// 		val = val.Elem()
// 		for i, v := range opts[0] {
// 			fieldVa := val.FieldByName(i)
// 			fieldVa.Set(reflect.ValueOf(v))
// 		}
// 	}
// 	return &l
// }

// // setOptions
// //
// //	@param opts
// func (l *LogoTemplate) setOptions(opts map[string]string) {
// 	val := reflect.ValueOf(&l)
// 	val = val.Elem()
// 	for i, v := range opts {
// 		fieldVa := val.FieldByName(i)
// 		fieldVa.Set(reflect.ValueOf(v))
// 	}
// }

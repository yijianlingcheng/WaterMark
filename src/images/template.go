package images

import (
	"errors"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var tpls []WaterMarkTemplate

// LoadTemplate 加载模板配置文件
func LoadTemplate() {
	viper.UnmarshalKey("tpls", &tpls)
}

// GetTemplates 获取模板列表
//
//	@return []WaterMarkTemplate
func GetTemplates() []WaterMarkTemplate {
	return tpls
}

// findTemplate 查找模板
//
//	@param id
//	@return WaterMarkTemplate
//	@return error
func findTemplateById(id string) (WaterMarkTemplate, error) {
	for _, v := range tpls {
		if v.ID == id {
			return v, nil
		}
	}
	return WaterMarkTemplate{}, errors.New(id + ":模板文件不存在")
}

// WaterMarkTemplate 水印模板
type WaterMarkTemplate struct {

	// Stackblur 是否高斯模糊
	Stackblur bool

	// BlurRadius 高斯模糊半径
	BlurRadius int

	// ID
	ID string

	// Type 模板类型(BOTTOM_LOGO_LEFT_AUTO,BOTTOM_LOGO_CENTER_AUTO,BOTTOM_LOGO_RIGHT_AUTO,STACK_BLUR_AUTO)
	Type string

	// Description 模板描述
	Description string

	// BorderTemplate 边框模板
	BorderTemplate *BorderTemplate

	// LogoTemplate logo模板
	LogoTemplate *LogoTemplate

	// WordsTemplate 文字模板
	WordsTemplate *WordsTemplate

	// SeparateTamplate 分隔符模板
	SeparateTamplate *SeparateTamplate
}

// // newWaterMarkTemplate 初始化一个水印模板
// //
// //	@param t
// //	@return *WaterMarkTemplate
// func newWaterMarkTemplate(t string) *WaterMarkTemplate {
// 	w := WaterMarkTemplate{}
// 	w.Type = t
// 	w.BorderTemplate = newBorderTemplate()
// 	w.LogoTemplate = newLogoTemplate()
// 	w.WordsTemplate = newWordsTemplate()
// 	w.SeparateTamplate = newSeparateTamplate()
// 	return &w
// }

// newEmptyWaterMarkTemplate
//
//	@return WaterMarkTemplate
func newEmptyWaterMarkTemplate() *WaterMarkTemplate {
	w := WaterMarkTemplate{}
	return &w
}

// TemplateToYaml 转yaml
//
//	@return string
func (w *WaterMarkTemplate) TemplateToYaml() string {
	bytes, err := yaml.Marshal(w)
	if err != nil {
		return ""
	}
	return string(bytes)
}

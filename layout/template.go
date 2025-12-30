package layout

import (
	"encoding/json"
	"os"

	"WaterMark/internal"
	"WaterMark/pkg"
)

type (
	FrameLayouts struct {
		List []FrameLayout `json:"list"`
	}

	// 布局.
	FrameLayout struct {
		TextOneFontColor      string `json:"text_one_font_color"`
		Type                  string `json:"frame_type"`
		Layout                string `json:"frame_layout"`
		TextFourFontFile      string `json:"text_four_font_file"`
		TextFourFontColor     string `json:"text_four_font_color"`
		TextFourContent       string `json:"text_four_content"`
		TextThreeFontFile     string `json:"text_three_font_file"`
		BgColor               string `json:"bg_color"`
		TextThreeFontColor    string `json:"text_three_font_color"`
		TextThreeContent      string `json:"text_three_content"`
		TextTwoFontFile       string `json:"text_two_font_file"`
		Name                  string `json:"frame_name"`
		TextTwoFontColor      string `json:"text_two_font_color"`
		TextTwoContent        string `json:"text_two_content"`
		TextOneContent        string `json:"text_one_content"`
		TextOneFontFile       string `json:"text_one_font_file"`
		SeparatorColor        string `json:"separator_color"`
		LogoRatio             int    `json:"logo_ratio"`
		TextRatio             int    `json:"text_ratio"`
		LogoMarginRight       int    `json:"logo_margin_right"`
		TextThreeMarginRight  int    `json:"text_three_margin_right"`
		TextOneMarginLeft     int    `json:"text_one_margin_left"`
		TextOneMarginRight    int    `json:"text_one_margin_right"`
		TextOneMarginTop      int    `json:"text_one_margin_top"`
		TextOneMarginBottom   int    `json:"text_one_margin_bottom"`
		LogoMarginBottom      int    `json:"logo_margin_bottom"`
		TextTwoFontSize       int    `json:"text_two_font_size"`
		LogoMarginTop         int    `json:"logo_margin_top"`
		LogoMarginLeft        int    `json:"logo_margin_left"`
		TextTwoMarginLeft     int    `json:"text_two_margin_left"`
		TextTwoMarginRight    int    `json:"text_two_margin_right"`
		TextTwoMarginTop      int    `json:"text_two_margin_top"`
		TextTwoMarginBottom   int    `json:"text_two_margin_bottom"`
		LogoHeight            int    `json:"logo_height"`
		TextThreeFontSize     int    `json:"text_three_font_size"`
		LogoWidth             int    `json:"logo_width"`
		MainMarginBottom      int    `json:"main_margin_bottom"`
		TextThreeMarginLeft   int    `json:"text_three_margin_left"`
		TextOneFontSize       int    `json:"text_one_font_size"`
		TextThreeMarginTop    int    `json:"text_three_margin_top"`
		TextThreeMarginBottom int    `json:"text_three_margin_bottom"`
		MainMarginTop         int    `json:"main_margin_top"`
		TextFourFontSize      int    `json:"text_four_font_size"`
		MainMarginRight       int    `json:"main_margin_right"`
		MainMarginLeft        int    `json:"main_margin_left"`
		TextFourMarginLeft    int    `json:"text_four_margin_left"`
		TextFourMarginRight   int    `json:"text_four_margin_right"`
		TextFourMarginTop     int    `json:"text_four_margin_top"`
		TextFourMarginBottom  int    `json:"text_four_margin_bottom"`
		SeparatorWidth        int    `json:"separator_width"`
		SeparatorHeight       int    `json:"separator_height"`
		SeparatorMarginLeft   int    `json:"separator_margin_left"`
		SeparatorMarginRight  int    `json:"separator_margin_right"`
		SeparatorMarginTop    int    `json:"separator_margin_top"`
		SeparatorMarginBottom int    `json:"separator_margin_bottom"`
		BorderRadius          int    `json:"border_radius"`
		Isblur                bool   `json:"is_blur"`
	}
)

var frameLayouts *FrameLayouts

// 根据名称查找布局.
func FindLayoutByName(name string) (FrameLayout, pkg.EError) {
	for i := range frameLayouts.List {
		if frameLayouts.List[i].Name == name {
			return frameLayouts.List[i], pkg.NoError
		}
	}

	return FrameLayout{}, pkg.LayoutNotFindError
}

// 根据名称查找布局.
func MustFindLayoutByName(name string) FrameLayout {
	for i := range frameLayouts.List {
		if frameLayouts.List[i].Name == name {
			return frameLayouts.List[i]
		}
	}

	return FrameLayout{}
}

// 获取全部的模板.
func GetAllLayout() []FrameLayout {
	loadandInitLayout()

	return frameLayouts.List
}

// 加载并初始化布局.
func loadandInitLayout() pkg.EError {
	// 确保frameLayouts不为nil
	if frameLayouts == nil {
		frameLayouts = &FrameLayouts{}
	}

	file := internal.GetMainLayoutPath()
	if !internal.PathExists(file) {
		return pkg.NewErrors(pkg.FILE_NOT_EXIST_ERROR, file+":布局文件不存在")
	}

	layoutStr, err := os.ReadFile(file)
	if err != nil {
		return pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, file+":布局文件打开失败")
	}

	err = json.Unmarshal(layoutStr, &frameLayouts)
	if err != nil {
		return pkg.NewErrors(pkg.FILE_NOT_READ_ERROR, file+":布局文件json解析失败")
	}

	return pkg.NoError
}

// 重新加载模板.
func ReloadandInitLayout() pkg.EError {
	frameLayouts = &FrameLayouts{}

	return loadandInitLayout()
}

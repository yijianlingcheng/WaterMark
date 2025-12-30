package native

import (
	"WaterMark/layout"
	"WaterMark/pkg"
)

type (
	borderStrategy interface {
		initLayoutValue(fm baseFrame) pkg.EError
		drawBorder(fm baseFrame) pkg.EError
	}

	borderBlurStrategy interface {
		initLayoutValue(fm baseFrame) pkg.EError
		drawBorder(fm baseFrame) pkg.EError
	}

	// 固定-模板.
	fixedBottomLogoTextLayoutBorder struct {
		Strategy borderStrategy
		baseBottomLogoTextLayoutBorder
		IsRight bool
	}

	// 自适应-模板.
	autoBottomLogoTextLayoutBorder struct {
		Strategy borderStrategy
		baseBottomLogoTextLayoutBorder
		IsRight      bool
		HasSeparator bool
	}
	// 自适应-均衡-模板.
	autoBottomLogoTextAverageLayoutBorder struct {
		Strategy borderStrategy
		autoBottomLogoTextLayoutBorder
	}

	// 简约布局-模板.
	simpleBottomLogoTextCenterBorder struct {
		Strategy borderStrategy
		baseBottomLogoTextLayoutBorder
		HasLogo bool
	}

	// 高斯模糊模板.
	blurBottomTextCenterLayout struct {
		Strategy borderBlurStrategy
		baseBottomLogoTextLayoutBorder
	}

	SimpleBorderFactory struct{}
)

// 获取模板对应的策略.
func (simple *SimpleBorderFactory) createBorder(name string) borderStrategy {
	// 根据模板名称查询配置的模板
	t := layout.MustFindLayoutByName(name)

	switch t.Type {
	// 固定-左logo-布局
	case "fixed_bottom_logo_text_left_layout":
		return &fixedBottomLogoTextLayoutBorder{}
	// 固定-右logo-布局
	case "fixed_bottom_logo_text_right_layout":
		return &fixedBottomLogoTextLayoutBorder{
			IsRight: true,
		}
	// 经典-自动-左logo-无分割线-布局
	case "auto_bottom_logo_text_left_no_separator_layout":
		return &autoBottomLogoTextLayoutBorder{}
	// 经典-自动-右logo-无分割线-布局
	case "auto_bottom_logo_text_right_no_separator_layout":
		return &autoBottomLogoTextLayoutBorder{
			IsRight: true,
		}
	// 经典-自动-左logo-布局
	case "auto_bottom_logo_text_left_layout":
		return &autoBottomLogoTextLayoutBorder{
			HasSeparator: true,
		}
	// 经典-自动-右logo-布局
	case "auto_bottom_logo_text_right_layout":
		return &autoBottomLogoTextLayoutBorder{
			IsRight:      true,
			HasSeparator: true,
		}
	case "auto_bottom_logo_text_average_layout":
		return &autoBottomLogoTextAverageLayoutBorder{}

	// 简约-居中-无logo-布局
	case "simple_bottom_text_center_layout":
		return &simpleBottomLogoTextCenterBorder{}

	// 简约-居中-布局
	case "simple_bottom_logo_text_center_layout":
		return &simpleBottomLogoTextCenterBorder{
			HasLogo: true,
		}
	case "blur_bottom_text_center_layout":
		return &blurBottomTextCenterLayout{}
	}

	// 如果没有匹配到任何布局类型，返回默认布局
	return &fixedBottomLogoTextLayoutBorder{}
}

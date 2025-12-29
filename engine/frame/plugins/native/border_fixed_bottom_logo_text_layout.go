package native

import (
	"WaterMark/pkg"
)

// 画边框.
func (b *fixedBottomLogoTextLayoutBorder) drawBorder(fm baseFrame) pkg.EError {
	// 画logo
	b.drawLogo(fm)

	// 画水印文字
	b.drawWords(fm)

	// 画分隔符
	b.drawSeparator(fm)

	return pkg.NoError
}

// 固定布局,不进行任何操作.
func (b *fixedBottomLogoTextLayoutBorder) initLayoutValue(fm baseFrame) pkg.EError {
	return pkg.NoError
}

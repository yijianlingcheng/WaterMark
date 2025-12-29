package native

import (
	"image"
	"image/draw"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"WaterMark/pkg"
)

// TextBrush 画笔.
type textBrush struct {
	FontType  *truetype.Font
	FontColor *image.Uniform
	FontSize  float64
}

// newTextBrush 初始化画笔
//
//	fontFilePath 字体文件路径
//	fontSize 字体size
//	fontColor 字体颜色
func newTextBrush(fontFilePath string, fontSize float64, fontColor *image.Uniform) (*textBrush, pkg.EError) {
	if fontFilePath == "" {
		return &textBrush{}, pkg.NoError
	}
	// 从缓存中加载字体
	fontType, err := loadTextFontWithCache(fontFilePath)
	if pkg.HasError(err) {
		return nil, err
	}

	return &textBrush{FontType: fontType, FontSize: fontSize, FontColor: fontColor}, pkg.NoError
}

// drawFontOnRGBA 图片插入文字
//
//	rgba 需要添加文字的图片
//	pt 文字对应的坐标
//	content 文字内容
func (fb *textBrush) drawFontOnRGBA(rgba draw.Image, pt image.Point, content string) pkg.EError {
	if fb.FontType == nil {
		return pkg.NoError
	}
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fb.FontType)
	c.SetHinting(font.HintingFull)
	c.SetFontSize(fb.FontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fb.FontColor)
	_, err := c.DrawString(content, freetype.Pt(pt.X+10, pt.Y+int(c.PointToFixed(fb.FontSize)>>6)))
	if err != nil {
		return pkg.NewErrors(pkg.IMAGE_TEXT_DRAW_TXT_ERROR, content+":绘制失败,原因:"+err.Error())
	}

	return pkg.NoError
}

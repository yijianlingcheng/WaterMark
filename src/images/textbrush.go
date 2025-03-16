package images

import (
	"WaterMark/src/logs"
	"image"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// TextBrush 画笔
type textBrush struct {

	// FontType 字体文件
	FontType *truetype.Font

	// FontSize 字体size
	FontSize float64

	// FontColor 字体颜色
	FontColor *image.Uniform
}

// newTextBrush 初始化画笔
//
//	@param FontFilePath 字体文件路径
//	@param FontSize 字体size
//	@param FontColor 字体颜色
//	@return *textBrush
//	@return error
func newTextBrush(FontFilePath string, FontSize float64, FontColor *image.Uniform) (*textBrush, error) {
	fontFile, err := os.ReadFile(FontFilePath)
	if err != nil {
		logs.Errors.Println(FontFilePath + "文件打开失败:" + err.Error())
		return nil, err
	}
	fontType, err := truetype.Parse(fontFile)
	if err != nil {
		logs.Errors.Println(FontFilePath + "文件解析失败:" + err.Error())
		return nil, err
	}
	return &textBrush{FontType: fontType, FontSize: FontSize, FontColor: FontColor}, nil
}

// drawFontOnRGBA 图片插入文字
//
//	@param rgba 需要添加文字的图片
//	@param pt 文字对应的坐标
//	@param content 文字内容
func (fb *textBrush) drawFontOnRGBA(rgba *image.RGBA, pt image.Point, content string) {
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fb.FontType)
	c.SetHinting(font.HintingFull)
	c.SetFontSize(fb.FontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fb.FontColor)
	c.DrawString(content, freetype.Pt(pt.X, pt.Y))
}

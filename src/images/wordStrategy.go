package images

import (
	"image"
	"sort"
)

// WordStrategy
type WordStrategy interface {
	drawWords(w *WaterMark)
}

// BottomLeftWordStrategy
type BottomLeftWordStrategy struct {
	Strategy WordStrategy
}

// drawWord implements WordStrategy.
//
//	@param w
func (b *BottomLeftWordStrategy) drawWords(w *WaterMark) {

	logoT := w.WaterMarkTemplate.LogoTemplate
	borderT := w.WaterMarkTemplate.BorderTemplate
	WordsT := w.WaterMarkTemplate.WordsTemplate
	SepT := w.WaterMarkTemplate.SeparateTamplate

	sminx := 0
	smaxy := w.SourceHeight

	x := sminx + borderT.LeftWidth + logoT.Width + WordsT.FirstFontMarginLeft
	y := smaxy + borderT.TopHeight + WordsT.FirstFontMarginTop

	// 存在分隔符,单独计算第一排文字间距
	if SepT.Exist {
		x = x + SepT.Width + SepT.MarginLeft + SepT.MarginRight
	}

	ftextBrush, _ := newTextBrush(WordsT.FirstFontFile, float64(WordsT.FirstFontSize), &image.Uniform{WordsT.FirstFontColor})
	stextBrush, _ := newTextBrush(WordsT.SecondFontFile, float64(WordsT.SecondFontSize), &image.Uniform{WordsT.SecondFontColor})

	// 第一行左边文字
	ftextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("One"))

	// 第二行左边文字
	x = sminx + borderT.LeftWidth + logoT.Width + WordsT.SecondFontMarginLeft

	// 存在分隔符,单独计算第一排文字间距
	if SepT.Exist {
		x = x + SepT.Width + SepT.MarginLeft + SepT.MarginRight
	}
	y = smaxy + borderT.TopHeight + WordsT.SecondFontMarginTop
	stextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("Two"))

	// 第一行右边文字
	x1 := w.Draw.Bounds().Dx() - borderT.RightWidth - WordsT.FirstFontMarginRight
	y = smaxy + borderT.TopHeight + WordsT.FirstFontMarginTop
	ftextBrush.drawFontOnRGBA(w.Draw, image.Pt(x1, y), w.getWords("Three"))

	// 第二行右边文字
	x2 := w.Draw.Bounds().Dx() - borderT.RightWidth - WordsT.SecondFontMarginRight
	y = smaxy + borderT.TopHeight + WordsT.SecondFontMarginTop
	stextBrush.drawFontOnRGBA(w.Draw, image.Pt(x2, y), w.getWords("Four"))
}

// BottomCenterWordStrategy
type BottomCenterWordStrategy struct {
	Strategy WordStrategy
}

// drawWord implements WordStrategy.
//
//	@param w
func (b *BottomCenterWordStrategy) drawWords(w *WaterMark) {

	logoT := w.WaterMarkTemplate.LogoTemplate
	borderT := w.WaterMarkTemplate.BorderTemplate
	WordsT := w.WaterMarkTemplate.WordsTemplate
	SepT := w.WaterMarkTemplate.SeparateTamplate

	sminx := 0
	smaxy := w.SourceHeight

	x := sminx + borderT.LeftWidth + WordsT.FirstFontMarginLeft
	y := smaxy + borderT.TopHeight + WordsT.FirstFontMarginTop

	ftextBrush, _ := newTextBrush(WordsT.FirstFontFile, float64(WordsT.FirstFontSize), &image.Uniform{WordsT.FirstFontColor})
	stextBrush, _ := newTextBrush(WordsT.SecondFontFile, float64(WordsT.SecondFontSize), &image.Uniform{WordsT.SecondFontColor})

	// 第一行左边文字
	ftextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("One"))

	// 第二行左边文字
	x = sminx + borderT.LeftWidth + WordsT.SecondFontMarginLeft
	y = smaxy + borderT.TopHeight + WordsT.SecondFontMarginTop
	stextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("Two"))

	// 第一行右边文字
	x = w.Draw.Bounds().Dx() - borderT.RightWidth - logoT.MarginRight + logoT.Width

	// 存在分隔符,单独计算第一排文字间距
	if SepT.Exist {
		x = x + SepT.Width + SepT.MarginLeft + SepT.MarginRight
	}

	y = smaxy + borderT.TopHeight + WordsT.FirstFontMarginTop
	ftextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("Three"))

	// 第二行右边文字
	y = smaxy + borderT.TopHeight + WordsT.SecondFontMarginTop
	stextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("Four"))
}

// BottomRightWordStrategy
type BottomRightWordStrategy struct {
	Strategy WordStrategy
}

// drawWord implements WordStrategy.
//
//	@param w
func (b *BottomRightWordStrategy) drawWords(w *WaterMark) {

	logoT := w.WaterMarkTemplate.LogoTemplate
	borderT := w.WaterMarkTemplate.BorderTemplate
	WordsT := w.WaterMarkTemplate.WordsTemplate
	SepT := w.WaterMarkTemplate.SeparateTamplate

	sminx := 0
	smaxy := w.SourceHeight

	x := sminx + borderT.LeftWidth + WordsT.FirstFontMarginLeft
	y := smaxy + borderT.TopHeight + WordsT.FirstFontMarginTop

	ftextBrush, _ := newTextBrush(WordsT.FirstFontFile, float64(WordsT.FirstFontSize), &image.Uniform{WordsT.FirstFontColor})
	stextBrush, _ := newTextBrush(WordsT.SecondFontFile, float64(WordsT.SecondFontSize), &image.Uniform{WordsT.SecondFontColor})

	// 第一行左边文字
	ftextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("One"))

	// 第二行左边文字
	x = sminx + borderT.LeftWidth + WordsT.SecondFontMarginLeft
	y = smaxy + borderT.TopHeight + WordsT.SecondFontMarginTop
	stextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("Two"))

	// 第一行右边文字
	x = w.Draw.Bounds().Dx() - borderT.RightWidth - logoT.Width - logoT.MarginLeft - WordsT.FirstFontMarginRight

	// 存在分隔符,单独计算第一排文字间距
	if SepT.Exist {
		x = x - SepT.MarginRight - SepT.Width - SepT.MarginLeft
	}

	y = smaxy + borderT.TopHeight + WordsT.FirstFontMarginTop
	ftextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("Three"))

	// 第二行右边文字
	x = w.Draw.Bounds().Dx() - borderT.RightWidth - logoT.Width - logoT.MarginLeft - WordsT.SecondFontMarginRight

	// 存在分隔符,单独计算文字间距
	if SepT.Exist {
		x = x - SepT.MarginRight - SepT.Width - SepT.MarginLeft
	}
	y = smaxy + borderT.TopHeight + WordsT.SecondFontMarginTop
	stextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("Four"))
}

// StackblurWordStrategy
type StackblurWordStrategy struct {
	Strategy WordStrategy
}

// drawWords
//
//	@param w
func (b *StackblurWordStrategy) drawWords(w *WaterMark) {
	logoT := w.WaterMarkTemplate.LogoTemplate
	borderT := w.WaterMarkTemplate.BorderTemplate
	WordsT := w.WaterMarkTemplate.WordsTemplate
	SepT := w.WaterMarkTemplate.SeparateTamplate

	sminx := 0
	smaxy := w.Draw.Bounds().Dy()

	x := sminx + borderT.LeftWidth + logoT.MarginLeft + logoT.Width + logoT.MarginRight + WordsT.FirstFontMarginLeft
	y := smaxy - borderT.BottomHeight + WordsT.FirstFontMarginTop

	ftextBrush, _ := newTextBrush(WordsT.FirstFontFile, float64(WordsT.FirstFontSize), &image.Uniform{WordsT.FirstFontColor})
	stextBrush, _ := newTextBrush(WordsT.SecondFontFile, float64(WordsT.SecondFontSize), &image.Uniform{WordsT.SecondFontColor})

	// 第一行左边文字
	ftextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("One"))

	// 第一行右边文字
	x = sminx + borderT.LeftWidth + logoT.MarginLeft + logoT.Width + logoT.MarginRight + WordsT.SecondFontMarginLeft
	if SepT.Exist {
		x = x + SepT.MarginLeft + SepT.Width + SepT.MarginRight
	}
	y = smaxy - borderT.BottomHeight + WordsT.FirstFontMarginTop

	stextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("Two"))
}

// SimpleWordFactory
type SimpleWordFactory struct {
}

// create
//
//	@param ext
//	@return WordStrategy
func (simple *SimpleWordFactory) create(ext string) WordStrategy {
	switch ext {
	case "BOTTOM_LOGO_LEFT":
		return &BottomLeftWordStrategy{}
	case "BOTTOM_LOGO_CENTER":
		return &BottomCenterWordStrategy{}
	case "BOTTOM_LOGO_RIGHT":
		return &BottomRightWordStrategy{}
	case "STACK_BLUR":
		return &StackblurWordStrategy{}
	case "BOTTOM_LOGO_LEFT_AUTO":
		return &BottomLeftWordAutoStrategy{}
	case "BOTTOM_LOGO_CENTER_AUTO":
		return &BottomLeftWordAutoStrategy{}
	case "BOTTOM_LOGO_RIGHT_AUTO":
		return &BottomLeftWordAutoStrategy{}
	case "STACK_BLUR_AUTO":
		return &BottomLeftWordAutoStrategy{}
	}
	return nil
}

// BottomLeftWordAutoStrategy
type BottomLeftWordAutoStrategy struct {
	Strategy WordStrategy
}

func (b *BottomLeftWordAutoStrategy) drawWords(w *WaterMark) {
	// 计算边距,字体
	b.calculateLeftAutoWordsT(w)

	logoT := w.WaterMarkTemplate.LogoTemplate
	borderT := w.WaterMarkTemplate.BorderTemplate
	WordsT := w.WaterMarkTemplate.WordsTemplate

	sminx := 0
	smaxy := w.SourceHeight

	x := sminx + borderT.LeftWidth + logoT.Width + WordsT.FirstFontMarginLeft
	y := smaxy + borderT.TopHeight + WordsT.FirstFontMarginTop

	ftextBrush, _ := newTextBrush(WordsT.FirstFontFile, float64(WordsT.FirstFontSize), &image.Uniform{WordsT.FirstFontColor})
	stextBrush, _ := newTextBrush(WordsT.SecondFontFile, float64(WordsT.SecondFontSize), &image.Uniform{WordsT.SecondFontColor})

	// 第一行左边文字
	ftextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("One"))

	// 第二行左边文字
	x = sminx + borderT.LeftWidth + logoT.Width + WordsT.SecondFontMarginLeft
	y = smaxy + borderT.TopHeight + WordsT.SecondFontMarginTop
	stextBrush.drawFontOnRGBA(w.Draw, image.Pt(x, y), w.getWords("Two"))

	// 第一行右边文字
	x1 := w.Draw.Bounds().Dx() - borderT.RightWidth - WordsT.FirstFontMarginRight
	y = smaxy + borderT.TopHeight + WordsT.FirstFontMarginTop
	ftextBrush.drawFontOnRGBA(w.Draw, image.Pt(x1, y), w.getWords("Three"))

	// 第二行右边文字
	x2 := w.Draw.Bounds().Dx() - borderT.RightWidth - WordsT.SecondFontMarginRight
	y = smaxy + borderT.TopHeight + WordsT.SecondFontMarginTop
	stextBrush.drawFontOnRGBA(w.Draw, image.Pt(x2, y), w.getWords("Four"))
}

// calculateLeftAutoWordsT 计算字体边距,字体大小
//
//	@param w
func (b *BottomLeftWordAutoStrategy) calculateLeftAutoWordsT(w *WaterMark) {
	// 边框模板
	borderT := w.WaterMarkTemplate.BorderTemplate

	// 字体大小
	fontSize := borderT.BottomHeight / 4

	// 右边距
	marginRightA := []int{len(w.getWords("Three")), len(w.getWords("Four"))}
	sort.Ints(marginRightA)
	marginRight := int(float64(marginRightA[len(marginRightA)-1]*fontSize) * 0.6)

	// 上边距
	marginTop := int(float64(borderT.BottomHeight) / 2.5)

	// 对对象赋值,方便后续计算
	w.WaterMarkTemplate.WordsTemplate = newWordsTemplate().WithFontSize(fontSize).WithMarginRight(marginRight).WithMarginTop(marginTop)
}

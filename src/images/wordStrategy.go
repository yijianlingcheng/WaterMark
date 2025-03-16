package images

import (
	"image"
	"sort"
)

// wordStrategy
type wordStrategy interface {
	drawWords(w *WaterMark)
}

// bottomLeftWord
type bottomLeftWord struct {
	Strategy wordStrategy
}

// drawWord implements WordStrategy.
//
//	@param w
func (b *bottomLeftWord) drawWords(w *WaterMark) {

	logoT := w.WT.LogoT
	borderT := w.WT.BorderT
	WordsT := w.WT.WordsT
	SepT := w.WT.SeparateT

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

// bottomCenterWord
type bottomCenterWord struct {
	Strategy wordStrategy
}

// drawWord implements WordStrategy.
//
//	@param w
func (b *bottomCenterWord) drawWords(w *WaterMark) {

	logoT := w.WT.LogoT
	borderT := w.WT.BorderT
	WordsT := w.WT.WordsT
	SepT := w.WT.SeparateT

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

// bottomRightWord
type bottomRightWord struct {
	Strategy wordStrategy
}

// drawWord implements WordStrategy.
//
//	@param w
func (b *bottomRightWord) drawWords(w *WaterMark) {

	logoT := w.WT.LogoT
	borderT := w.WT.BorderT
	WordsT := w.WT.WordsT
	SepT := w.WT.SeparateT

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

// stackblurWord
type stackblurWord struct {
	Strategy wordStrategy
}

// drawWords
//
//	@param w
func (b *stackblurWord) drawWords(w *WaterMark) {
	logoT := w.WT.LogoT
	borderT := w.WT.BorderT
	WordsT := w.WT.WordsT
	SepT := w.WT.SeparateT

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
//	@return wordStrategy
func (simple *SimpleWordFactory) create(ext string) wordStrategy {
	switch ext {
	case "BOTTOM_LOGO_LEFT":
		return &bottomLeftWord{}
	case "BOTTOM_LOGO_CENTER":
		return &bottomCenterWord{}
	case "BOTTOM_LOGO_RIGHT":
		return &bottomRightWord{}
	case "STACK_BLUR":
		return &stackblurWord{}
	case "BOTTOM_LOGO_LEFT_AUTO":
		return &bottomLeftWordAuto{}
	case "BOTTOM_LOGO_CENTER_AUTO":
		return &bottomLeftWordAuto{}
	case "BOTTOM_LOGO_RIGHT_AUTO":
		return &bottomLeftWordAuto{}
	case "STACK_BLUR_AUTO":
		return &bottomLeftWordAuto{}
	}
	return nil
}

// bottomLeftWordAuto
type bottomLeftWordAuto struct {
	Strategy wordStrategy
}

func (b *bottomLeftWordAuto) drawWords(w *WaterMark) {
	// 计算边距,字体
	b.calculateLeftAutoWordsT(w)

	logoT := w.WT.LogoT
	borderT := w.WT.BorderT
	WordsT := w.WT.WordsT

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
func (b *bottomLeftWordAuto) calculateLeftAutoWordsT(w *WaterMark) {
	// 边框模板
	borderT := w.WT.BorderT

	// 字体大小
	fontSize := borderT.BottomHeight / 4

	// 右边距
	marginRightA := []int{len(w.getWords("Three")), len(w.getWords("Four"))}
	sort.Ints(marginRightA)
	marginRight := int(float64(marginRightA[len(marginRightA)-1]*fontSize) * 0.6)

	// 上边距
	marginTop := int(float64(borderT.BottomHeight) / 2.5)

	// 对对象赋值,方便后续计算
	w.WT.WordsT = newWordsTemplate().WithFontSize(fontSize).WithMarginRight(marginRight).WithMarginTop(marginTop)
}

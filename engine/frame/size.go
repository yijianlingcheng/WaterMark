package frame

type PhotoSize struct {
	BorderLeftWidth    int
	BorderRightWidth   int
	BorderTopHeight    int
	BorderBottomHeight int
	SourceWidth        int
	SourceHeight       int
	Isblur             int
	BorderRadius       int
}

// 返回一个图片尺寸.
func NewPhotoSize(size map[string]int) PhotoSize {
	borderLeftWidth, blwOk := size["borderLeftWidth"]
	if !blwOk {
		borderLeftWidth = 0
	}
	borderRightWidth, brwOk := size["borderRightWidth"]
	if !brwOk {
		borderRightWidth = 0
	}
	borderTopHeight, bthOk := size["borderTopHeight"]
	if !bthOk {
		borderTopHeight = 0
	}
	borderBottomHeight, bbhOk := size["borderBottomHeight"]
	if !bbhOk {
		borderBottomHeight = 0
	}
	sourceWidth, swOk := size["sourceWidth"]
	if !swOk {
		sourceWidth = 0
	}
	sourceHeight, shOk := size["sourceHeight"]
	if !shOk {
		sourceHeight = 0
	}
	isBlur, ibOk := size["isBlur"]
	if !ibOk {
		isBlur = 0
	}
	borderRadius, brOk := size["borderRadius"]
	if !brOk {
		borderRadius = 0
	}

	return PhotoSize{
		BorderLeftWidth:    borderLeftWidth,
		BorderRightWidth:   borderRightWidth,
		BorderTopHeight:    borderTopHeight,
		BorderBottomHeight: borderBottomHeight,
		SourceWidth:        sourceWidth,
		SourceHeight:       sourceHeight,
		Isblur:             isBlur,
		BorderRadius:       borderRadius,
	}
}

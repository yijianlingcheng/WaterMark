package native

import (
	"image"
	"image/color"
)

type borderRadius struct {
	p image.Point // 矩形右下角位置
	r int
}

func (c *borderRadius) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *borderRadius) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.p.X, c.p.Y)
}

// 对每个像素点进行色值设置，分别处理矩形的四个角，在四个角的内切圆的外侧，色值设置为全透明，其他区域不透明.
func (c *borderRadius) At(x, y int) color.Color {
	var xx, yy, rr float64

	// 左上
	if x <= c.r && y <= c.r {
		xx, yy, rr = float64(c.r-x)+0.5, float64(y-c.r)+0.5, float64(c.r)

		if xx*xx+yy*yy >= rr*rr {
			return color.Alpha{}
		}

		return color.Alpha{A: 255}
	}
	// 右上
	if x >= c.p.X-c.r && y <= c.r {
		xx, yy, rr = float64(x-(c.p.X-c.r))+0.5, float64(y-c.r)+0.5, float64(c.r)

		if xx*xx+yy*yy >= rr*rr {
			return color.Alpha{}
		}

		return color.Alpha{A: 255}
	}
	// 左下
	if x <= c.r && y >= c.p.Y-c.r {
		xx, yy, rr = float64(c.r-x)+0.5, float64(y-(c.p.Y-c.r))+0.5, float64(c.r)

		if xx*xx+yy*yy >= rr*rr {
			return color.Alpha{}
		}

		return color.Alpha{A: 255}
	}
	// 右下
	if x >= c.p.X-c.r && y >= c.p.Y-c.r {
		xx, yy, rr = float64(x-(c.p.X-c.r))+0.5, float64(y-(c.p.Y-c.r))+0.5, float64(c.r)

		if xx*xx+yy*yy >= rr*rr {
			return color.Alpha{}
		}

		return color.Alpha{A: 255}
	}

	return color.Alpha{A: 255}
}

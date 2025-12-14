package native

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/fogleman/gg"

	"WaterMark/message"
	"WaterMark/pkg"
)

// 文字内容宽高缓存.
type textContentXAndYCache struct {
	xCache map[string]int
	yCache map[string]int
	mtx    sync.Mutex
}

var (
	textContentCache     *textContentXAndYCache
	textContentCacheOnce sync.Once
)

// 获取文字内容对应的width,每次都需要重新计算.
func getTextContentSize(fontSize int, fontFile, content string) (int, int) {
	// 利用gg库计算文字宽度
	width, height := fontSize*len(content), fontSize*2
	dc := gg.NewContext(width, height)
	err := dc.LoadFontFace(fontFile, float64(fontSize))
	if err != nil {
		message.SendErrorMsg(fontFile + ":字体文件不存在")

		return 0, 0
	}
	w, h := dc.MeasureString(content)

	return int(w), int(h)
}

// 获取指定宽度,指定字体文件下的最大宽度.
func getTextContentMaxSize(width int, content string) int {
	maxFontSize := width / len(content)
	for range 2 {
		maxFontSize = maxFontSize * 96 / 72
	}

	return maxFontSize
}

// 获取文字内容对应的width.
func getTextContentXAndY(fontSize int, fontFile, content string) (int, int) {
	// 延迟初始化
	textContentCacheOnce.Do(func() {
		textContentCache = &textContentXAndYCache{
			xCache: make(map[string]int),
			yCache: make(map[string]int),
		}
	})

	// 计算cache key
	key := pkg.GetStrMD5(fmt.Sprintf("%d%s%s", fontSize, fontFile, content))
	// 取数据
	textContentCache.mtx.Lock()
	width, xok := textContentCache.xCache[key]
	height, yok := textContentCache.yCache[key]
	textContentCache.mtx.Unlock()

	if xok && yok {
		return width, height
	}
	width, height = getTextContentSize(fontSize, fontFile, content)

	// 写入缓存
	textContentCache.mtx.Lock()
	textContentCache.xCache[key] = width
	textContentCache.yCache[key] = height
	textContentCache.mtx.Unlock()

	return width, height
}

// 画边框上的logo.
func drawBorderLogo(fm *photoFrame, logoImage image.Image, startX, startY, endX, endY int) {
	draw.Draw(
		fm.borderDraw,
		image.Rect(startX, startY, endX, endY),
		logoImage,
		image.Pt(0, 0),
		draw.Src,
	)
}

// 字符串颜色转RGBA.
func strColor2RGBA(s string) color.RGBA {
	if s == "" {
		s = COLOR
	}
	list := strings.Split(s, ",")
	r0, _ := strconv.ParseUint(list[0], 10, 8)
	r1, _ := strconv.ParseUint(list[1], 10, 8)
	r2, _ := strconv.ParseUint(list[2], 10, 8)
	r3, _ := strconv.ParseUint(list[3], 10, 8)

	return color.RGBA{uint8(r0), uint8(r1), uint8(r2), uint8(r3)}
}

// 取绝对值.
func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

// 画线.
func drawLine(img draw.Image, start, end image.Point, c color.Color) {
	dx := abs(end.X - start.X)
	dy := abs(end.Y - start.Y)
	sx, sy := 0, 0
	if start.X < end.X {
		sx = 1
	} else if start.X > end.X {
		sx = -1
	}
	if start.Y < end.Y {
		sy = 1
	} else if start.Y > end.Y {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(start.X, start.Y, c)
		if start.X == end.X && start.Y == end.Y {
			break
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			start.X += sx
		}
		if e2 < dx {
			err += dx
			start.Y += sy
		}
	}
}

// 保存图片.
func saveImage(saveImageFile string, image *image.RGBA, quality int) {
	file, err := os.Create(saveImageFile)
	if err != nil {
		message.SendErrorMsg(saveImageFile + ":图片打开失败")

		return
	}
	defer file.Close()

	err = jpeg.Encode(file, image, &jpeg.Options{
		Quality: quality,
	})
	if err != nil {
		message.SendErrorMsg(saveImageFile + "图片写入失败:" + err.Error())
	}

	runtime.GC()
}

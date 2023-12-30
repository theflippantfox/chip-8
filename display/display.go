package display

import (
	tm "github.com/buger/goterm"
)

const Height = 32
const Width = 64

type Display struct {
	Buffer [Height * Width]uint8
}

func NewDisplay() *Display {
	return &Display{}
}

func ClearDisplay(gfx *Display) {
	for i := 0; i < Height*Width; i++ {
		gfx.Buffer[i] = 0
	}
	tm.Clear()
}

func FetchPixel(gfx *Display, x uint16, y uint16) uint8 {
	return gfx.Buffer[x+(y*64)]
}

func SetPixel(gfx *Display, x uint16, y uint16, val uint8) {
	gfx.Buffer[x+(y*64)] = val
}

func XORPixel(gfx *Display, x uint16, y uint16) {
	gfx.Buffer[x+(y*64)] ^= 1
}

func Render(gfx *Display) {
	tm.MoveCursor(1, 1)
	for i := 0; i < Height*Width; i++ {
		if gfx.Buffer[i] == 1 {
			tm.Print("*")
		} else {
			tm.Print(" ")
		}

		if i > 0 && i%64 == 0 {
			tm.Println()
		}

	}
	tm.Flush()
}

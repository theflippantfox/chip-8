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

func (gfx *Display) ClearDisplay() {
	for i := 0; i < Height*Width; i++ {
		gfx.Buffer[i] = 0
	}
	tm.Clear()
}

func (gfx *Display) FetchPixel(x uint16, y uint16) uint8 {
	return gfx.Buffer[x+(y*Width)]
}

func (gfx *Display) SetPixel(x uint16, y uint16, val uint8) {
	gfx.Buffer[x+(y*Width)] = val
}

func (gfx *Display) XORPixel(x uint16, y uint16) {
	gfx.Buffer[x+(y*Width)] ^= 1
}

func (gfx *Display) Render() {
	tm.MoveCursor(1, 1)
	for i := 0; i < Height*Width; i++ {
		if gfx.Buffer[i] == 1 {
			tm.Print("*")
		} else {
			tm.Print(" ")
		}

		if i > 0 && i%Width == 0 {
			tm.Println()
		}

	}
	tm.Flush()
}

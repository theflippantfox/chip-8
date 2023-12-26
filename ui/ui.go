package ui

import (
	"fmt"
	"time"

	tm "github.com/buger/goterm"
)

func Renderer(framebuffer [32 * 64]bool) {
	tm.MoveCursor(1, 2) //Using 2 to avoild overlapping with the prompt
	for i := 0; i < 32*64; i++ {
		if framebuffer[i] {
			fmt.Print("*")
		} else {
			fmt.Print(" ")
		}

		if i%64 == 0 {
			fmt.Print("\n")
		}
	}
	tm.Flush()
	time.Sleep(time.Second)
}

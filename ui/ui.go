package ui

import (
	"fmt"
	"time"

	tm "github.com/buger/goterm"
)

func InitGUI() {
	screen := [32][64]bool{}
	for {
		tm.Clear()

		tm.MoveCursor(1, 1)

		for i := 0; i < 32; i++ {
			for j := 0; j < 64; j++ {
				screen[i][j] = false
			}
		}

		for i := 0; i < 32; i++ {
			for j := 0; j < 64; j++ {
				if screen[i][j] == false {
					fmt.Print(".")
				} else {
					fmt.Print("*")
				}
			}

			fmt.Println("")
		}
		tm.Flush() // Call it every time at the end of rendering
		time.Sleep(5 * time.Second)
	}
}

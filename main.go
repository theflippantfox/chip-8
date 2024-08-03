package main

import (
	"chip8/cmd/cpu"
	"chip8/cmd/display"
	"chip8/cmd/memory"
	"os"
)

func main() {
	c := cpu.NewChip()
	m := mem.NewMemory()
	gfx := display.NewDisplay()

	roamArg := os.Args[1]
	roam := "tests/bin/3-corax+.ch8"
	if roamArg != "" {
		roam = roamArg
	}

	m.LoadROMtoMemory(roam)

	for i := 0; i == i; i++ {
		c.EmulateCycle(m, gfx)
	}
}

package main

import (
	"chip8/cpu"
	"chip8/display"
	"chip8/memory"
)

func main() {
	c := cpu.NewChip()
	m := mem.NewMemory()
	gfx := display.NewDisplay()

	m.LoadROMtoMemory("tests/bin/3-corax+.ch8")

	for i := 0; i == i; i++ {
		c.EmulateCycle(m, gfx)
	}
}

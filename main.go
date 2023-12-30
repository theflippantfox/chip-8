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

	mem.LoadROMtoMemory(m, "roms/games/Timebomb.ch8")
    
    for i:=0; i==i; i++ {
		cpu.Execute(c, m, gfx)
    }

}

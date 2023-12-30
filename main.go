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

	mem.LoadROMtoMemory(m, "roms/programs/Chip8 emulator Logo [Garstyciuks].ch8")
    
    for i:=0; i==i; i++ {
		cpu.EmulateCycle(c, m, gfx)
    }
}

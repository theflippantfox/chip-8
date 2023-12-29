package main

import (
	"chip8/cpu"
	"chip8/memory"
	"fmt"
)

func main() {
	c := cpu.NewChip()
	m := mem.NewMemory()

	mem.LoadROMtoMemory(m, "roms/programs/Chip8 emulator Logo [Garstyciuks].ch8")
    // fmt.Println(m.Memory)
    
    for i:=0; i==i; i++ {
        fmt.Println("Cycle: ", i) 
		cpu.Execute(c, m)
    }

}

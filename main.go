package main

import (
	"chip8/cpu"
	"chip8/memory"
	"fmt"
)

func main() {
	c := cpu.NewChip()
	m := mem.NewMemory()

	mem.LoadROMtoMemory(m, "roms/demos/Trip8 Demo (2008) [Revival Studios].ch8")

    for i:=0; i==i; i++ {
        fmt.Println("Cycle: ", i) 
		cpu.Execute(c, m)
    }
}

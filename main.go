package main

import (
	cpu "chip8/cpu"
)

func main() {
	chip8 := cpu.InitChip8()

	cpu.Cycle(&chip8)
    cpu.DumpMemory(&chip8)
}

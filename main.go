package main

import (
	cpu "chip8/cpu"
	"time"
	"fmt"
)

func main() {
	chip8 := cpu.ChipContext{}
	success := cpu.InitChip8(&chip8)

	if success {
		fmt.Println("Successfully Initialized Chip-8")
		fmt.Println(chip8.PC, cpu.GetOpCode(&chip8))
	}

	for {
		cpu.Cycle(&chip8)
	    time.Sleep(time.Second)
    }

}

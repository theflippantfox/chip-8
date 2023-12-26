package main

import (
	//"bufio"
	//"os"
    cpu "chip8/cpu"
    ui "chip8/ui"
)

func main() {
	//step := 'y'
	for {
		ui.Renderer(cpu.Frame())

		//reader := bufio.NewReader(os.Stdin)
		//step, _, _ = reader.ReadRune()
	}
}

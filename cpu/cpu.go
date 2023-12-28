package cpu

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type ChipContext struct {
	opcode uint16 //Stores the current opcode

	memory [4096]uint8 // Chip-8 had 4Kilobytes (4096 Bytes) of memory

	V [16]uint8 // 16 8-Bit registers from V0 to VF

	stack [16]uint16 // Stack

	// Special Registers
	I         uint16 // Address Pointer
	PC        uint16 // Program Counter
	SP        uint8  // Stack Pointer
	delay_reg uint8  // Delay timer
	sound_reg uint8  // Sound timer

	framebuffer [64 * 32]bool // Display framebuffer top-left to bottom-right
	keys        []uint8
	fontset     [80]uint8
}

func InitChip8() ChipContext {
	c := ChipContext{}
	Reset(&c)

	loadROM(&c)
	return c
}

func Reset(chip8 *ChipContext) {
	// Variable Declarations
	entry_point := 0x200 // ROMs are loaded at memory location 0x200/512

	// Set Defaults
	chip8.PC = uint16(entry_point) // Program Counter should start at the starting point of ROM

	// All registers should be at 0/false state
	chip8.opcode = 0x00E0
	chip8.I = 0
	chip8.SP = 0
	chip8.delay_reg = 0
	chip8.sound_reg = 0

	for i := 0; i < len(chip8.memory); i++ {
		chip8.memory[i] = 0
	}

	for i := 0; i < len(chip8.V); i++ {
		chip8.V[i] = 0
	}

	for i := 0; i < len(chip8.stack); i++ {
		chip8.stack[i] = 0
	}

	for i := 0; i < len(chip8.framebuffer); i++ {
		chip8.framebuffer[i] = false
	}

	set_fontset(chip8)
}

func loadROM(chip8 *ChipContext) {
	// Open the file in read-only mode
	file, err := os.Open("roms/demos/Zero Demo [zeroZshadow, 2007].ch8")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			fmt.Println(string(buf[:n]))
		}
	}

	fi, err := file.Stat()
	if err != nil {
		panic(err)
	}

	for i := 0; i < int(fi.Size()); i++ {
		chip8.memory[i+0x200] = buf[i]
	}
}

func DumpMemory(c *ChipContext) {
	fmt.Println(c.memory)
}

func increment_PC(chip8 *ChipContext) {
	chip8.PC += 2
}

func Cycle(chip8 *ChipContext) {
	chip8.opcode = uint16(chip8.memory[chip8.PC]<<8 | chip8.memory[chip8.PC+1])

    fmt.Println(chip8.opcode)
	instruction_set(chip8)

	increment_PC(chip8)
}

func instruction_set(chip8 *ChipContext) {
	first_digit_of_opcode := chip8.opcode >> 12

	switch first_digit_of_opcode {
	case 0x0: // When opcode starts with 0x0
		{
			if chip8.opcode == 0x00E0 { // If opcode = 0x00E0, clear the framebuffer
				for i := 0; i < len(chip8.framebuffer); i++ {
					chip8.framebuffer[i] = true
				}
			} else if chip8.opcode == 0x00EE { // If opode = 0x00EE, set the PC to the address at the top of the stack, then subtracts 1 from the SP.
				chip8.SP -= 1
				chip8.PC = chip8.stack[chip8.SP]
			}
		}
	case 0x1: // If opcode = 0x1nnn, set the program counter to nnn.

		{
			chip8.PC = chip8.opcode & 0x0FFF
		}

	case 0x2: // The interpreter increments the SP, then puts the current PC on the top of the stack. The PC is then set to nnn.
		{
			chip8.stack[chip8.SP] = chip8.PC
			chip8.SP += 1
			chip8.PC = chip8.opcode & 0x0FFF
		}
	case 0x3: // Skip next instruction if Vx = kk.
		{
			x := (chip8.opcode & 0x0F00) >> 8

			if chip8.V[x] == uint8(chip8.opcode)&0x00FF {
				increment_PC(chip8)
			}
		}
	case 0x4: // Skip next instruction if Vx != kk
		{
			x := (chip8.opcode & 0x0F00) >> 8

			if chip8.V[x] != uint8(chip8.opcode)&0x00FF {
				increment_PC(chip8)
			}
		}
	case 0x5: // Skip next instruction if Vx = Vy.
		{
			x := (chip8.opcode & 0x0F00) >> 8
			y := (chip8.opcode & 0x00F0) >> 4

			if chip8.V[x] == chip8.V[y] {
				increment_PC(chip8)
			}

		}

	case 0x6: //Set Vx = kk.
		{
			x := (chip8.opcode & 0xF00) >> 8
			kk := byte(chip8.opcode)

			chip8.V[x] = kk
		}

	case 0x7: //Set Vx = Vx + kk.
		{
			x := (chip8.opcode & 0xF00) >> 8
			kk := byte(chip8.opcode)

			chip8.V[x] += kk
		}

	case 0x8:
		{
			x := (chip8.opcode & 0x0F00) >> 8
			y := (chip8.opcode & 0x00F0) >> 4
			m := chip8.opcode & 0x000F

			switch m {
			case 0:
				chip8.V[x] = chip8.V[y]
			case 1:
				chip8.V[x] |= chip8.V[y]
			case 2:
				chip8.V[x] &= chip8.V[y]
			case 3:
				chip8.V[x] ^= chip8.V[y]
			case 4:
				{
					sum := chip8.V[x]
					sum += chip8.V[y]

					if sum > 255 {
						chip8.V[0xF] = 1
					} else {
						chip8.V[0xF] = 0
					}

					chip8.V[x] = byte(sum)
				}
			case 5:
				{
					if chip8.V[x] > chip8.V[y] {
						chip8.V[0xF] = 1
					} else {
						chip8.V[0xF] = 0
					}

					chip8.V[x] -= chip8.V[y]
				}
			case 6:
				{
					if chip8.V[x]&0x01 == 0x01 {
						chip8.V[0xF] = 1
					} else {
						chip8.V[0xF] = 0
					}

					chip8.V[x] /= 2
				}
			case 7:
				{
					if chip8.V[y] > chip8.V[x] {
						chip8.V[0xF] = 1
					} else {
						chip8.V[0xF] = 0
					}

					chip8.V[x] = chip8.V[y] - chip8.V[x]
				}
			case 14:
				{
					if chip8.V[x]&0x80 == 0x80 {
						chip8.V[0xF] = 1
					} else {
						chip8.V[0xF] = 0
					}

					chip8.V[x] *= 2
				}
			}
		}

	case 0x9:
		{
			x := (chip8.opcode & 0x0F00) >> 8
			y := (chip8.opcode & 0x00F0) >> 4

			if chip8.V[x] != chip8.V[y] {
				increment_PC(chip8)
			}

		}
	case 0xA:
		{
			chip8.I = chip8.opcode & 0x0FFF
			increment_PC(chip8)
		}
	case 0xB:
		{
			chip8.PC = (chip8.opcode & 0x0FFF) + uint16(chip8.V[0])
		}
	case 0xC: //Set Vx = random byte AND kk.
		{
			x := (chip8.opcode & 0x0F00) >> 8
			kk := byte(chip8.opcode)

			chip8.V[x] = kk + byte(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(255))
		}
	case 0xD:
		{
			chip8.V[0xF] = 0

			xx := (chip8.opcode & 0x0F00) >> 8
			yy := (chip8.opcode & 0x00F0) >> 4
			nn := (chip8.opcode & 0x000F)

			valX := chip8.V[xx]
			valY := chip8.V[yy]

			var j uint16 = 0
			var i uint16 = 0

			for j = 0; j < nn; j++ {
				pixel := chip8.memory[chip8.I+j]
				for i = 0; i < 8; i++ {
					msb := 0x80

					if pixel&uint8(msb>>i) != 0 {
						tX := (valX + uint8(i)) % 64
						tY := (valY + uint8(j)) & 32

						idx := tX + (tY * 64)

						if chip8.framebuffer[idx] == false {
							chip8.framebuffer[idx] = true
						}

					}
				}
			}
		}
	case 0xE:
		{
			x := (chip8.opcode & 0x0F00) >> 8
			kk := chip8.opcode & 0x00FF

			if kk == 0x9E {
				if chip8.keys[chip8.V[x]] == 1 {
					increment_PC(chip8)
				} else if kk == 0xA1 {
					if chip8.keys[chip8.V[x]] != 1 {
						increment_PC(chip8)
					}
				}
			}
		}
	case 0xF:
		{
			x := (chip8.opcode & 0x0F00) >> 8
			kk := chip8.opcode & 0x00FF

			if kk == 0x07 {
				chip8.V[x] = chip8.delay_reg
			} else if kk == 0x0A {
				key_pressed := false

				for i := 0; i < len(chip8.keys); i++ {
					if chip8.keys[i] != 0 {
						chip8.V[x] = uint8(i)
						key_pressed = true
						break
					}
				}

				if !key_pressed {
					return
				}
			} else if kk == 0x15 {
				chip8.delay_reg = chip8.V[x]
			} else if kk == 0x18 {
				chip8.sound_reg = chip8.V[x]
			} else if kk == 0x1E {
				chip8.I += uint16(chip8.V[x])
			} else if kk == 0x29 {
				chip8.I = uint16(chip8.V[x]) * 0x5
			} else if kk == 0x33 {
				chip8.memory[chip8.I] = chip8.V[x] / 100
				chip8.memory[chip8.I+1] = (chip8.V[x] / 10) % 10
				chip8.memory[chip8.I+2] = chip8.V[x] % 100
			} else if kk == 0x55 {
				var i uint16 = 0
				for i = 0; i < x; i++ {
					chip8.memory[chip8.I+i] = chip8.V[i]
				}
			} else if kk == 0x65 {
				var i uint16 = 0
				for i = 0; i < x; i++ {
					chip8.V[i] = chip8.memory[chip8.I+i]
				}
			}
		}
	}

}

func set_fontset(chip8 *ChipContext) { // Sets fonts for the chip
	chip8.fontset = [80]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
}

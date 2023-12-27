package cpu

import (
// "chip8/ui"
)

type ChipContext struct {
	opcode uint16

	// Chip-8 had 4Kilobytes (4096 Bytes) of memory
	memory [4096]uint8

	// 16 8-Bit registers from V0 to VF
	V [16]uint8

	// Stack
	stack [16]uint16

	// Special Registers
	I         uint16 // Address Pointer
	PC        uint16 // Program Counter
	SP        uint8  // Stack Pointer
	delay_reg uint8  // Delay timer
	sound_reg uint8  // Sound timer

	// Display framebuffer top-left to bottom-right
	framebuffer [64 * 32]bool

	fontset [80]uint8
}

func InitChip8(chip8 *ChipContext) bool {
	// Variable Declarations
	entry_point := 0x200 // ROMs are loaded at memory location 0x200/512

	// TODO: Load ROM into memory

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

	return true
}

func increment_PC(chip8 *ChipContext) {
	chip8.PC += 2
}

func GetOpCode(chip8 *ChipContext) uint16 {
	return chip8.opcode
}

func Cycle(chip8 *ChipContext) {
	chip8.opcode = uint16(chip8.memory[chip8.PC]<<8 | chip8.memory[chip8.PC+1])

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
	case 0x1:
		{
			chip8.PC = chip8.opcode & 0x0FFF // If opcode = 0x1nnn, set the program counter to nnn.
		}

	case 0x2:
		{
			chip8.stack[chip8.SP] = chip8.PC
			chip8.SP += 1
			chip8.PC = chip8.opcode & 0x0FFF
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

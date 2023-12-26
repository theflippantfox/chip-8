package cpu

const Byte = 8
const Kilobytes = 1024 * Byte

type ChipContext struct {
	// Chip-8 had 4Kilobytes (4096 Bytes) of memory
	memory [4 * Kilobytes]uint8

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
}

func Frame() [32 * 64]bool {
	chip := ChipContext{}
	return testFrameBufferPrimer(chip.framebuffer)
}

func testFrameBufferPrimer(buffer [32 * 64]bool) ([32*64]bool) {
    for i := 0; i < 64 * 32; i++ {
        if i%7==0 {
            buffer[i] = true
        }
    }

    return buffer
}

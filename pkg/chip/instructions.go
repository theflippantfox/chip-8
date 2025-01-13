package chip

import "math/rand"

func (c Chip) instruction() {
	ins := (c.oc & 0xF000) >> 12
	switch ins {
	case 0x0:
		{
			m := c.oc & 0x000F
			if m == 0x0 { // If Operation Code (oc) is 0x00E0 then clear the screen
				c.display.clearDisplay()
				c.renderCycle = true
			} else if m == 0xE {
				// The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
				c.cpu.sp -= 1
				c.cpu.pc = c.cpu.stack[c.cpu.sp]
			}
		}
	case 0x1: // 1nnn - JP addr
		{
			// Jump to location nnn.
			// The interpreter sets the program counter to nnn.
			c.cpu.pc = c.oc & 0x0FFF
		}
	case 0x2: // 2nnn - CALL addr
		{
			// Call subroutine at nnn.
			// The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
			c.cpu.stack[c.cpu.sp] = c.cpu.pc
			c.cpu.sp += 1
			c.cpu.pc = c.oc & 0x0FFF

		}
	case 0x3: // 3xkk - SE Vx, byte
		{
			// Skip next instruction if Vx = kk.
			// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
			x := (c.oc & 0x0F00) >> 8
			if c.cpu.vx[x] == c.oc&0x00FF {
				c.cpu.increment_pc()
			}

		}
	case 0x4: // 4xkk - SNE Vx, byte
		{
			// Skip next instruction if Vx != kk.
			// The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
			x := (c.oc & 0x0F00) >> 8
			if c.cpu.vx[x] != c.oc&0x00FF {
				c.cpu.increment_pc()
			}

		}
	case 0x5: // 5xy0 - SE Vx, Vy
		{
			// Skip next instruction if Vx = Vy.
			// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
			x := (c.oc & 0x0F00) >> 8
			y := (c.oc & 0x00F0) >> 4

			if c.cpu.vx[x] == c.cpu.vx[y] {
				c.cpu.increment_pc()
			}

		}
	case 0x6: // 6xkk - LD Vx, byte
		{
			// Set Vx = kk.
			// The interpreter puts the value kk into register Vx.
			x := (c.oc & 0x0F00) >> 8
			c.cpu.vx[x] = (c.oc & 0x00FF)

		}
	case 0x7: // 7xkk - ADD Vx, byte
		{
			// Set Vx = Vx + kk.
			// Adds the value kk to the value of register Vx, then stores the result in Vx.
			x := (c.oc & 0x0F00) >> 8
			c.cpu.vx[x] += (c.oc & 0x00FF)
		}
	case 0x8: // 8xym
		{
			x := (c.oc & 0x0F00) >> 8
			y := (c.oc & 0x00F0) >> 4
			m := c.oc & 0x000F

			switch m {
			case 0x0: // 8xy0 - LD Vx, Vy
				{
					//Set Vx = Vy.
					//Stores the value of register Vy in register Vx.
					c.cpu.vx[x] = c.cpu.vx[y]
				}
			case 0x1: // 8xy1 - OR Vx, Vy
				{
					// Set Vx = Vx OR Vy.
					// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx.
					// A bitwise OR compares the corrseponding bits from two values, and if either bit is 1,
					// then the same bit in the result is also 1. Otherwise, it is 0.

					c.cpu.vx[x] |= c.cpu.vx[y]
				}
			case 0x2: // 8xy2 - AND Vx, Vy
				{
					// Set Vx = Vx AND Vy.
					// Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
					// A bitwise AND compares the corrseponding bits from two values, and if both bits are 1,
					// then the same bit in the result is also 1. Otherwise, it is 0.
					c.cpu.vx[x] &= c.cpu.vx[y]
				}
			case 0x3: // Set Vx = Vx XOR Vy.
				{
					// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx.
					// An exclusive OR compares the corrseponding bits from two values, and if the bits are not both the same,
					// then the corresponding bit in the result is set to 1. Otherwise, it is 0.
					c.cpu.vx[x] ^= c.cpu.vx[y]
				}
			case 0x4: // 8xy4 - ADD Vx, Vy
				{
					// Set Vx = Vx + Vy, set VF = carry.
					// The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1,
					// otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
					c.cpu.vx[x] += +c.cpu.vx[y]
					c.cpu.vx[0xF] = 0
					if c.cpu.vx[y] > 0xFF-c.cpu.vx[x] {
						c.cpu.vx[0xF] = 1
					}
				}
			case 0x5: // 8xy5 - SUB Vx, Vy
				{
					// Set Vx = Vx - Vy, set VF = NOT borrow.
					// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
					c.cpu.vx[0xF] = 0
					if c.cpu.vx[x] > c.cpu.vx[y] {
						c.cpu.vx[0xF] = 1
					}
					c.cpu.vx[x] -= c.cpu.vx[y]
				}
			case 0x6: // 8xy6 - SHR Vx {, Vy}
				{
					// Set Vx = Vx SHR 1.
					// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
					c.cpu.vx[0xF] = c.cpu.vx[x] & 0x1
					c.cpu.vx[x] >>= 1 // >> divides by 2
				}
			case 0x7: // 8xy7 - SUBN Vx, Vy
				{
					// Set Vx = Vy - Vx, set VF = NOT borrow.
					// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
					c.cpu.vx[0xF] = 0
					if c.cpu.vx[y] > c.cpu.vx[x] {
						c.cpu.vx[0xF] = 1
					}

					c.cpu.vx[x] = c.cpu.vx[y] - c.cpu.vx[x]
				}
			case 0xE: // 8xyE - SHL Vx {, Vy}
				{
					// Set Vx = Vx SHL 1.
					// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
					c.cpu.vx[0xF] = ((c.cpu.vx[x] & 0x0F00) >> 8) >> 7
					c.cpu.vx[x] <<= 1
				}
			} // Inner Switch end
		} // Case end
	case 0x9: // 9xy0 - SNE Vx, Vy
		{
			// Skip next instruction if Vx != Vy.
			// The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
			x := (c.oc & 0x0F00) >> 8
			y := (c.oc & 0x00F0) >> 4

			if c.cpu.vx[x] != c.cpu.vx[y] {
				c.cpu.increment_pc()
			}
		}
	case 0xA: // Annn - LD I, addr
		{
			// Set I = nnn.
			// The value of register I is set to nnn.
			c.cpu.i = c.oc & 0x0FFF
		}
	case 0xB: // Bnnn - JP V0, addr
		{
			// Jump to location nnn + V0.
			// The program counter is set to nnn plus the value of V0.
			c.cpu.pc = (c.oc & 0x0FFF) + c.cpu.vx[0]
		}
	case 0xC: // Cxkk - RND Vx, byte
		{
			// Set Vx = random byte AND kk.
			// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk.
			// The results are stored in Vx. See instruction 8xy2 for more information on AND.
			x := (c.oc & 0x0F00) >> 8
			kk := c.oc & 0x00FF

			c.cpu.vx[x] = uint16(rand.Intn(255)) + kk
		}
	case 0xD: // Dxyn - DRW Vx, Vy, nibble
		{
			x := c.cpu.vx[(c.oc&0x0F00)>>8]
			y := c.cpu.vx[(c.oc&0x00F0)>>4]
			n := c.oc & 0x000F

			var i uint16 = 0
			var j uint16 = 0
			c.cpu.vx[0xF] = 0
			for j = 0; j < n; j++ {
				pixel := c.mem.fetch(c.cpu.i + j) // Get the pixel from memory

				for i = 0; i < 8; i++ {
					// check if the current pixel will be drawn by ANDING it to 1 aka
					// check if the pixel is set to 1 (This will scan through the byte,
					// one bit at the time)
					if pixel&(0x80>>i) != 0 {
						// since the pixel will be drawn, check the destination location in
						// gfx for collision aka verify if that location is flipped on (== 1)
						if c.display.fetchPixel(x+i, y+j) == 1 {
							c.cpu.vx[0xF] = 1
						}
						c.display.XORPixel(x+i, y+j)
					}
				}
			}
			c.renderCycle = true
		}
	case 0xE: // Ex9E - SKP Vx
		{
			c.cpu.increment_pc()
			// mode := c.oc & 0x00FF
			// if mode == 0x9E {
			// 	// x := (c.oc & 0x0F00) >> 8
			// 	// TODO: Add keypress and value check
			// 	if c.isKeyDown && true { // Skip next instruction if key with the value of Vx is pressed.
			//
			// 		increment_pc(c)
			// 	}
			// } else if mode == 0xA1 {
			// 	if !c.isKeyDown && false {
			// 		increment_pc(c)
			// 	}
			// }

		}
	case 0xF:
		{
			x := (c.oc & 0x0F00) >> 8
			kk := c.oc & 0x00FF

			if kk == 0x07 {
				c.cpu.vx[x] = uint16(c.delayTimer)
			} else if kk == 0x0A {
				key_pressed := false

				if !key_pressed {
					return
				}
			} else if kk == 0x15 {
				c.delayTimer = uint8(c.cpu.vx[x])
			} else if kk == 0x18 {
				c.soundTimer = uint8(c.cpu.vx[x])
			} else if kk == 0x1E {
				c.cpu.i += uint16(c.cpu.vx[x])
			} else if kk == 0x29 {
				c.cpu.i = uint16(c.cpu.vx[x]) * 0x5
			} else if kk == 0x33 {
				c.mem[c.cpu.i] = uint8(c.cpu.vx[x] / 100)
				c.mem[c.cpu.i+1] = uint8((c.cpu.vx[x] / 10) % 10)
				c.mem[c.cpu.i+2] = uint8(c.cpu.vx[x] % 100)
			} else if kk == 0x55 {
				var i uint16 = 0
				for i = 0; i < x; i++ {
					c.mem[c.cpu.i+i] = uint8(c.cpu.vx[i])
				}
			} else if kk == 0x65 {
				var i uint16 = 0
				for i = 0; i < x; i++ {
					c.cpu.vx[i] = uint16(c.mem[c.cpu.i+i])
				}
			}
		}
	}
}

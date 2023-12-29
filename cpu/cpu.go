package cpu

import (
	mem "chip8/memory"
	"fmt"
	"math/rand"
	"time"
)

type Chip struct {
	oc uint16

	vx [16]uint16 // Chip-8 has 16 general purpose resisters referred as Vx where x is a hexadecimal digit from 0 to F

	stack [16]uint16 // Used to store the address that the interpreter shoud return to when finished with a subroutine

	pc uint16 // Stores currently executing address
	i  uint16 // Index register used to store address
	sp uint8  //points to the topmost level of the stack

	delay_timer uint8 // Used for delay timer. When non-zero it automatically decremented at a rate of 60Hz.
	sound_timer uint8 // Used for sound timer. When non-zero it automatically decremented at a rate of 60Hz.

	isKeyDown     bool
	isRenderCycle bool
}

func NewChip() *Chip {
	return &Chip{
		pc: 0x200,
	}
}

func Reset(c *Chip, m *mem.Mem) bool {
	c.pc = 0x200 // Programs should start at 0x200

	c.oc = 0
	c.sp = 0

	c.delay_timer = 0
	c.sound_timer = 0

	for i := 0; i < len(m.Memory); i++ {
		m.Memory[i] = 0
	}

	for i := 0; i < len(c.stack); i++ {
		c.stack[i] = 0
	}

	for i := 0; i < len(c.vx); i++ {
		c.vx[i] = 0
	}

	return true
}

func increment_pc(c *Chip) {
	c.pc += 2
}

func Execute(c *Chip, m *mem.Mem) {
	c.oc = mem.Fetch(m, c.pc)
	fmt.Println(c.oc, c.pc, c.sp, c.vx, c.stack)

	increment_pc(c)

	Instruction(c, m)

	time.Sleep(time.Second / 10)
	fmt.Println()
}

func Instruction(c *Chip, m *mem.Mem) {
	ins := (c.oc & 0xF000) >> 12
	switch ins {
	case 0x0:
		{
			m := c.oc & 0x000F
			fmt.Printf("0x0")
			if m == 0x0 { // If Operation Code (oc) is 0x00E0 then clear the screen
				// TODO: Clear Screen
			} else if m == 0xE {
				// The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
				c.pc = uint16(c.sp)
				c.sp -= 1
			}
		}
	case 0x1: // 1nnn - JP addr
		{
			// Jump to location nnn.
			// The interpreter sets the program counter to nnn.
			fmt.Printf("0x1")
			c.pc = c.oc & 0x0FFF
		}
	case 0x2: // 2nnn - CALL addr
		{
			// Call subroutine at nnn.
			// The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
			fmt.Printf("0x2")
			c.sp += 1
			c.stack[c.sp] = c.pc
			c.pc = c.oc & 0x0FFF

		}
	case 0x3: // 3xkk - SE Vx, byte
		{
			// Skip next instruction if Vx = kk.
			// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
			fmt.Printf("0x3")
			x := (c.oc & 0x0F00) >> 8
			if c.vx[x] == c.oc&0x00FF {
				increment_pc(c)
			}

		}
	case 0x4: // 4xkk - SNE Vx, byte
		{
			// Skip next instruction if Vx != kk.
			// The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
			fmt.Printf("0x4")
			x := (c.oc & 0x0F00) >> 8
			if c.vx[x] != c.oc&0x00FF {
				increment_pc(c)
			}

		}
	case 0x5: // 5xy0 - SE Vx, Vy
		{
			// Skip next instruction if Vx = Vy.
			// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
			fmt.Printf("0x5")
			x := (c.oc & 0x0F00) >> 8
			y := (c.oc & 0x00F0) >> 4

			if c.vx[x] == c.vx[y] {
				increment_pc(c)
			}

		}
	case 0x6: // 6xkk - LD Vx, byte
		{
			// Set Vx = kk.
			// The interpreter puts the value kk into register Vx.
			fmt.Printf("0x6")
			x := (c.oc & 0x0F00) >> 8
			c.vx[x] = (c.oc & 0x00FF)

		}
	case 0x7: // 7xkk - ADD Vx, byte
		{
			// Set Vx = Vx + kk.
			// Adds the value kk to the value of register Vx, then stores the result in Vx.
			fmt.Printf("0x7")
			x := (c.oc & 0x0F00) >> 8
			c.vx[x] += (c.oc & 0x00FF)
		}
	case 0x8: // 8xym
		{
			fmt.Printf("0x8")
			x := (c.oc & 0x0F00) >> 8
			y := (c.oc & 0x00F0) >> 4
			m := c.oc & 0x000F

			switch m {
			case 0x0: // 8xy0 - LD Vx, Vy
				{
					//Set Vx = Vy.
					//Stores the value of register Vy in register Vx.
					c.vx[x] = c.vx[y]
				}
			case 0x1: // 8xy1 - OR Vx, Vy
				{
					// Set Vx = Vx OR Vy.
					// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx.
					// A bitwise OR compares the corrseponding bits from two values, and if either bit is 1,
					// then the same bit in the result is also 1. Otherwise, it is 0.

					c.vx[x] |= c.vx[y]
				}
			case 0x2: // 8xy2 - AND Vx, Vy
				{
					// Set Vx = Vx AND Vy.
					// Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
					// A bitwise AND compares the corrseponding bits from two values, and if both bits are 1,
					// then the same bit in the result is also 1. Otherwise, it is 0.
					c.vx[x] &= c.vx[y]
				}
			case 0x3: // Set Vx = Vx XOR Vy.
				{
					// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx.
					// An exclusive OR compares the corrseponding bits from two values, and if the bits are not both the same,
					// then the corresponding bit in the result is set to 1. Otherwise, it is 0.
					c.vx[x] ^= c.vx[y]
				}
			case 0x4: // 8xy4 - ADD Vx, Vy
				{
					// Set Vx = Vx + Vy, set VF = carry.
					// The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1,
					// otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
					sum := c.vx[x] + c.vx[y]
					c.vx[0xF] = 0
					if sum > 0x0008 {
						c.vx[0xF] = 1
					}
					c.vx[x] = sum
				}
			case 0x5: // 8xy5 - SUB Vx, Vy
				{
					// Set Vx = Vx - Vy, set VF = NOT borrow.
					// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
					c.vx[0xF] = 0
					if c.vx[x] > c.vx[y] {
						c.vx[0xF] = 1
					}
					c.vx[x] -= c.vx[y]
				}
			case 0x6: // 8xy6 - SHR Vx {, Vy}
				{
					// Set Vx = Vx SHR 1.
					// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
					c.vx[0xF] = 0
					if c.vx[x]&0x000F == 0x1 {
						c.vx[0xF] = 1
					}

					c.vx[x] /= 2
				}
			case 0x7: // 8xy7 - SUBN Vx, Vy
				{
					// Set Vx = Vy - Vx, set VF = NOT borrow.
					// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
					c.vx[0xF] = 0
					if c.vx[y] > c.vx[x] {
						c.vx[0xF] = 1
					}

					c.vx[x] = c.vx[y] - c.vx[x]
				}
			case 0xE: // 8xyE - SHL Vx {, Vy}
				{
					// Set Vx = Vx SHL 1.
					// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
					c.vx[0xF] = 0
					if ((c.vx[x] & 0xF000) >> 12) == 1 {
						c.vx[0xF] = 1
					}

					c.vx[x] *= 2
				}
			} // Inner Switch end
		} // Case end
	case 0x9: // 9xy0 - SNE Vx, Vy
		{
			// Skip next instruction if Vx != Vy.
			// The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
			fmt.Printf("0x9")
			x := (c.oc & 0x0F00) >> 8
			y := (c.oc & 0x00F0) >> 4

			if c.vx[x] != c.vx[y] {
				increment_pc(c)
			}
		}
	case 0xA: // Annn - LD I, addr
		{
			// Set I = nnn.
			// The value of register I is set to nnn.
			fmt.Printf("0xA")
			c.i = c.oc & 0x0FFF
		}
	case 0xB: // Bnnn - JP V0, addr
		{
			// Jump to location nnn + V0.
			// The program counter is set to nnn plus the value of V0.
			fmt.Printf("0xB")
			c.pc = (c.oc & 0x0FFF) + c.vx[0]
		}
	case 0xC: // Cxkk - RND Vx, byte
		{
			// Set Vx = random byte AND kk.
			// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk.
			// The results are stored in Vx. See instruction 8xy2 for more information on AND.
			fmt.Printf("0xC")
			x := (c.oc & 0x0F00) >> 8
			kk := c.oc & 0x00FF

			c.vx[x] = uint16(rand.Intn(255)) + kk
		}
	case 0xD: // Dxyn - DRW Vx, Vy, nibble
		{
			fmt.Printf("0xD")

			// TODO: Add dxyn function
		}
	case 0xE: // Ex9E - SKP Vx
		{
			fmt.Printf("0xE")
            increment_pc(c)
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
			fmt.Printf("0xF")

			x := (c.oc & 0x0F00) >> 8
			kk := c.oc & 0x00FF

			if kk == 0x07 {
				c.vx[x] = uint16(c.delay_timer)
			} else if kk == 0x0A {
				key_pressed := false

				

				if !key_pressed {
					return
				}
			} else if kk == 0x15 {
				c.delay_timer = uint8(c.vx[x])
			} else if kk == 0x18 {
				c.sound_timer = uint8(c.vx[x])
			} else if kk == 0x1E {
				c.i += uint16(c.vx[x])
			} else if kk == 0x29 {
				c.i = uint16(c.vx[x]) * 0x5
			} else if kk == 0x33 {
				m.Memory[c.i] = uint8(c.vx[x] / 100)
				m.Memory[c.i+1] = uint8((c.vx[x] / 10) % 10) 
				m.Memory[c.i+2] = uint8(c.vx[x] % 100) 
			} else if kk == 0x55 {
				var i uint16 = 0
				for i = 0; i < x; i++ {
					m.Memory[c.i+i] = uint8(c.vx[i])
				}
			} else if kk == 0x65 {
				var i uint16 = 0
				for i = 0; i < x; i++ {
					c.vx[i] = uint16(m.Memory[c.i+i])
				}
			}
        }
	}
}

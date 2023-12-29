package cpu

import (
	"chip8/memory"
	"fmt"
	"time"
)

type Chip struct {
	oc uint16

	vx [16]uint16 // Chip-8 has 16 general purpose resisters referred as Vx where x is a hexadecimal digit from 0 to F

	stack [16]uint16 // Used to store the address that the interpreter shoud return to when finished with a subroutine

	pc uint16 // Stores currently executing address
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
	fmt.Println(c.oc, c.pc, c.vx)
    Instruction(c, m)
	time.Sleep(time.Second * 2)
}

func Instruction(c *Chip, m *mem.Mem) {
	ins := (c.oc & 0xF000) >> 12
	switch ins {
	case 0x0:
		{
			m := c.oc & 0x000F
			if m == 0x0 { // If Operation Code (oc) is 0x00E0 then clear the screen
				// TODO: Clear Screen
			} else if m == 0xE {
				// The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
				c.pc = uint16(c.sp)
				c.sp -= 1
			}
			increment_pc(c)
		}
	case 0x1:
		{
			c.pc = c.oc & 0x0FFF
			increment_pc(c)
		}
	case 0x2:
		{
		}
	case 0x3:
		{
		}
	case 0x4:
		{
		}
	case 0x5:
		{
		}
	case 0x6:
		{
		}
	case 0x7:
		{
		}
	case 0x8:
		{
		}
	case 0x9:
		{
		}
	case 0xA:
		{
		}
	case 0xB:
		{
		}
	case 0xC:
		{
		}
	case 0xD:
		{
		}
	case 0xE:
		{
		}
	case 0xF:
		{
		}

	}
}

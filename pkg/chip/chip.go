package chip

import (
	"os"
	"time"
)

type Chip struct {
	oc         uint16
	delayTimer uint8
	soundTimer uint8

	keyDown     bool
	renderCycle bool

	cpu     CPU
	mem     Memory
	display Display
}

func Init() {
	c := Chip{
		mem: Memory{},
		cpu: CPU{
			pc: 0x200,
		},
		display: Display{},
	}

	roamArg := os.Args[1]
	roam := ""
	if roamArg != "" {
		roam = roamArg
	}

	c.mem.loadROMtoMemory(roam)

	for i := 0; i == i; i++ {
		c.emulateCycle()
	}
}

func (c Chip) reset() bool {
	c.cpu.pc = 0x200 // Programs should start at 0x200

	c.oc = 0
	c.cpu.sp = 0

	c.delayTimer = 0
	c.soundTimer = 0

	for i := 0; i < len(c.mem); i++ {
		c.mem[i] = 0
	}

	for i := 0; i < len(c.cpu.stack); i++ {
		c.cpu.stack[i] = 0
	}

	for i := 0; i < len(c.cpu.vx); i++ {
		c.cpu.vx[i] = 0
	}

	return true
}

func (c Chip) emulateCycle() {
	c.oc = c.mem.fetch(c.cpu.pc)
	c.renderCycle = false
	c.cpu.increment_pc()

	c.instruction()

	if c.delayTimer > 0 {
		c.delayTimer -= 1
	}

	if c.soundTimer > 0 {
		c.soundTimer -= 1
	}

	time.Sleep(time.Second / 7000)

	if c.renderCycle {
		c.display.render()
	}
}

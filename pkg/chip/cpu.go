package chip

type CPU struct {
	stack [16]uint16
	vx    [16]uint16
	pc    uint16
	i     uint16
	sp    uint8
}

func (cpu CPU) increment_pc() {
	cpu.pc += 2
}

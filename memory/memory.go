package mem

import (
	"fmt"
	"io"
	"os"
)

type Mem struct {
	// The Chip-8 language is capable of accessing up to 4KB (4,096 bytes) of RAM,
	// from location 0x000 (0) to 0xFFF (4095). The first 512 bytes, from 0x000 to 0x1FF,
	// are where the original interpreter was located, and should not be used by programs.

	Memory [4096]uint8
}

func NewMemory() *Mem {
	return &Mem{}
}

func Reset(m *Mem) {
	for i := 0; i < len(m.Memory); i++ {
		m.Memory[i] = 0
	}
}

func LoadROMtoMemory(m *Mem, path string) {
	f, err := os.Open(path)
	if err != nil {
	}
	defer f.Close()
	buf := make([]uint8, 4096)
	stats, err := f.Stat()
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			var i int64 = 0
			for i = 0; i <= stats.Size() && i <= 4096-512; i++ {
				m.Memory[0x0200+i] = buf[i]
			}
		}
	}
}

func Fetch(m *Mem, pc uint16) uint16 {
	data := (uint16(m.Memory[pc]) << 8) | uint16(m.Memory[pc+1])
	return data
}

func Put(m *Mem, addr uint16, val uint16) {
	m.Memory[addr] = uint8(val & 0xFF)
	m.Memory[addr+1] = uint8(val >> 8)
}

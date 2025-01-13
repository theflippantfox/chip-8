package chip

import (
	"fmt"
	"io"
	"os"
)

type Memory [4096]uint8

func (m Memory) loadROMtoMemory(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
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
				m[0x0200+i] = buf[i]
			}
		}
	}
}

func (m Memory) fetch(pc uint16) uint16 {
	data := (uint16(m[pc]) << 8) | uint16(m[pc+1])
	return data
}

func (m Memory) put(addr uint16, val uint16) {
	m[addr] = uint8(val & 0xFF)
	m[addr+1] = uint8(val >> 8)
}

package chip 

type Display [64][32]uint8

func (d Display) clearDisplay() {}

func (d Display) render() {
	for i := 0; i < 64; i++ {
		for j := 0; j <= 32; j++ {
			print(d[i][j])
		}
	}
}

func (d Display) fetchPixel(x, y uint16) uint8 {
	return d[x][y]
}
func (d Display) XORPixel(x, y uint16) {
	d[x][y] ^= d[x][y]
}

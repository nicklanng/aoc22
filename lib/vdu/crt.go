package vdu

import "fmt"

type CRT struct {
	cpu    *CPU
	width  int
	height int

	scan   int
	pixels []byte
}

func NewCRT(cpu *CPU, width, height int) *CRT {
	return &CRT{
		cpu: cpu,
		width: width,
		height: height,
		pixels: make([]byte, width*height),
	}
}

func (crt *CRT) ExecuteCycle() {
	x := crt.cpu.XRegister()
	if x < (crt.scan%crt.width)-1 || x > (crt.scan%crt.width)+1{
		crt.pixels[crt.scan] = '.'
	} else {
		crt.pixels[crt.scan] = '#'
	}
	crt.scan = (crt.scan + 1) % len(crt.pixels)
}

func (crt *CRT) Render() {
	for y := 0; y < crt.height; y++ {
		for x := 0; x < crt.width; x++ {
			fmt.Print(string(crt.pixels[y*crt.width+x]) + " ")
		}
		fmt.Println()
	}
}


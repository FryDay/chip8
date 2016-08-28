package chip8

import (
	"fmt"
	"math/rand"
)

type Chip8 struct {
	opcode     uint16
	memory     [4096]byte
	v          [16]byte
	i          uint16
	pc         uint16
	draw       bool
	display    [64 * 32]byte
	delayTimer byte
	soundTimer byte
	stack      [16]uint16
	sp         uint16
	key        [16]byte
}

var font = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func (c *Chip8) Initialize() {
	c.pc = 0x200
	c.opcode = 0
	c.i = 0
	c.sp = 0

	c.draw = false
	for i := 0; i < 16; i++ {
		c.stack[i] = 0
		c.v[i] = 0
	}
	for i := range c.memory {
		c.memory[i] = 0
	}

	for i := 0; i < 80; i++ {
		c.memory[i] = font[i]
	}
}

func (c *Chip8) Cycle() {
	c.opcode = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	a := uint16(c.opcode & 0xf000)
	b := uint16(c.opcode & 0x0fff)
	xReg := uint16(c.opcode & 0x0f00 >> 8)
	yReg := uint16(c.opcode & 0x00f0 >> 4)

	// fmt.Printf("0x%X\n", c.opcode)

	switch a {
	case 0x0000:
		switch c.opcode {
		case opSystem:
			panic(fmt.Sprintln("The internet tells me this shouldn't happen..."))
		case opClear:
			for i := range c.display {
				c.display[i] = 0
			}
			c.draw = true
		case opReturn:
			c.sp--
			c.pc = c.stack[c.sp]
			c.pc += 2
		}
	case opJump:
		c.pc = b
	case opCall:
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = b
	case opSkipIfEqual:
		if c.v[xReg] == byte(c.opcode&0xff00) {
			c.pc += 2
		}
		c.pc += 2
	case opSkipIfNotEqual:
		if c.v[xReg] != byte(c.opcode&0xff00) {
			c.pc += 2
		}
		c.pc += 2
	case opSkipIfEqualRegister:
		if c.v[xReg] == c.v[yReg] {
			c.pc += 2
		}
	case opSetValue:
		c.v[xReg] = c.v[yReg]
		c.pc += 2
	case opAddValue:
		c.v[xReg] = c.v[xReg] & +byte(c.opcode&0xff00)
		c.pc += 2
	case 0x8000:
		switch c.opcode & 0x000f {
		case 0x0:
			c.v[xReg] = c.v[yReg]
			c.pc += 2
		case 0x1:
			c.v[xReg] |= c.v[yReg]
			c.pc += 2
		case 0x2:
			c.v[xReg] &= c.v[yReg]
			c.pc += 2
		case 0x3:
			c.v[xReg] ^= c.v[yReg]
			c.pc += 2
		case 0x4:
			if c.v[(c.opcode&0x00f0)>>4] > (0xff - c.v[(c.opcode&0x0f00)>>8]) {
				c.v[0xf] = 1
			} else {
				c.v[0xf] = 0
			}
			c.v[(c.opcode&0x0f00)>>8] += c.v[(c.opcode&0x00f0)>>4]
			c.pc += 2
		case 0x5:
			if c.v[xReg] < c.v[yReg] {
				c.v[0xf] = 0
			} else {
				c.v[0xf] = 1
			}
			c.v[xReg] = c.v[xReg] & -c.v[yReg]
			c.pc += 2
		case 0x6:
			c.v[0xf] = c.v[xReg] & 1
			c.v[xReg] >>= 1
			c.pc += 2
		case 0x7:
			if c.v[yReg] < c.v[xReg] {
				c.v[0xf] = 0
			} else {
				c.v[0xf] = 1
			}
			c.v[xReg] = c.v[yReg] & -c.v[xReg]
			c.pc += 2
		case 0xe:
			c.v[0xf] = (c.v[xReg] & 0x08) >> 7
			c.v[xReg] <<= 1
			c.pc += 2
		}
	case opSkipIfNotEqualRegister:
		if c.v[xReg] != c.v[yReg] {
			c.pc += 2
		}
		c.pc += 2
	case opSetIndex:
		c.i = c.opcode & 0x0fff
		c.pc += 2
	case opJumpRelative:
		c.pc = b + uint16(c.v[0])
	case opAndRandom:
		c.v[xReg] = byte(rand.Int()) & byte(c.opcode&0xff00)
		c.pc += 2
	case opDraw:
		x := uint16(c.v[(c.opcode&0x0f00)>>8])
		y := uint16(c.v[(c.opcode&0x00f0)>>4])
		height := uint16(c.opcode & 0x000f)
		var pixel uint16
		c.v[0xf] = 0
		for yLine := uint16(0); yLine < height; yLine++ {
			pixel = uint16(c.memory[c.i+yLine])
			for xLine := uint16(0); xLine < 8; xLine++ {
				if (pixel & (0x80 >> xLine)) != 0 {
					if c.display[x+xLine+((y+yLine)*64)] == 1 {
						c.v[0xf] = 1
						c.display[x+xLine+((y+yLine)*64)] ^= 1
					}
				}
			}
		}
		c.draw = true
		c.pc += 2
	case 0xe000:
		switch c.opcode & 0x00ff {
		case 0x9e:
			if c.key[c.v[(c.opcode&0x0f00)>>8]] != 0 {
				c.pc += 2
			}
			c.pc += 2
		case 0xa1:
			if c.key[c.v[(c.opcode&0x0f00)>>8]] == 0 {
				c.pc += 2
			}
			c.pc += 2
		}
	case 0xf000:
		last1 := c.opcode & 0x000f
		last2 := c.opcode & 0x00ff
		switch last1 {
		case 0x7:
			c.v[xReg] = c.delayTimer
			c.pc += 2
		case 0xa:
			//TODO Implement
		}

		switch last2 {
		case 0x15:
			c.delayTimer = c.v[xReg]
		case 0x18:
			c.soundTimer = c.v[xReg]
		case 0x1e:
			if (uint16(c.v[xReg]) + c.i) > uint16(0xfff) {
				c.v[0xf] = 1
			} else {
				c.v[0xf] = 0
			}
			c.i += uint16(c.v[xReg])
			c.pc += 2
		case 0x29:
			c.i = uint16(c.v[xReg] * 5)
			c.pc += 2
		case 0x33:
			c.memory[c.i] = c.v[(c.opcode&0x0f00)>>8] / 100
			c.memory[c.i+1] = (c.v[(c.opcode&0x0f00)>>8] / 10) % 10
			c.memory[c.i+2] = (c.v[(c.opcode&0x0f00)>>8] % 100) % 10
			c.pc += 2
		case 0x55:
			for i := uint16(0); i < xReg; i++ {
				c.memory[c.i+i] = c.v[i]
			}
			c.pc += 2
		case 0x65:
			for i := uint16(0); i < xReg; i++ {
				c.v[i] = c.memory[c.i+i]
			}
			c.pc += 2
		}
	default:
		panic(fmt.Sprintf("Unknown opcode: 0x%X", c.opcode))
	}

	if c.delayTimer > 0 {
		c.delayTimer--
	}
	if c.soundTimer > 0 {
		fmt.Println("BEEP")
		c.soundTimer--
	}
}

func (c *Chip8) LoadROM(r []byte) {
	for i := range r {
		c.memory[i+512] = r[i]
	}
}

func (c *Chip8) Draw() bool {
	return c.draw
}

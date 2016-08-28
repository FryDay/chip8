package chip8

import (
	"fmt"
	"math/rand"
	"time"
)

type Chip8 struct {
	opcode     uint16
	memory     [4096]byte
	v          [16]byte
	index      uint16
	pc         uint16
	Draw       bool
	Display    [64 * 32]byte
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
	c.index = 0
	c.sp = 0

	c.Draw = false
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

	switch a {
	case 0x0000:
		switch c.opcode {
		case 0x0000: // Unused
			panic(fmt.Sprintln("The internet tells me this shouldn't happen..."))
		case 0x00e0: // Clear
			for i := range c.Display {
				c.Display[i] = 0
			}
			c.Draw = true
			c.pc += 2
		case 0x00ee: // Return
			c.sp--
			c.pc = c.stack[c.sp]
			c.pc += 2
		}
	case 0x1000: // Jump
		c.pc = b
	case 0x2000: // Call
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = b
	case 0x3000: // Skip if Equal
		if c.v[xReg] == byte(c.opcode&0x00ff) {
			c.pc += 2
		}
		c.pc += 2
	case 0x4000: // Skip if Not Equal
		if c.v[xReg] != byte(c.opcode&0x00ff) {
			c.pc += 2
		}
		c.pc += 2
	case 0x5000: // Skip if x Register Equals y Register
		if c.v[xReg] == c.v[yReg] {
			c.pc += 2
		}
		c.pc += 2
	case 0x6000: // Set Value
		c.v[xReg] = byte(c.opcode & 0x00ff)
		c.pc += 2
	case 0x7000: // Add Value
		c.v[xReg] += byte(c.opcode & 0x00ff)
		c.pc += 2
	case 0x8000:
		switch c.opcode & 0x000f {
		case 0x0: // Set x Register to y Register
			c.v[xReg] = c.v[yReg]
			c.pc += 2
		case 0x1: // Or
			c.v[xReg] |= c.v[yReg]
			c.pc += 2
		case 0x2: // And
			c.v[xReg] &= c.v[yReg]
			c.pc += 2
		case 0x3: // Xor
			c.v[xReg] ^= c.v[yReg]
			c.pc += 2
		case 0x4: // Add y Register to x Register
			if c.v[yReg] > (0xff - c.v[xReg]) {
				c.v[0xf] = 1
			} else {
				c.v[0xf] = 0
			}
			c.v[xReg] += c.v[yReg]
			c.pc += 2
		case 0x5: // Subtract y Register from x Register
			if c.v[yReg] > c.v[xReg] {
				c.v[0xf] = 0
			} else {
				c.v[0xf] = 1
			}
			c.v[xReg] -= c.v[yReg]
			c.pc += 2
		case 0x6: // Shift x Register Right by 1
			c.v[0xf] = c.v[xReg] & 1
			c.v[xReg] >>= 1
			c.pc += 2
		case 0x7: // Sets x Register to y Register Minus x Register
			if c.v[yReg] < c.v[xReg] {
				c.v[0xf] = 0
			} else {
				c.v[0xf] = 1
			}
			c.v[xReg] = c.v[yReg] - c.v[xReg]
			c.pc += 2
		case 0xe: // Shift x Register Left by 1
			c.v[0xf] = c.v[xReg] >> 7
			c.v[xReg] <<= 1
			c.pc += 2
		}
	case 0x9000: // Skip if x Register Not Equal y Register
		if c.v[xReg] != c.v[yReg] {
			c.pc += 2
		}
		c.pc += 2
	case 0xa000: // Sets i to NNN
		c.index = c.opcode & 0x0fff
		c.pc += 2
	case 0xb000: // Jumps to Address NNN Plus v0
		c.pc = b + uint16(c.v[0])
	case 0xc000: // Sets x Register to Bitwise And on Random and NNN
		c.v[xReg] = byte(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(255)) & byte(c.opcode&0x00ff)
		c.pc += 2
	case 0xd000: // Draw
		x := uint16(c.v[(c.opcode&0x0f00)>>8])
		y := uint16(c.v[(c.opcode&0x00f0)>>4])
		height := uint16(c.opcode & 0x000f)
		var pixel uint16
		c.v[0xf] = 0
		for yLine := uint16(0); yLine < height; yLine++ {
			pixel = uint16(c.memory[c.index+yLine])
			for xLine := uint16(0); xLine < 8; xLine++ {
				if (pixel & (0x80 >> xLine)) != 0 {
					if c.Display[x+xLine+((y+yLine)*64)] == 1 {
						c.v[0xf] = 1
					}
					c.Display[x+xLine+((y+yLine)*64)] ^= 1
				}
			}
		}
		c.Draw = true
		c.pc += 2
	case 0xe000:
		switch c.opcode & 0x00ff {
		case 0x9e: // Skip If Key Pressed
			// if c.key[c.v[xReg]] != 0 {
			// 	c.pc += 2
			// }
			//TODO Implement
			c.pc += 2
		case 0xa1: // Skip If Key Not Pressed
			// if c.key[c.v[(c.opcode&0x0f00)>>8]] == 0 {
			// 	c.pc += 2
			// }
			//TODO Implement
			c.pc += 4 // Will be 2
		}
	case 0xf000:
		switch c.opcode & 0x00ff {
		case 0x07: // Store Delay Timer
			c.v[xReg] = c.delayTimer
			c.pc += 2
		case 0x0a: // Await Key Press
			//TODO Implement
			c.pc += 2
		case 0x15: // Set Delay Timer
			c.delayTimer = c.v[xReg]
			c.pc += 2
		case 0x18: // Set Sound Timer
			c.soundTimer = c.v[xReg]
			c.pc += 2
		case 0x1e: // Add Index
			if (uint16(c.v[xReg]) + c.index) > uint16(0xfff) {
				c.v[0xf] = 1
			} else {
				c.v[0xf] = 0
			}
			c.index += uint16(c.v[xReg])
			c.pc += 2
		case 0x29: // Set Index Font Character
			c.index = uint16(c.v[xReg] * 0x5)
			c.pc += 2
		case 0x33: // Store BCD
			c.memory[c.index] = c.v[xReg] / 100
			c.memory[c.index+1] = (c.v[xReg] / 10) % 10
			c.memory[c.index+2] = (c.v[xReg] % 100) % 10
			c.pc += 2
		case 0x55: // Write Memory
			for i := uint16(0); i <= xReg; i++ {
				c.memory[c.index+i] = c.v[i]
			}

			c.pc += 2
		case 0x65: // Read Memory
			for i := uint16(0); i <= xReg; i++ {
				c.v[i] = c.memory[c.index+i]
			}
			c.index += uint16(c.v[xReg]) + 1
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

func (c *Chip8) LoadROM(rom []byte) {
	progSpace := c.memory[512:]
	copy(progSpace, rom)
}

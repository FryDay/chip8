package chip8

import (
	"fmt"
	"math/rand"
)

type Chip8 struct {
	Opcode     uint16
	Memory     [4096]byte
	V          [16]byte
	I          uint16
	PC         uint16
	Draw       bool
	Display    [64 * 32]byte
	DelayTimer byte
	SoundTimer byte
	Stack      [16]uint16
	SP         uint16
	Key        [16]byte
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
	c.PC = 0x200
	c.Opcode = 0
	c.I = 0
	c.SP = 0

	c.Draw = false
	for i := 0; i < 16; i++ {
		c.Stack[i] = 0
		c.V[i] = 0
	}
	for i := range c.Memory {
		c.Memory[i] = 0
	}

	for i := 0; i < 80; i++ {
		c.Memory[i] = font[i]
	}
}

func (c *Chip8) Cycle() {
	c.Opcode = uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
	a := uint16(c.Opcode & 0xf000)
	b := uint16(c.Opcode & 0x0fff)
	xReg := uint16(c.Opcode & 0xf0ff)
	yReg := uint16(c.Opcode & 0xff0f)

	switch a {
	case System:
		panic(fmt.Sprintln("The internet tells me this shouldn't happen..."))

	case Clear:
		for i := range c.Display {
			c.Display[i] = 0
		}
		c.Draw = true

	case Return:
		c.SP--
		c.PC = c.Stack[c.SP]
		c.PC += 2

	case Jump:
		c.PC = b

	case Call:
		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = b

	case SkipIfEqual:
		if c.V[xReg] == byte(c.Opcode&0xff00) {
			c.PC += 2
		}
		c.PC += 2

	case SkipIfNotEqual:
		if c.V[xReg] != byte(c.Opcode&0xff00) {
			c.PC += 2
		}
		c.PC += 2

	case SkipIfEqualRegister:
		if c.V[xReg] == c.V[yReg] {
			c.PC += 2
		}

	case SetValue:
		c.V[xReg] = c.V[yReg]
		c.PC += 2

	case AddValue:
		c.V[xReg] = c.V[xReg] & +byte(c.Opcode&0xff00)
		c.PC += 2

	case SetRegister:
		c.V[xReg] = c.V[yReg]
		c.PC += 2

	case Or:
		c.V[xReg] |= c.V[yReg]
		c.PC += 2

	case And:
		c.V[xReg] &= c.V[yReg]
		c.PC += 2

	case Xor:
		c.V[xReg] ^= c.V[yReg]
		c.PC += 2

	case AddRegister:
		if c.V[(c.Opcode&0x00f0)>>4] > (0xff - c.V[(c.Opcode&0x0f00)>>8]) {
			c.V[0xf] = 1
		} else {
			c.V[0xf] = 0
		}
		c.V[(c.Opcode&0x0f00)>>8] += c.V[(c.Opcode&0x00f0)>>4]
		c.PC += 2

	case SubtractYFromX:
		if c.V[xReg] < c.V[yReg] {
			c.V[0xf] = 0
		} else {
			c.V[0xf] = 1
		}
		c.V[xReg] = c.V[xReg] & -c.V[yReg]
		c.PC += 2

	case ShiftRight:
		c.V[0xf] = c.V[xReg] & 1
		c.V[xReg] >>= 1
		c.PC += 2

	case SubtractXFromY:
		if c.V[yReg] < c.V[xReg] {
			c.V[0xf] = 0
		} else {
			c.V[0xf] = 1
		}
		c.V[xReg] = c.V[yReg] & -c.V[xReg]
		c.PC += 2

	case ShiftLeft:
		c.V[0xf] = (c.V[xReg] & 0x08) >> 7
		c.V[xReg] <<= 1
		c.PC += 2

	case SkipIfNotEqualRegister:
		if c.V[xReg] != c.V[yReg] {
			c.PC += 2
		}
		c.PC += 2

	case SetIndex:
		c.I = c.Opcode & 0x0fff
		c.PC += 2

	case JumpRelative:
		c.PC = b + uint16(c.V[0])

	case AndRandom:
		c.V[xReg] = byte(rand.Int()) & byte(c.Opcode&0xff00)
		c.PC += 2

	case Draw:
		x := uint16(c.V[(c.Opcode&0x0f00)>>8])
		y := uint16(c.V[(c.Opcode&0x00f0)>>4])
		height := uint16(c.Opcode & 0x000f)
		var pixel uint16
		c.V[0xf] = 0
		for yLine := uint16(0); yLine < height; yLine++ {
			pixel = uint16(c.Memory[c.I+yLine])
			for xLine := uint16(0); xLine < 8; xLine++ {
				if (pixel & (0x80 >> xLine)) != 0 {
					if c.Display[x+xLine+((y+yLine)*64)] == 1 {
						c.V[0xf] = 1
						c.Display[x+xLine+((y+yLine)*64)] ^= 1
					}
				}
			}
		}
		c.Draw = true
		c.PC += 2

	case SkipIfKeyPressed:
		if c.Key[c.V[(c.Opcode&0x0f00)>>8]] != 0 {
			c.PC += 2
		}
		c.PC += 2

	case SkipIfKeyNotPressed:
		if c.Key[c.V[(c.Opcode&0x0f00)>>8]] == 0 {
			c.PC += 2
		}
		c.PC += 2

	case StoreDelayTimer:
		c.V[xReg] = c.DelayTimer
		c.PC += 2

	case AwaitKeyPress:
		//TODO Implement

	case SetDelayTimer:
		c.DelayTimer = c.V[xReg]

	case SetSoundTimer:
		c.SoundTimer = c.V[xReg]

	case AddIndex:
		if (uint16(c.V[xReg]) + c.I) > uint16(0xfff) {
			c.V[0xf] = 1
		} else {
			c.V[0xf] = 0
		}
		c.I += uint16(c.V[xReg])
		c.PC += 2

	case SetIndexFontCharacter:
		c.I = uint16(c.V[xReg] * 5)
		c.PC += 2

	case StoreBCD:
		c.Memory[c.I] = c.V[(c.Opcode&0x0f00)>>8] / 100
		c.Memory[c.I+1] = (c.V[(c.Opcode&0x0f00)>>8] / 10) % 10
		c.Memory[c.I+2] = (c.V[(c.Opcode&0x0f00)>>8] % 100) % 10
		c.PC += 2

	case WriteMemory:
		for i := uint16(0); i < xReg; i++ {
			c.Memory[c.I+i] = c.V[i]
		}
		c.PC += 2

	case ReadMemory:
		for i := uint16(0); i < xReg; i++ {
			c.V[i] = c.Memory[c.I+i]
		}
		c.PC += 2

	default:
		panic(fmt.Sprintf("Unknown opcode: 0x%X", c.Opcode))
	}

	if c.DelayTimer > 0 {
		c.DelayTimer--
	}
	if c.SoundTimer > 0 {
		fmt.Println("BEEP")
		c.SoundTimer--
	}
}

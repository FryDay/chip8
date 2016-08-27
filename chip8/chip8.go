package chip8

import "fmt"

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
	//clear Stack
	//clear registers V0-VF
	//clear memory

	// Load fontset
	for i := 0; i < 80; i++ {
		c.Memory[i] = font[i]
	}
}

func (c *Chip8) Cycle() {
	c.Opcode = uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
	a := uint16(c.Opcode & 0xf000)
	//b := byte(c.Opcode & 0x0fff)

	switch a {
	case System:
		panic(fmt.Sprintln("The internet tells me this shouldn't happen..."))

	case Clear:
		for i := range c.Display {
			c.Display[i] = 0
		}
		c.Draw = true

	case Return:
	case Jump:
	case Call:
		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = c.Opcode & 0x0fff

	case SkipIfEqual:
	case SkipIfNotEqual:
	case SkipIfEqualRegister:
	case SetValue:
	case AddValue:
	case SetRegister:
	case Or:
	case And:
	case Xor:
	case AddRegister:
		if c.V[(c.Opcode&0x00f0)>>4] > (0xff - c.V[(c.Opcode&0x0f00)>>8]) {
			c.V[0xf] = 1
		} else {
			c.V[0xf] = 0
		}
		c.V[(c.Opcode&0x0f00)>>8] += c.V[(c.Opcode&0x00f0)>>4]
		c.PC += 2

	case SubtractYFromX:
	case ShiftRight:
	case SubtractXFromY:
	case ShiftLeft:
	case SkipIfNotEqualRegister:
	case SetIndex:
		c.I = c.Opcode & 0x0fff
		c.PC += 2

	case JumpRelative:
	case AndRandom:
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
			c.PC += 4
		} else {
			c.PC += 2
		}

	case SkipIfKeyNotPressed:
		if c.Key[c.V[(c.Opcode&0x0f00)>>8]] != 0 {
			c.PC += 2
		} else {
			c.PC += 4
		}

	case StoreDelayTimer:
	case AwaitKeyPress:
	case SetDelayTimer:
	case SetSoundTimer:
	case AddIndex:
	case SetIndexFontCharacter:
	case StoreBCD:
		c.Memory[c.I] = c.V[(c.Opcode&0x0f00)>>8] / 100
		c.Memory[c.I+1] = (c.V[(c.Opcode&0x0f00)>>8] / 10) % 10
		c.Memory[c.I+2] = (c.V[(c.Opcode&0x0f00)>>8] % 100) % 10
		c.PC += 2

	case WriteMemory:
	case ReadMemory:
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

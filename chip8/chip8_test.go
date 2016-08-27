package chip8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Cycle_fetch_opcode(t *testing.T) {
	chip8 := Chip8{}
	chip8.Initialize()

	chip8.Memory[chip8.PC] = 0xA2
	chip8.Memory[chip8.PC+1] = 0xF0
	assert.Equal(t, chip8.Opcode, uint16(0x0))
	chip8.Cycle()
	assert.Equal(t, chip8.Opcode, uint16(0xa2f0))
}

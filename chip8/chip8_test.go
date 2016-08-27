package chip8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Cycle_fetch_opcode(t *testing.T) {
	chip8 := Chip8{}
	chip8.Initialize()

	chip8.memory[chip8.pc] = 0xA2
	chip8.memory[chip8.pc+1] = 0xF0
	assert.Equal(t, chip8.opcode, uint16(0x0))
	chip8.Cycle()
	assert.Equal(t, chip8.opcode, uint16(0xa2f0))
}

package main

import (
	"fmt"

	"github.com/FryDay/chip8/chip8"
)

func main() {
	chip8 := chip8.Chip8{}
	//graphics
	//input

	chip8.Initialize()
	//loadRom

	chip8.Memory[chip8.PC] = 0xA2
	chip8.Memory[chip8.PC+1] = 0xF0
	fmt.Printf("0x%X\n", chip8.Opcode)
	chip8.Cycle()
	fmt.Printf("0x%X\n", chip8.Opcode)

	// for {
	// 	chip8.Cycle()
	//if draw flag then draw
	//set keys
	// }
}

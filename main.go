package main

import (
	"io/ioutil"

	"github.com/FryDay/chip8/chip8"
)

func main() {
	chip8 := chip8.Chip8{}
	//graphics
	//input

	chip8.Initialize()
	rom, _ := ioutil.ReadFile("./roms/TETRIS")
	chip8.LoadROM(rom)

	for {
		chip8.Cycle()
		// if draw flag then draw
		// set keys
	}
}

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
	rom, _ := ioutil.ReadFile("./roms/BLINKY")
	chip8.LoadROM(rom)

	//fmt.Println(chip8.Memory)

	for {
		chip8.Cycle()
		// if draw flag then draw
		// set keys
	}
}

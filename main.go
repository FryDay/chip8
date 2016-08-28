package main

import (
	"io/ioutil"
	"runtime"

	"github.com/FryDay/chip8/chip8"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	chip8 := chip8.Chip8{}
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 320, "CHIP-8", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	//input

	chip8.Initialize()
	rom, _ := ioutil.ReadFile("./roms/TETRIS")
	chip8.LoadROM(rom)

	for !window.ShouldClose() {
		chip8.Cycle()
		if chip8.Draw() {
			//draw stuff
		}
		window.SwapBuffers()
		glfw.PollEvents()
		// set keys
	}
}

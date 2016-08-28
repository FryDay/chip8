package main

import (
	"io/ioutil"
	"runtime"

	"github.com/FryDay/chip8/chip8"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	screenWidth          = 640
	screenHeight         = 320
	zoom         float32 = 10
)

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(screenWidth, screenHeight, "CHIP-8", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0.0, screenWidth, screenHeight, 0.0, 1.0, -1.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	chip8 := chip8.Chip8{}

	//input

	chip8.Initialize()
	rom, _ := ioutil.ReadFile("./roms/TETRIS")
	chip8.LoadROM(rom)

	for !window.ShouldClose() {
		chip8.Cycle()
		if chip8.Draw {
			render(chip8.Display[:])
			chip8.Draw = false
		}

		window.SwapBuffers()
		glfw.PollEvents()
		// set keys
	}
}

func render(d []byte) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Begin(gl.QUADS)

	var row float32
	var col float32
	for i := range d {
		if col > 31 {
			row++
			col = 0
		}
		if d[i] == 1 {
			gl.Vertex2f(col*zoom, row*zoom)
			gl.Vertex2f(col*zoom, row*zoom+zoom)
			gl.Vertex2f(col*zoom+zoom, row*zoom+zoom)
			gl.Vertex2f(col*zoom+zoom, row*zoom)
		}
		col++
	}

	gl.End()
}

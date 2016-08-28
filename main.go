package main

import (
	"io/ioutil"
	"runtime"

	"github.com/FryDay/chip8/chip8"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	screenWidth  = 1280
	screenHeight = 640
)

var zoom float32

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

	zoom = screenWidth / 64

	chip8 := chip8.Chip8{}

	//input

	chip8.Initialize()
	rom, _ := ioutil.ReadFile("./roms/INVADERS")
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
	var row float32
	var col float32
	var colZoom float32
	var rowZoom float32

	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Begin(gl.QUADS)
	for i := range d {
		if col > 63 {
			row++
			col = 0
		}
		if d[i] == 1 {
			colZoom = col * zoom
			rowZoom = row * zoom
			gl.Vertex2f(colZoom, rowZoom)
			gl.Vertex2f(colZoom, rowZoom+zoom)
			gl.Vertex2f(colZoom+zoom, rowZoom+zoom)
			gl.Vertex2f(colZoom+zoom, rowZoom)
		}
		col++
	}
	gl.End()
}

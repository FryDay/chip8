package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/FryDay/chip8/chip8"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	screenWidth  float32 = 1280
	screenHeight float32 = 640
	zoom         float32
	romPtr       = flag.String("rom", "", "ROM to load (Required)")
	widthPtr     = flag.Int("width", 1280, "Screen width (Optional)")
)

func init() {
	flag.StringVar(romPtr, "r", "", "Same as -rom")
	flag.IntVar(widthPtr, "w", 1280, "Same as -width")
	runtime.LockOSThread()
}

func main() {
	flag.Parse()

	if *romPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if _, err := os.Stat("./roms/" + strings.ToUpper(*romPtr)); os.IsNotExist(err) {
		fmt.Println(*romPtr, "is not a valid ROM")
		os.Exit(1)
	}
	if float32(*widthPtr) != screenWidth {
		screenWidth = float32(*widthPtr)
		screenHeight = screenWidth / 2
	}

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(int(screenWidth), int(screenHeight), "CHIP-8", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	w, h := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(w), int32(h))
	gl.ClearColor(0, 0, 0, 1)
	window.SetKeyCallback(onKey)

	zoom = screenWidth / 64

	chip8 := chip8.Chip8{}

	//input

	chip8.Initialize()
	rom, _ := ioutil.ReadFile("./roms/" + strings.ToUpper(*romPtr))
	chip8.LoadROM(rom)

	for !window.ShouldClose() {
		glfw.PollEvents()

		chip8.Cycle()
		if chip8.Draw {
			render(chip8.Display[:])
			window.SwapBuffers()
			chip8.Draw = false
		}
		// set keys
	}
}

func render(d []byte) {
	var row float32
	var col float32
	// var colZoom float32
	// var rowZoom float32

	gl.Clear(gl.COLOR_BUFFER_BIT)

	for i := range d {
		if col > 63 {
			row++
			col = 0
		}
		if d[i] == 1 {
			// colZoom = col * zoom
			// rowZoom = row * zoom
			col++
		}
	}
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mod glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

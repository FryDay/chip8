package graphics

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	vao, vbo      uint32
	shaderProgram uint32
	verts         = []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
)

// Initialize ...
func Initialize(width, height int) {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Viewport(0, 0, int32(width), int32(height))
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	shaderProgram = newProgram(vertexShaderSource, fragmentShaderSource)
	gl.UseProgram(shaderProgram)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

// Render ...
func Render(d []byte) {
	// var row float32
	// var col float32
	// var colZoom float32
	// var rowZoom float32

	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(shaderProgram)
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)

	// for i := range d {
	// 	if col > 63 {
	// 		row++
	// 		col = 0
	// 	}
	// 	if d[i] == 1 {
	// 		// colZoom = col * zoom
	// 		// rowZoom = row * zoom
	// 		col++
	// 	}
	// }
}

// Cleanup ...
func Cleanup() {
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
}

func newProgram(vSource, fSource string) uint32 {
	vertexShader := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	defer gl.DeleteShader(vertexShader)
	defer gl.DeleteShader(fragmentShader)

	shaderProgram := gl.CreateProgram()

	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	var status int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		panic(fmt.Errorf("Failed to link program: %v", log))
	}

	return shaderProgram
}

func compileShader(shaderSource string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)

	source, free := gl.Strs(shaderSource)
	gl.ShaderSource(shader, 1, source, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		panic(fmt.Errorf("Failed to compile %v: %v", shaderSource, log))
	}
	return shader
}

package main

import (
	"C"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/exp/f32"
)

// #include <stdlib.h>
import "C"

const vertexSource = `//#version 120 // OpenGL 2.1.
//#version 100 // WebGL.

attribute vec3 aVertexPosition;

uniform mat4 uMVMatrix;
uniform mat4 uPMatrix;

void main() {
	gl_Position = uPMatrix * uMVMatrix * vec4(aVertexPosition, 1.0);
}
`

const fragmentSource = `//#version 120 // OpenGL 2.1.
//#version 100 // WebGL.

void main() {
	gl_FragColor = vec4(1.0, 1.0, 1.0, 1.0);
}
`

// func goxjsTest() {
// 	// Initialize GL
// 	fatal(glfw.Init(gl.ContextWatcher))
// 	defer glfw.Terminate()
//
// 	// Create Window Object
// 	var windowSize = [2]int{640, 480}
// 	glfw.WindowHint(glfw.Samples, 8)
// 	window, err := glfw.CreateWindow(windowSize[0], windowSize[1], "", nil, nil)
// 	fatal(err)
//
// 	window.MakeContextCurrent()
//
// 	// Print GL Configuration Information
// 	dumpGLConfig()
//
// 	// Logic for Moving Mouse Cursor
// 	cursorPos := [2]float32{200, 200}
// 	cursorPosCallback := func(_ *glfw.Window, x, y float64) {
// 		cursorPos[0], cursorPos[1] = float32(x), float32(y)
// 	}
// 	window.SetCursorPosCallback(cursorPosCallback)
//
// 	// Callback for when Framebuffer Changes;
// 	// Framebuffers are the objects that contain image and rendering data.
// 	// Assumedly this would be triggered on window resize
// 	framebufferSizeCallback := func(w *glfw.Window, framebufferSize0, framebufferSize1 int) {
// 		gl.Viewport(0, 0, framebufferSize0, framebufferSize1)
// 		windowSize[0], windowSize[1] = w.GetSize()
// 	}
// 	window.SetFramebufferSizeCallback(framebufferSizeCallback)
// 	var framebufferSize [2]int
// 	framebufferSize[0], framebufferSize[1] = window.GetFramebufferSize()
// 	framebufferSizeCallback(window, framebufferSize[0], framebufferSize[1])
//
// 	// Clear Screen Color
// 	gl.ClearColor(0.8, 0.3, 0.01, 1)
//
// 	// Create and Set Active a new Program using our custom Shaders
// 	program, err := glutil.CreateProgram(vertexSource, fragmentSource)
// 	fatal(err)
//
// 	gl.ValidateProgram(program)
// 	if gl.GetProgrami(program, gl.VALIDATE_STATUS) != gl.TRUE {
// 		fatal(fmt.Errorf("gl validate status: %s", gl.GetProgramInfoLog(program)))
// 	}
// 	gl.UseProgram(program)
//
// 	// Get Uniform Locations
// 	pMatrixUniform := gl.GetUniformLocation(program, "uPMatrix")
// 	mvMatrixUniform := gl.GetUniformLocation(program, "uMVMatrix")
//
// 	// Load Triangle Vertex Data into Shaders
// 	triangleVertexPositionBuffer := gl.CreateBuffer()
// 	gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexPositionBuffer)
// 	vertices := f32.Bytes(binary.LittleEndian,
// 		0, 0, 0,
// 		300, 100, 0,
// 		0, 100, 0,
// 	)
// 	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)
// 	itemSize := 3
// 	itemCount := 3
//
// 	// Setup Vertex Attribute Arrays
// 	vertexPositionAttrib := gl.GetAttribLocation(program, "aVertexPosition")
// 	gl.EnableVertexAttribArray(vertexPositionAttrib)
// 	gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)
//
// 	// Check for Errors
// 	if err := gl.GetError(); err != 0 {
// 		fatal(fmt.Errorf("gl error: %v", err))
// 	}
//
// 	// Main Render Loop
// 	for !window.ShouldClose() {
//
// 		// Clear
// 		gl.Clear(gl.COLOR_BUFFER_BIT)
//
// 		// Get Perspective Transformation
// 		pMatrix := mgl32.Ortho2D(0, float32(windowSize[0]), float32(windowSize[1]), 0)
//
// 		// Get Model View Matrix
// 		mvMatrix := mgl32.Translate3D(cursorPos[0], cursorPos[1], 0)
//
// 		// Load Uniform Values for Render
// 		gl.UniformMatrix4fv(pMatrixUniform, pMatrix[:])
// 		gl.UniformMatrix4fv(mvMatrixUniform, mvMatrix[:])
// 		gl.DrawArrays(gl.TRIANGLES, 0, itemCount)
//
// 		// Draw New Scene
// 		window.SwapBuffers()
//
// 		// Poll for GL Interaction Events
// 		glfw.PollEvents()
// 	}
// }

// type mat4 [4][4]float32
//
// func (m mat4) mul(m2 mat4) mat4 {
// 	var mres mat4
// 	for i := 0; i < 4; i++ {
// 		for j := 0; j < 4; j++ {
// 			var s float32
// 			for k := 0; k < 4; k++ {
// 				s += m[i][k] * m2[k][j]
// 			}
// 			mres[i][j] = s
// 		}
// 	}
// 	return mres
// }
//
// func (m mat4) vec() [16]float32 {
// 	var res [16]float32
// 	for i := 0; i < 4; i++ {
// 		for j := 0; j < 4; j++ {
// 			res[4*i+j] = m[i][j]
// 		}
// 	}
// 	return res
// }
//
// func (m mat4) String() string {
// 	const format = `[
// 	[ %v, %v, %v, %v ],
// 	[ %v, %v, %v, %v ],
// 	[ %v, %v, %v, %v ],
// 	[ %v, %v, %v, %v ],
// ]
// `
// 	var vs = make([]interface{}, 16)
// 	for i, f := range m.vec() {
// 		vs[i] = f
// 	}
// 	return fmt.Sprintf(format, vs...)
// }
//
// func test2() {
// 	m1 := mat4{
// 		{1.0, 0.0, 0.0},
// 		{0.0, 1.0, 0.0},
// 		{0.0, 0.0, 1.0},
// 	}
//
// 	m2 := mat4{
// 		{1.0, 0.0, 0.0},
// 		{0.0, 1.0, 0.0},
// 		{0.0, 0.0, 1.0},
// 	}
//
// 	m3 := m1.mul(m2)
// 	fmt.Printf("m1 * m2 = %v", m3)
// }

type GLEnum uint32

const (
	VENDOR   GLEnum = gl.VENDOR
	RENDERER        = gl.RENDERER
	VERSION         = gl.VERSION
	SAMPLES         = gl.SAMPLES

	SHADING_LANGUAGE_VERSION = gl.SHADING_LANGUAGE_VERSION
)

func dumpGLConfig() {
	fmt.Printf("OpenGL:\n")
	fmt.Printf(" * Vendor:    %s\n", GetString(gl.VENDOR))
	fmt.Printf(" * Renderer:  %s\n", GetString(gl.RENDERER))
	fmt.Printf(" * Version:   %s\n", GetString(gl.VERSION))
	fmt.Printf(" * Samples:   %v\n", GetString(gl.SAMPLES))

	fmt.Printf("GLSL:\n")
	fmt.Printf(" * Shading Language Version: %s\n", GetString(gl.SHADING_LANGUAGE_VERSION))
}

func Str(str string) *uint8 {
	if !strings.HasSuffix(str, "\x00") {
		panic("str argument missing null terminator: " + str)
	}
	header := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return (*uint8)(unsafe.Pointer(header.Data))
}

func GoString(cstr *uint8) string {
	return C.GoString((*C.char)(unsafe.Pointer(cstr)))
}

func GetString(n uint32) string {
	return GoString(gl.GetString(n))
}

// type Program uint32

// type Shader uint32

func GetShaderInfoLog(s uint32) string {
	var logLength int32
	gl.GetShaderiv(uint32(s), gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetShaderInfoLog(uint32(s), logLength, nil, &logBuffer[0])
	return GoString(&logBuffer[0])
}

func GetProgramInfoLog(p uint32) string {
	var logLength int32
	gl.GetProgramiv(p, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetProgramInfoLog(p, logLength, nil, &logBuffer[0])
	return GoString(&logBuffer[0])
}

func LoadShader(typ GLEnum, src string) (uint32, error) {
	// shader := gl.CreateShader(typ)
	shader := gl.CreateShader(uint32(typ))
	if shader == 0 {
		return 0, fmt.Errorf("glutil: could not create shader (type %v)", typ)
	}

	// gl.ShaderSource(shader, src)
	glsource, free := gl.Strs(src + "\x00")
	gl.ShaderSource(shader, 1, glsource, nil)
	free()

	gl.CompileShader(shader)

	var shaderi int32
	gl.GetShaderiv(shader, uint32(typ), &shaderi)
	if shaderi == 0 {
		defer gl.DeleteShader(shader)
		return 0, fmt.Errorf("shader compile: %s", GetShaderInfoLog(shader))
	}

	return 0, nil
}

func CreateProgram(vs, fs string) (uint32, error) {
	p := gl.CreateProgram()
	if p == 0 {
		return 0, fmt.Errorf("glutil: no programs available")
	}

	vshader, err := LoadShader(gl.VERTEX_SHADER, vs)
	if err != nil {
		return 0, err
	}

	fshader, err := LoadShader(gl.FRAGMENT_SHADER, fs)
	if err != nil {
		gl.DeleteShader(vshader)
		return 0, err
	}

	gl.AttachShader(p, vshader)
	gl.AttachShader(p, fshader)
	gl.LinkProgram(p)

	gl.DeleteShader(vshader)
	gl.DeleteShader(fshader)

	var programi int32
	gl.GetProgramiv(p, gl.LINK_STATUS, &programi)
	if programi == 0 {
		defer gl.DeleteProgram(p)
		return 0, fmt.Errorf("glutil: %s", GetProgramInfoLog(p))
	}

	return p, nil
}

func goglTest() error {
	// Initialize GL
	glfw.Init()
	defer glfw.Terminate()

	// var m glfw.Monitor
	// var share glfw.Window
	var windowSize = [2]int{640, 480}
	glfw.WindowHint(glfw.Samples, 8)
	// w, err := glfw.CreateWindow(windowSize[0], windowSize[1], "Test Title", &m, &share)
	w, err := glfw.CreateWindow(windowSize[0], windowSize[1], "Test Title", nil, nil)
	if err != nil {
		return err
	}

	w.MakeContextCurrent()
	dumpGLConfig()

	// Set Cursor Change-Listener
	cursor := [2]float32{200, 200}
	w.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
		cursor[0], cursor[1] = float32(x), float32(y)
	})

	// Callback for when Framebuffer Changes;
	// Framebuffers are the objects that contain image and rendering data.
	// Assumedly this would be triggered on window resize
	w.SetFramebufferSizeCallback(func(w *glfw.Window, frameBufferSize0, frameBufferSize1 int) {
		x, y := 0, 0
		gl.Viewport(int32(x), int32(y), int32(frameBufferSize0), int32(frameBufferSize1))
		windowSize[0], windowSize[1] = w.GetSize()
	})

	// 	var framebufferSize [2]int
	// 	framebufferSize[0], framebufferSize[1] = window.GetFramebufferSize()
	// 	framebufferSizeCallback(window, framebufferSize[0], framebufferSize[1])

	// Clear Screen Color
	var red, green, blue, alpha float32
	red, green, blue, alpha = 0.8, 0.3, 0.01, 1
	gl.ClearColor(red, green, blue, alpha)

	// Create and Set Active a new Program using our custom Shaders
	program, err := CreateProgram(vertexSource, fragmentSource)
	if err != nil {
		return nil
	}

	gl.ValidateProgram(program)

	var programi int32
	gl.GetProgramiv(program, uint32(gl.VALIDATE_STATUS), &programi)
	if programi != gl.TRUE {
		return fmt.Errorf("gl validate status: %s", GetProgramInfoLog(program))
	}
	gl.UseProgram(program)

	// Get Uniform Locations
	pMatrixUniform := gl.GetUniformLocation(program, Str("uPMatrix\x00"))
	mvMatrixUniform := gl.GetUniformLocation(program, Str("uMVMatrix\x00"))

	// Load Triangle Vertex Data into Shaders
	// triangleVertexPositionBuffer := gl.CreateBuffer()
	var triangleVertexPositionBuffer uint32
	gl.GenBuffers(1, &triangleVertexPositionBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexPositionBuffer)
	vertices := f32.Bytes(binary.LittleEndian,
		0, 0, 0,
		300, 100, 0,
		0, 100, 0,
	)
	gl.BufferData(gl.ARRAY_BUFFER, int(len(vertices)), gl.Ptr(&vertices[0]), gl.STATIC_DRAW)
	var itemSize int32 = 3
	var itemCount int32 = 3

	// Setup Vertex Attribute Arrays
	vertexPositionAttrib := uint32(gl.GetAttribLocation(program, Str("aVertexPosition\x00")))
	gl.EnableVertexAttribArray(vertexPositionAttrib)
	var stride int32 = 0
	var offset = gl.PtrOffset(0)
	gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, stride, offset)

	// Check for Errors
	if err := gl.GetError(); err != 0 {
		return fmt.Errorf("gl error: %v", err)
	}

	// Main Render Loop
	for !w.ShouldClose() {

		// Clear
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Get Perspective Transformation
		pMatrix := mgl32.Ortho2D(0, float32(windowSize[0]), float32(windowSize[1]), 0)

		// Get Model View Matrix
		mvMatrix := mgl32.Translate3D(cursor[0], cursor[1], 0)

		// Load Uniform Values for Render
		gl.UniformMatrix4fv(pMatrixUniform, int32(len(pMatrix)/(4*4)), false, &pMatrix[0])
		gl.UniformMatrix4fv(mvMatrixUniform, int32(len(mvMatrix)/(4*4)), false, &mvMatrix[0])
		gl.DrawArrays(gl.TRIANGLES, int32(0), itemCount)

		// Draw New Scene
		w.SwapBuffers()

		// Poll for GL Interaction Events
		glfw.PollEvents()
	}

	return nil
}

func main() {
	// goxjsTest()
	// test2()

	// err := goglTest()
	// if err != nil {
	// 	panic(err)
	// }
}

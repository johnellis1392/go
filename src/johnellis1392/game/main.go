package main

import (
	"C"
	"encoding/binary"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/exp/f32"
)

// #include <stdlib.h>
import "C"

const vertexSource = `
#version 120 // OpenGL 2.1.
//#version 100 // WebGL.

attribute vec3 aVertexPosition;

uniform mat4 uMVMatrix;
uniform mat4 uPMatrix;

void main() {
	gl_Position = uPMatrix * uMVMatrix * vec4(aVertexPosition, 1.0);
}
`

const fragmentSource = `
#version 120 // OpenGL 2.1.
//#version 100 // WebGL.

void main() {
	gl_FragColor = vec4(1.0, 1.0, 1.0, 1.0);
}
`

func init() {
	runtime.LockOSThread()
}

//type GLEnum uint32

const (
	VENDOR   uint32 = gl.VENDOR
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
	// fmt.Printf(" * Samples:   %v\n", GetString(gl.SAMPLES))

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

// func GetShaderInfoLog(s uint32) string {
// 	var logLength int32
// 	gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &logLength)
// 	if logLength == 0 {
// 		return "[no log found]"
// 	}
//
// 	logBuffer := make([]uint8, logLength)
// 	gl.GetShaderInfoLog(s, logLength, nil, &logBuffer[0])
// 	return GoString(&logBuffer[0])
// }

// func GetProgramInfoLog(p uint32) string {
// 	var logLength int32
// 	gl.GetProgramiv(p, gl.INFO_LOG_LENGTH, &logLength)
// 	if logLength == 0 {
// 		return "[no log found]"
// 	}
//
// 	logBuffer := make([]uint8, logLength)
// 	gl.GetProgramInfoLog(p, logLength, nil, &logBuffer[0])
// 	return GoString(&logBuffer[0])
// }

// func LoadShader(typ uint32, src string) (uint32, error) {
// 	shader := gl.CreateShader(typ)
// 	if shader == 0 {
// 		return 0, fmt.Errorf("glutil: could not create shader (type %v)", typ)
// 	}
//
// 	glsource, free := gl.Strs(src + "\x00")
// 	gl.ShaderSource(shader, 1, glsource, nil)
// 	free()
//
// 	gl.CompileShader(shader)
//
// 	var shaderi int32
// 	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &shaderi)
// 	if shaderi == 0 {
// 		defer gl.DeleteShader(shader)
// 		return 0, fmt.Errorf("%s", GetShaderInfoLog(shader))
// 	}
//
// 	return shader, nil
// }

// func CreateProgram(vs, fs string) (uint32, error) {
// 	p := gl.CreateProgram()
// 	if p == 0 {
// 		return 0, fmt.Errorf("glutil: no programs available")
// 	}
//
// 	vshader, err := LoadShader(gl.VERTEX_SHADER, vs)
// 	if err != nil {
// 		return 0, fmt.Errorf("an error occurred while processing vertex shader: %s", err.Error())
// 	}
//
// 	fshader, err := LoadShader(gl.FRAGMENT_SHADER, fs)
// 	if err != nil {
// 		gl.DeleteShader(vshader)
// 		return 0, fmt.Errorf("an error occurred while processing fragment shader: %s", err.Error())
// 	}
//
// 	gl.AttachShader(p, vshader)
// 	gl.AttachShader(p, fshader)
// 	gl.LinkProgram(p)
//
// 	gl.DeleteShader(vshader)
// 	gl.DeleteShader(fshader)
//
// 	var programi int32
// 	gl.GetProgramiv(p, gl.LINK_STATUS, &programi)
// 	if programi == 0 {
// 		defer gl.DeleteProgram(p)
// 		return 0, fmt.Errorf("glutil: %s", GetProgramInfoLog(p))
// 	}
//
// 	return p, nil
// }

func Terminate() {
	glfw.Terminate()
	runtime.UnlockOSThread()
}

func goglTest() error {
	var err error

	// Initialize GL
	if err = glfw.Init(); err != nil {
		return err
	}
	defer Terminate()

	glfw.WindowHint(glfw.Hint(glfw.Samples), 8)
	var windowSize = [2]int{640, 480}
	w, err := glfw.CreateWindow(windowSize[0], windowSize[1], "Test Title", nil, nil)
	if err != nil {
		return err
	}

	w.MakeContextCurrent()
	if err = gl.Init(); err != nil {
		return err
	}

	dumpGLConfig()

	// Set Cursor Change-Listener
	cursor := [2]float32{200, 200}
	w.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
		cursor[0], cursor[1] = float32(x), float32(y)
	})

	// Callback for when Framebuffer Changes;
	// Framebuffers are the objects that contain image and rendering data.
	// Assumedly this would be triggered on window resize
	framebufferSizeCallback := func(w *glfw.Window, frameBufferSize0, frameBufferSize1 int) {
		x, y := 0, 0
		gl.Viewport(int32(x), int32(y), int32(frameBufferSize0), int32(frameBufferSize1))
		windowSize[0], windowSize[1] = w.GetSize()
	}
	w.SetFramebufferSizeCallback(framebufferSizeCallback)
	var framebufferSize [2]int
	framebufferSize[0], framebufferSize[1] = w.GetFramebufferSize()
	framebufferSizeCallback(w, framebufferSize[0], framebufferSize[1])

	// Clear Screen Color
	var red, green, blue, alpha float32
	red, green, blue, alpha = 0.8, 0.3, 0.01, 1
	gl.ClearColor(red, green, blue, alpha)

	// Create and Set Active a new Program using our custom Shaders
	program, err := CreateProgram(vertexSource, fragmentSource)
	if err != nil {
		return err
	}
	// program := CreateProgram()

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
	gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// Check for Errors
	if glerr := gl.GetError(); glerr != 0 {
		return fmt.Errorf("gl error: %v", glerr)
	} else {
		fmt.Printf("No Errors Found...\n")
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

// func simpleExample() {
// 	err := glfw.Init()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer glfw.Terminate()
//
// 	w, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	w.MakeContextCurrent()
//
// 	for !w.ShouldClose() {
// 		w.SwapBuffers()
// 		glfw.PollEvents()
// 	}
// }
//
// func simpleExample2() {
// 	err := glfw.Init()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer glfw.Terminate()
//
// 	w, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	w.MakeContextCurrent()
//
// 	// Need to Re-Init GL after changing Context
// 	if err = gl.Init(); err != nil {
// 		panic(err)
// 	}
//
// 	gl.ClearColor(1.0, 0.0, 0.0, 0.5)
// 	gl.Clear(gl.COLOR_BUFFER_BIT)
//
// 	for !w.ShouldClose() {
// 		gl.Clear(gl.COLOR_BUFFER_BIT)
//
// 		w.SwapBuffers()
// 		glfw.PollEvents()
// 	}
// }

func main() {
	// goxjsTest()

	if err := goglTest(); err != nil {
		panic(err)
	}

	// simpleExample()

	// simpleExample2()
}

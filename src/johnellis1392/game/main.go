package main

import (
	"C"
	"encoding/binary"
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/exp/f32"
)

// const vertexSource = `
// #version 120 // OpenGL 2.1.
// //#version 100 // WebGL.
//
// attribute vec3 aVertexPosition;
//
// uniform mat4 uMVMatrix;
// uniform mat4 uPMatrix;
//
// void main() {
// 	gl_Position = uPMatrix * uMVMatrix * vec4(aVertexPosition, 1.0);
// }
// `

// const fragmentSource = `
// #version 120 // OpenGL 2.1.
// //#version 100 // WebGL.
//
// void main() {
// 	gl_FragColor = vec4(1.0, 1.0, 1.0, 1.0);
// }
// `

const (
	vertexShaderLoc   = "./main.vertex.glsl"
	fragmentShaderLoc = "./main.fragment.glsl"
)

func init() {
	runtime.LockOSThread()
}

func dumpGLConfig() {
	fmt.Printf("OpenGL:\n")
	fmt.Printf(" * Vendor:    %s\n", GetString(gl.VENDOR))
	fmt.Printf(" * Renderer:  %s\n", GetString(gl.RENDERER))
	fmt.Printf(" * Version:   %s\n", GetString(gl.VERSION))
	// fmt.Printf(" * Samples:   %v\n", GetString(gl.SAMPLES))

	fmt.Printf("GLSL:\n")
	fmt.Printf(" * Shading Language Version: %s\n", GetString(gl.SHADING_LANGUAGE_VERSION))
}

func Terminate() {
	glfw.Terminate()
	runtime.UnlockOSThread()
}

// func goglTest() error {
// 	var err error
//
// 	// Initialize GL
// 	if err = glfw.Init(); err != nil {
// 		return err
// 	}
// 	defer Terminate()
//
// 	glfw.WindowHint(glfw.Hint(glfw.Samples), 8)
// 	var windowSize = [2]int{640, 480}
// 	w, err := glfw.CreateWindow(windowSize[0], windowSize[1], "Test Title", nil, nil)
// 	if err != nil {
// 		return err
// 	}
//
// 	w.MakeContextCurrent()
// 	if err = gl.Init(); err != nil {
// 		return err
// 	}
//
// 	dumpGLConfig()
//
// 	// Set Cursor Change-Listener
// 	cursor := [2]float32{200, 200}
// 	w.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
// 		cursor[0], cursor[1] = float32(x), float32(y)
// 	})
//
// 	// Callback for when Framebuffer Changes;
// 	// Framebuffers are the objects that contain image and rendering data.
// 	// Assumedly this would be triggered on window resize
// 	framebufferSizeCallback := func(w *glfw.Window, frameBufferSize0, frameBufferSize1 int) {
// 		x, y := 0, 0
// 		gl.Viewport(int32(x), int32(y), int32(frameBufferSize0), int32(frameBufferSize1))
// 		windowSize[0], windowSize[1] = w.GetSize()
// 	}
// 	w.SetFramebufferSizeCallback(framebufferSizeCallback)
// 	var framebufferSize [2]int
// 	framebufferSize[0], framebufferSize[1] = w.GetFramebufferSize()
// 	framebufferSizeCallback(w, framebufferSize[0], framebufferSize[1])
//
// 	// Clear Screen Color
// 	var red, green, blue, alpha float32
// 	red, green, blue, alpha = 0.8, 0.3, 0.01, 1
// 	gl.ClearColor(red, green, blue, alpha)
//
// 	// Create and Set Active a new Program using our custom Shaders
// 	program, err := CreateProgram(vertexSource, fragmentSource)
// 	if err != nil {
// 		return err
// 	}
//
// 	gl.ValidateProgram(program)
//
// 	var programi int32
// 	gl.GetProgramiv(program, uint32(gl.VALIDATE_STATUS), &programi)
// 	if programi != gl.TRUE {
// 		return fmt.Errorf("gl validate status: %s", GetProgramInfoLog(program))
// 	}
// 	gl.UseProgram(program)
//
// 	// Get Uniform Locations
// 	pMatrixUniform := gl.GetUniformLocation(program, Str("uPMatrix\x00"))
// 	mvMatrixUniform := gl.GetUniformLocation(program, Str("uMVMatrix\x00"))
//
// 	// Load Triangle Vertex Data into Shaders
// 	// triangleVertexPositionBuffer := gl.CreateBuffer()
// 	var triangleVertexPositionBuffer uint32
// 	gl.GenBuffers(1, &triangleVertexPositionBuffer)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexPositionBuffer)
// 	vertices := f32.Bytes(binary.LittleEndian,
// 		0, 0, 0,
// 		300, 100, 0,
// 		0, 100, 0,
// 	)
// 	gl.BufferData(gl.ARRAY_BUFFER, int(len(vertices)), gl.Ptr(&vertices[0]), gl.STATIC_DRAW)
// 	var itemSize int32 = 3
// 	var itemCount int32 = 3
//
// 	// Setup Vertex Attribute Arrays
// 	vertexPositionAttrib := uint32(gl.GetAttribLocation(program, Str("aVertexPosition\x00")))
// 	gl.EnableVertexAttribArray(vertexPositionAttrib)
// 	gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, 0, gl.PtrOffset(0))
//
// 	// Check for Errors
// 	if glerr := gl.GetError(); glerr != 0 {
// 		return fmt.Errorf("gl error: %v", glerr)
// 	} else {
// 		fmt.Printf("No Errors Found...\n")
// 	}
//
// 	// Main Render Loop
// 	for !w.ShouldClose() {
//
// 		// Clear
// 		gl.Clear(gl.COLOR_BUFFER_BIT)
//
// 		// Get Perspective Transformation
// 		pMatrix := mgl32.Ortho2D(0, float32(windowSize[0]), float32(windowSize[1]), 0)
//
// 		// Get Model View Matrix
// 		mvMatrix := mgl32.Translate3D(cursor[0], cursor[1], 0)
//
// 		// Load Uniform Values for Render
// 		gl.UniformMatrix4fv(pMatrixUniform, int32(len(pMatrix)/(4*4)), false, &pMatrix[0])
// 		gl.UniformMatrix4fv(mvMatrixUniform, int32(len(mvMatrix)/(4*4)), false, &mvMatrix[0])
// 		gl.DrawArrays(gl.TRIANGLES, int32(0), itemCount)
//
// 		// Draw New Scene
// 		w.SwapBuffers()
//
// 		// Poll for GL Interaction Events
// 		glfw.PollEvents()
// 	}
//
// 	return nil
// }

func goglTest2() error {
	var err error

	// Initialize GL
	if err = glfw.Init(); err != nil {
		return err
	}
	defer Terminate()

	window, err := CreateWindow()
	if err != nil {
		return err
	}

	// Set GL Context & Initialize GL
	window.MakeContextCurrent()

	// Clear Screen Color
	window.Clear(Color{0.8, 0.3, 0.01, 1})

	// Create and Set Active a new Program using our custom Shaders
	program := CreateProgram()

	// Compile Program
	program.AddShader(gl.VERTEX_SHADER, vertexShaderLoc)
	program.AddShader(gl.FRAGMENT_SHADER, fragmentShaderLoc)
	if err := program.Compile(); err != nil {
		return err
	}

	if err := program.Validate(); err != nil {
		return err
	}

	// Set Current Program
	program.Use()

	// Get Uniform Locations
	pMatrixUniform := gl.GetUniformLocation(program.ID, Str("uPMatrix\x00"))
	mvMatrixUniform := gl.GetUniformLocation(program.ID, Str("uMVMatrix\x00"))

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
	vertexPositionAttrib := uint32(gl.GetAttribLocation(program.ID, Str("aVertexPosition\x00")))
	gl.EnableVertexAttribArray(vertexPositionAttrib)
	gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// Main Render Loop
	for !window.ShouldClose() {

		// Clear
		// TODO: Replace this with API
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Get Perspective Transformation
		ww, wh := float32(window.Size.Width), float32(window.Size.Height)
		pMatrix := mgl32.Ortho2D(0, ww, wh, 0)

		// Get Model View Matrix
		cx, cy := float32(window.Cursor.X), float32(window.Cursor.Y)
		mvMatrix := mgl32.Translate3D(cx, cy, 0)

		// Load Uniform Values for Render
		gl.UniformMatrix4fv(pMatrixUniform, int32(len(pMatrix)/(4*4)), false, &pMatrix[0])
		gl.UniformMatrix4fv(mvMatrixUniform, int32(len(mvMatrix)/(4*4)), false, &mvMatrix[0])
		gl.DrawArrays(gl.TRIANGLES, int32(0), itemCount)

		// Draw New Scene
		window.SwapBuffers()

		// Poll for GL Interaction Events
		window.PollEvents()
	}

	return nil
}

func main() {
	if err := goglTest2(); err != nil {
		panic(err)
	}
}

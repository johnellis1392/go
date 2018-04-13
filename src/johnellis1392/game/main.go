package main

import (
	"encoding/binary"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/gl/glutil"
	"github.com/goxjs/glfw"
	"golang.org/x/mobile/exp/f32"
)

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

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func dumpGLConfig() {
	fmt.Printf("OpenGL:\n")
	fmt.Printf(" * Vendor:    %s\n", gl.GetString(gl.VENDOR))
	fmt.Printf(" * Renderer:  %s\n", gl.GetString(gl.RENDERER))
	fmt.Printf(" * Version:   %s\n", gl.GetString(gl.VERSION))
	fmt.Printf(" * Samples:   %v\n", gl.GetString(gl.SAMPLES))

	fmt.Printf("GLSL:\n")
	fmt.Printf(" * Shading Language Version: %s\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))
}

func goxjsTest() {
	// Initialize GL
	fatal(glfw.Init(gl.ContextWatcher))
	defer glfw.Terminate()

	// Create Window Object
	var windowSize = [2]int{640, 480}
	glfw.WindowHint(glfw.Samples, 8)
	window, err := glfw.CreateWindow(windowSize[0], windowSize[1], "", nil, nil)
	fatal(err)

	window.MakeContextCurrent()

	// Print GL Configuration Information
	dumpGLConfig()

	// Logic for Moving Mouse Cursor
	cursorPos := [2]float32{200, 200}
	cursorPosCallback := func(_ *glfw.Window, x, y float64) {
		cursorPos[0], cursorPos[1] = float32(x), float32(y)
	}
	window.SetCursorPosCallback(cursorPosCallback)

	// Callback for when Framebuffer Changes;
	// Framebuffers are the objects that contain image and rendering data.
	// Assumedly this would be triggered on window resize
	framebufferSizeCallback := func(w *glfw.Window, framebufferSize0, framebufferSize1 int) {
		gl.Viewport(0, 0, framebufferSize0, framebufferSize1)
		windowSize[0], windowSize[1] = w.GetSize()
	}
	window.SetFramebufferSizeCallback(framebufferSizeCallback)
	var framebufferSize [2]int
	framebufferSize[0], framebufferSize[1] = window.GetFramebufferSize()
	framebufferSizeCallback(window, framebufferSize[0], framebufferSize[1])

	// Clear Screen Color
	gl.ClearColor(0.8, 0.3, 0.01, 1)

	// Create and Set Active a new Program using our custom Shaders
	program, err := glutil.CreateProgram(vertexSource, fragmentSource)
	fatal(err)

	gl.ValidateProgram(program)
	if gl.GetProgrami(program, gl.VALIDATE_STATUS) != gl.TRUE {
		fatal(fmt.Errorf("gl validate status: %s", gl.GetProgramInfoLog(program)))
	}
	gl.UseProgram(program)

	// Get Uniform Locations
	pMatrixUniform := gl.GetUniformLocation(program, "uPMatrix")
	mvMatrixUniform := gl.GetUniformLocation(program, "uMVMatrix")

	// Load Triangle Vertex Data into Shaders
	triangleVertexPositionBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexPositionBuffer)
	vertices := f32.Bytes(binary.LittleEndian,
		0, 0, 0,
		300, 100, 0,
		0, 100, 0,
	)
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)
	itemSize := 3
	itemCount := 3

	// Setup Vertex Attribute Arrays
	vertexPositionAttrib := gl.GetAttribLocation(program, "aVertexPosition")
	gl.EnableVertexAttribArray(vertexPositionAttrib)
	gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)

	// Check for Errors
	if err := gl.GetError(); err != 0 {
		fatal(fmt.Errorf("gl error: %v", err))
	}

	// Main Render Loop
	for !window.ShouldClose() {

		// Clear
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Get Perspective Transformation
		pMatrix := mgl32.Ortho2D(0, float32(windowSize[0]), float32(windowSize[1]), 0)

		// Get Model View Matrix
		mvMatrix := mgl32.Translate3D(cursorPos[0], cursorPos[1], 0)

		// Load Uniform Values for Render
		gl.UniformMatrix4fv(pMatrixUniform, pMatrix[:])
		gl.UniformMatrix4fv(mvMatrixUniform, mvMatrix[:])
		gl.DrawArrays(gl.TRIANGLES, 0, itemCount)

		// Draw New Scene
		window.SwapBuffers()

		// Poll for GL Interaction Events
		glfw.PollEvents()
	}
}

func main() {
	goxjsTest()
}

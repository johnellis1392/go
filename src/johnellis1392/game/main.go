package main

import (
	"C"
	"encoding/binary"
	"fmt"
	"johnellis1392/game/gl"
	"johnellis1392/game/math"
	"runtime"
	"time"
)
import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/glfw"
	"golang.org/x/mobile/exp/f32"
)

const (
	vertexShaderLoc   = "./main.vertex.glsl"
	fragmentShaderLoc = "./main.fragment.glsl"
)

func init() {
	runtime.LockOSThread()
}

func dumpGLConfig() {
	fmt.Printf("OpenGL:\n")
	fmt.Printf(" * Vendor:    %s\n", gl.GetString(gl.Vendor))
	fmt.Printf(" * Renderer:  %s\n", gl.GetString(gl.Renderer))
	fmt.Printf(" * Version:   %s\n", gl.GetString(gl.Version))

	fmt.Printf("GLSL:\n")
	fmt.Printf(" * Shading Language Version: %s\n", gl.GetString(gl.ShadingLanguageVersion))
}

func Terminate() {
	glfw.Terminate()
	runtime.UnlockOSThread()
}

type Triangle struct {
	GL                   GLPos
	Vertices             [3]math.Vec3
	Program              gl.Program
	VertexPositionBuffer uint32
	VertexPositionAttrib uint32
}

func newTriangle(p gl.Program) Triangle {
	return Triangle{
		GL: GLPos{
			Pos:   math.Mat4Id(),
			Rot:   math.Mat4Id(),
			Scale: math.Mat4Id(),
		},
		Vertices: [3]math.Vec3{
			{0, 0, 0},
			{300, 100, 0},
			{0, 100, 0},
		},
		Program:              p,
		VertexPositionBuffer: 0, // Null for Now
		VertexPositionAttrib: 0,
	}
}

func spread(vs [3]math.Vec3) []float32 {
	var res []float32
	for _, v := range vs {
		res = append(res, v[0])
		res = append(res, v[1])
		res = append(res, v[2])
	}
	return res
}

// Init Init
func (t Triangle) Init(c gl.Context) error {
	// TODO:
	// * Bind vertex attrib array for points
	// * Bind Program / Bind Shaders
	// * Create Buffer & Load with Vertices
	// return nil

	// Create Triangle Vertex Buffer
	var triangleVertexPositionBuffer uint32
	gl.GenBuffers(1, &triangleVertexPositionBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexPositionBuffer)
	t.VertexPositionBuffer = triangleVertexPositionBuffer

	// Load Triangle Vertices
	vertices := f32.Bytes(binary.LittleEndian, spread(t.Vertices)...)
	gl.BufferData(gl.ARRAY_BUFFER, int(len(vertices)), gl.Ptr(&vertices[0]), gl.STATIC_DRAW)
	var itemSize int32 = 3

	// Setup Vertex Attribute Arrays
	vertexPositionAttrib := uint32(gl.GetAttribLocation(t.Program.ID, gl.Str("aVertexPosition\x00")))
	gl.EnableVertexAttribArray(vertexPositionAttrib)
	if glerr := gl.GetError(); glerr != 0 {
		return fmt.Errorf("gl error: %v", glerr)
	}

	gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, 0, gl.PtrOffset(0))
	t.VertexPositionAttrib = vertexPositionAttrib

	return nil
}

// Update Update
func (t Triangle) Update(dt time.Time) {
	// Set VAO to Window Cursor Position
}

// Render Render
func (t Triangle) Render(c gl.Context) {
	// Call GLDrawArrays
	gl.DrawArrays(gl.TRIANGLES, int32(0), int32(len(t.Vertices)))
}

// Destroy Destroy
func (t Triangle) Destroy() error {
	return nil
}

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

	// Set Current Program
	program.Use()

	// Get Uniform Locations
	pMatrixUniform := gl.GetUniformLocation(program.ID, Str("uPMatrix\x00"))
	mvMatrixUniform := gl.GetUniformLocation(program.ID, Str("uMVMatrix\x00"))

	var context gl.Context
	context = window.Context()

	// Create Triangle
	t := newTriangle(program)
	if err := t.Init(context); err != nil {
		return err
	}

	var updateTime time.Time
	updateTime = time.Now()

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

		// gl.DrawArrays(gl.TRIANGLES, int32(0), itemCount)
		// Render Triangle
		context = window.Context()
		t.Update(updateTime)
		t.Render(context)

		// Draw New Scene
		window.SwapBuffers()

		// Poll for GL Interaction Events
		window.PollEvents()

		// Recalculate Time
		dt := time.Since(updateTime)
		updateTime.Add(dt)

	}

	if err := t.Destroy(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := goglTest2(); err != nil {
		panic(err)
	}
}

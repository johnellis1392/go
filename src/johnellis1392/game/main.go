package main

import (
	"C"
	"fmt"
	"johnellis1392/game/gl"
	"johnellis1392/game/math"
	"runtime"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/glfw"
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

// Test Triangle implements GameObject interface
var _ GameObject = (Triangle*)(nil)

type Triangle struct {
	GL       GLPos
	Vertices [3]math.Vec3
	Program  gl.Program
	// VertexPositionBuffer uint32
	// VertexPositionAttrib uint32

	VertexBuffer gl.Buffer
	VertexAttrib gl.Attrib
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
		Program: p,
		// VertexPositionBuffer: 0, // Null for Now
		// VertexPositionAttrib: 0,
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
// func (t Triangle) Init(c gl.Context) error {
// 	// TODO:
// 	// * Bind vertex attrib array for points
// 	// * Bind Program / Bind Shaders
// 	// * Create Buffer & Load with Vertices
// 	// return nil
//
// 	// Create Triangle Vertex Buffer
// 	var triangleVertexPositionBuffer uint32
// 	gl.GenBuffers(1, &triangleVertexPositionBuffer)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexPositionBuffer)
// 	t.VertexPositionBuffer = triangleVertexPositionBuffer
//
// 	// Load Triangle Vertices
// 	vertices := f32.Bytes(binary.LittleEndian, spread(t.Vertices)...)
// 	gl.BufferData(gl.ARRAY_BUFFER, int(len(vertices)), gl.Ptr(&vertices[0]), gl.STATIC_DRAW)
// 	var itemSize int32 = 3
//
// 	// Setup Vertex Attribute Arrays
// 	vertexPositionAttrib := uint32(gl.GetAttribLocation(t.Program.ID, gl.Str("aVertexPosition\x00")))
// 	gl.EnableVertexAttribArray(vertexPositionAttrib)
// 	if glerr := gl.GetError(); glerr != 0 {
// 		return fmt.Errorf("gl error: %v", glerr)
// 	}
//
// 	gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, 0, gl.PtrOffset(0))
// 	t.VertexPositionAttrib = vertexPositionAttrib
//
// 	return nil
// }

// Init initializes the Triangle's data in GL.
func (t Triangle) Init(c gl.Context) error {

	t.VertexBuffer = c.CreateBuffer(gl.ArrayBuffer)
	t.VertexBuffer.Bind()

	t.VertexAttrib = c.GetAttrib("aVertexPosition")
	t.VertexAttrib.Enable()
	t.VertexAttrib.Bind(len(t.Vertices))

	return nil
}

// Update Update
func (t Triangle) Update(dt time.Time) {
	// Set VAO to Window Cursor Position
}

// Render Render
func (t Triangle) Render(c gl.Context) {
	// Call GLDrawArrays
	// gl.DrawArrays(gl.TRIANGLES, int32(0), int32(len(t.Vertices)))
	c.Draw(len(t.Vertices))
}

// Destroy Destroy
func (t Triangle) Destroy() error {
	return nil
}

func goglTest2() error {
	var err error

	// Initialize GL
	if err = gl.Init(); err != nil {
		return err
	}
	defer gl.Terminate()

	window, err := gl.CreateWindow()
	if err != nil {
		return err
	}

	window.MakeContextCurrent()
	var backgroundColor = math.Color{R: 0.8, G: 0.3, B: 0.01, A: 1}
	window.Clear(backgroundColor)
	program := gl.CreateProgram()

	// Compile Program
	program.AddShader(gl.VertexShader, vertexShaderLoc)
	program.AddShader(gl.FragmentShader, fragmentShaderLoc)
	if err := program.Compile(); err != nil {
		return err
	}

	program.Use()

	var context gl.Context
	context = window.Context()

	// Get Uniform Locations
	// pMatrixUniform := gl.GetUniformLocation(program.ID, Str("uPMatrix\x00"))
	// mvMatrixUniform := gl.GetUniformLocation(program.ID, Str("uMVMatrix\x00"))
	pMatrixUniform := context.GetUniform("uPMatrix")
	mvMatrixUniform := context.GetUniform("uMVMatrix")

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
		// gl.Clear(gl.COLOR_BUFFER_BIT)
		window.Clear(backgroundColor)

		// Get Perspective Transformation
		ww, wh := float32(window.Size.W), float32(window.Size.H)
		pMatrix := mgl32.Ortho2D(0, ww, wh, 0)

		// Get Model View Matrix
		cx, cy := float32(window.Cursor.X), float32(window.Cursor.Y)
		mvMatrix := mgl32.Translate3D(cx, cy, 0)

		// Load Uniform Values for Render
		// gl.UniformMatrix4fv(pMatrixUniform, int32(len(pMatrix)/(4*4)), false, &pMatrix[0])
		// gl.UniformMatrix4fv(mvMatrixUniform, int32(len(mvMatrix)/(4*4)), false, &mvMatrix[0])
		pMatrixUniform.Set(pMatrix)
		mvMatrixUniform.Set(mvMatrix)

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

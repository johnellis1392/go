package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const vertexShaderSource = `
#version 410
in vec3 vp;
void main() {
	gl_Position = vec4(vp, 1.0);
	//gl_Position = gl_ProjectionMatrix * gl_ModelViewMatrix * gl_Vertex;
}
` + "\x00"

const fragmentShaderSource = `
#version 410
out vec4 frag_colour;
void main() {
	frag_colour = vec4(1, 1, 1, 1);
}
` + "\x00"

const (
	width     = 500
	height    = 500
	threshold = 0.15
)

var (
	triangle = []float32{
		0, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}

	square = []float32{
		-0.5, 0.5, 0.0,
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,

		-0.5, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, -0.5, 0.0,
	}
)

type GLContext struct{}

func (ctx GLContext) BindBuffer(bufid int) {
	// ...
}

func (ctx GLContext) Draw() {
	// ...
}

type tempNode struct {
	bufid int
}

func (n *tempNode) Width() int  { return 100 }
func (n *tempNode) Height() int { return 100 }

func (n *tempNode) renderGl(ctx GLContext) {
	ctx.BindBuffer(n.bufid)
	ctx.Draw()
}

// //////////////////////////////////////// //
// //////////////////////////////////////// //
// //////////////////////////////////////// //

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Example App", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL Version:", version)

	vshader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fshader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	p := gl.CreateProgram()
	gl.AttachShader(p, vshader)
	gl.AttachShader(p, fshader)
	gl.LinkProgram(p)
	return p
}

// func draw(vao uint32, w *glfw.Window, p uint32) {
func draw(cells [][]*cell, w *glfw.Window, p uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(p)

	// gl.BindVertexArray(vao)
	// gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

	for x := range cells {
		for _, c := range cells[x] {
			c.draw()
		}
	}

	glfw.PollEvents()
	w.SwapBuffers()
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// ////////////////////////// //
// Conways Game of Life Stuff //
// ////////////////////////// //

type cell struct {
	drawable         uint32
	x, y             int
	alive, aliveNext bool
}

func (c *cell) draw() {
	if !c.alive {
		return
	}

	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

func (c *cell) checkState(cells [][]*cell) {
	c.alive = c.aliveNext
	c.aliveNext = c.alive

	liveCount := c.liveNeighbors(cells)
	if c.alive {
		// 1. Any live cell with fewer than two live neighbors dies
		if liveCount < 2 {
			c.aliveNext = false
		}

		// 2. Any live cells with two or three live neighbors lives
		if liveCount == 2 || liveCount == 3 {
			c.aliveNext = true
		}

		// 3. Any live cell with more than three live neighbors dies
		if liveCount > 3 {
			c.aliveNext = false
		}
	} else {
		// 4. Any dead cell with exactly three live neighbors becomes a live cell
		if liveCount == 3 {
			c.aliveNext = true
		}
	}
}

func (c *cell) liveNeighbors(cells [][]*cell) int {
	var liveCount int
	add := func(x, y int) {
		if x == len(cells) {
			x = 0
		} else if x == -1 {
			x = len(cells) - 1
		}

		if y == len(cells[x]) {
			y = 0
		} else if y == -1 {
			y = len(cells[x]) - 1
		}

		if cells[x][y].alive {
			liveCount++
		}
	}

	add(c.x-1, c.y)
	add(c.x+1, c.y)

	add(c.x, c.y-1)
	add(c.x, c.y+1)

	add(c.x-1, c.y+1)
	add(c.x+1, c.y+1)
	add(c.x-1, c.y-1)
	add(c.x+1, c.y-1)

	return liveCount
}

const (
	rows    = 10
	columns = 10
)

func newCell(x, y int) *cell {
	points := make([]float32, len(square), len(square))
	copy(points, square)

	for i := 0; i < len(points); i++ {
		var position, size float32
		switch i % 3 {
		case 0:
			size = 1.0
			position = float32(x) * size
		case 1:
			size = 1.0 / float32(rows)
			position = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (position * 2) - 1
		} else {
			points[i] = ((position + size) * 2) - 1
		}
	}

	return &cell{
		drawable: makeVao(points),
		x:        x,
		y:        y,
	}
}

func makeCells() [][]*cell {
	rand.Seed(time.Now().UnixNano())

	cells := make([][]*cell, rows, rows)
	for x := 0; x < rows; x++ {
		for y := 0; y < columns; y++ {
			c := newCell(x, y)
			c.alive = rand.Float64() < threshold
			c.aliveNext = c.alive
			cells[x] = append(cells[x], c)
		}
	}

	return cells
}

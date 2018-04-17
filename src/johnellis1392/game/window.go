package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	SAMPLES        = 0
	DEFAULT_WIDTH  = 640
	DEFAULT_HEIGHT = 480
	TITLE          = "Test Title"
)

type Color struct {
	R, G, B, A float32
}

type Point struct {
	X, Y float64
}

type Dim struct {
	Width, Height int
}

type Window struct {
	Pos       Point
	Size      Dim
	FrameSize Dim
	Cursor    Point
	w         *glfw.Window

	cursorPositionCallback  cursorPositionCallback
	framebufferSizeCallback framebufferSizeCallback
}

func NewWindow() (Window, error) {
	title := TITLE
	w, h, s := DEFAULT_WIDTH, DEFAULT_HEIGHT, SAMPLES
	var err error

	glfw.WindowHint(glfw.Hint(glfw.Samples), s)

	glw, err := glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		return Window{}, err
	}

	glw.MakeContextCurrent()
	if err = gl.Init(); err != nil {
		return Window{}, err
	}

	const cursorX, cursorY = 200, 200
	window := Window{
		w:      glw,
		Size:   Dim{w, h},
		Cursor: Point{cursorX, cursorY},
	}

	return window, nil
}

func (w Window) PollEvents() {
	glfw.PollEvents()
}

func (w Window) Clear(c Color) {
	gl.ClearColor(c.R, c.G, c.B, c.A)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

type cursorPositionCallback func(*glfw.Window, float64, float64)

type framebufferSizeCallback func(*glfw.Window, int, int)

func (w Window) SetCursorPositionCallback(c cursorPositionCallback) {
	w.cursorPositionCallback = c
}

func (w Window) SetFramebufferSizeCallback(c framebufferSizeCallback) {
	w.framebufferSizeCallback = c
}

func (w Window) onCursorPositionChange(glw *glfw.Window, x, y float64) {
	w.Cursor = Point{x, y}
	if w.cursorPositionCallback != nil {
		w.cursorPositionCallback(glw, x, y)
	}
}

func (w Window) onFramebufferSizeChange(glw *glfw.Window, width, height int) {
	w.FrameSize = Dim{width, height}
	ww, wh := glw.GetSize()
	w.Size = Dim{ww, wh}
	gl.Viewport(w.Pos.X, w.Pos.Y, width, height)

	if w.framebufferSizeCallback != nil {
		w.framebufferSizeCallback(glw, width, height)
	}
}

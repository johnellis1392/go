package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	DEFAULT_SAMPLES = 0
	DEFAULT_WIDTH   = 640
	DEFAULT_HEIGHT  = 480
	DEFAULT_TITLE   = "Test Title"
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

	cursorPositionCallback     cursorPositionCallback
	framebufferSizeCallback    framebufferSizeCallback
	makeContextCurrentCallback makeContextCurrentCallback
}

func CreateWindow() (*Window, error) {
	title := DEFAULT_TITLE
	w, h, s := DEFAULT_WIDTH, DEFAULT_HEIGHT, DEFAULT_SAMPLES
	var err error

	glfw.WindowHint(glfw.Hint(glfw.Samples), s)

	glwindow, err := glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		return nil, err
	}

	glwindow.MakeContextCurrent()
	if err = gl.Init(); err != nil {
		return nil, err
	}

	const cursorX, cursorY = 200, 200
	window := &Window{
		w:      glwindow,
		Size:   Dim{w, h},
		Cursor: Point{cursorX, cursorY},
	}

	glwindow.SetFramebufferSizeCallback(window.onFramebufferSizeChange)
	glwindow.SetCursorPosCallback(window.onCursorPositionChange)

	return window, nil
}

func (w *Window) PollEvents() {
	glfw.PollEvents()
}

func (w *Window) Clear(c Color) {
	gl.ClearColor(c.R, c.G, c.B, c.A)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

type cursorPositionCallback func(*glfw.Window, float64, float64)

type framebufferSizeCallback func(*glfw.Window, int, int)

type makeContextCurrentCallback func()

func (w *Window) SetCursorPositionCallback(c cursorPositionCallback) {
	w.cursorPositionCallback = c
}

func (w *Window) SetFramebufferSizeCallback(c framebufferSizeCallback) {
	w.framebufferSizeCallback = c
}

func (w *Window) SetMakeContextCurrentCallback(c makeContextCurrentCallback) {
	w.makeContextCurrentCallback = c
}

func (w *Window) resize(x, y float64, width, height int) {
	w.Pos = Point{float64(x), float64(y)}
	w.Size = Dim{width, height}
	w.FrameSize = Dim{width, height}
}

func (w *Window) Viewport(x, y float64, width, height int) {
	w.resize(x, y, width, height)
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}

func (w *Window) onCursorPositionChange(glwindow *glfw.Window, x, y float64) {
	w.Cursor = Point{x, y}
	if w.cursorPositionCallback != nil {
		w.cursorPositionCallback(glwindow, x, y)
	}
}

func (w *Window) onFramebufferSizeChange(glwindow *glfw.Window, width, height int) {
	w.Viewport(w.Pos.X, w.Pos.Y, width, height)
	if w.framebufferSizeCallback != nil {
		w.framebufferSizeCallback(glwindow, width, height)
	}
}

func (w *Window) onMakeContextCurrent() error {
	if err := gl.Init(); err != nil {
		return err
	}
	if w.makeContextCurrentCallback != nil {
		w.makeContextCurrentCallback()
	}
	return nil
}

func (w *Window) MakeContextCurrent() error {
	w.w.MakeContextCurrent()
	return w.onMakeContextCurrent()
}

func (w *Window) ShouldClose() bool {
	return w.w.ShouldClose()
}

func (w *Window) SwapBuffers() {
	w.w.SwapBuffers()
}

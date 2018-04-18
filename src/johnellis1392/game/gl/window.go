package gl

import (
	"johnellis1392/game/math"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	defaultSamples = 0
	defaultWidth   = 640
	defaultHeight  = 480
	defaultTitle   = ""
)

// Window is an abstraction on top of GLFW's native Window object that provides
// some nice abstractions for us to work with.
type Window struct {
	Pos       math.Point
	Size      math.Dim
	FrameSize math.Dim
	Cursor    math.Point
	w         *glfw.Window

	cursorPositionCallback     cursorPositionCallback
	framebufferSizeCallback    framebufferSizeCallback
	makeContextCurrentCallback makeContextCurrentCallback
}

// Canvas is an abstraction of the GLFW Window's Framebuffer object.
type Canvas struct{}

// CreateWindow CreateWindow
func CreateWindow() (*Window, error) {
	title := defaultTitle
	w, h, s := defaultWidth, defaultHeight, defaultSamples
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
		Size:   math.Dim{W: w, H: h},
		Cursor: math.Point{X: cursorX, Y: cursorY},
	}

	glwindow.SetFramebufferSizeCallback(window.onFramebufferSizeChange)
	glwindow.SetCursorPosCallback(window.onCursorPositionChange)

	return window, nil
}

// Context creates a new GL Context from the given Window object
// and returns it.
func (w *Window) Context() Context {
	return context{
		// program: w.program,
		// buffers:  []uint32{},
		// attribs:  []uint32{},
		// uniforms: []uint32{},
	}
}

// PollEvents calls GLFW's PollEvents function, which starts processing
// window events.
func (w *Window) PollEvents() {
	glfw.PollEvents()
}

// Clear clears the screen with the supplied color.
func (w *Window) Clear(c math.Color) {
	gl.ClearColor(c.R, c.G, c.B, c.A)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

type cursorPositionCallback func(*glfw.Window, float64, float64)

type framebufferSizeCallback func(*glfw.Window, int, int)

type makeContextCurrentCallback func()

// SetCursorPositionCallback adds a callback to the window which fires
// whenever the mouse cursor changes position.
func (w *Window) SetCursorPositionCallback(c cursorPositionCallback) {
	w.cursorPositionCallback = c
}

// SetFramebufferSizeCallback adds a callback to the window which fires
// when the Window's size changes, such as by manual resizing. Resize
// operations cause GL's Framebuffer object to adjust to the window's new
// size.
//
// The Framebuffer is the viewport in the Window which GL renders
// to.
func (w *Window) SetFramebufferSizeCallback(c framebufferSizeCallback) {
	w.framebufferSizeCallback = c
}

// SetMakeContextCurrentCallback adds a callback to the window which fires
// when the GLFW window's MakeContextCurrent function is called. MakeContextCurrent
// reassigns the active GL context that is associated with the window's
// Framebuffer, meaning GL will unbind the current context so that clients
// will need to re-bind a new context when this operation occurs.
func (w *Window) SetMakeContextCurrentCallback(c makeContextCurrentCallback) {
	w.makeContextCurrentCallback = c
}

func (w *Window) resize(x, y float64, width, height int) {
	w.Pos = math.Point{X: float64(x), Y: float64(y)}
	w.Size = math.Dim{W: width, H: height}
	w.FrameSize = math.Dim{W: width, H: height}
}

// Viewport reassigns the current viewport frame for the GL context.
func (w *Window) Viewport(x, y float64, width, height int) {
	w.resize(x, y, width, height)
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}

func (w *Window) onCursorPositionChange(glwindow *glfw.Window, x, y float64) {
	w.Cursor = math.Point{X: x, Y: y}
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

// MakeContextCurrent activates the current window's GL context.
func (w *Window) MakeContextCurrent() error {
	w.w.MakeContextCurrent()
	return w.onMakeContextCurrent()
}

// ShouldClose returns a boolean indicating whether or not the GLFW window
// object has been closed.
func (w *Window) ShouldClose() bool {
	return w.w.ShouldClose()
}

// SwapBuffers swaps the current Framebuffer's render buffers. This is the
// operation for performing double-buffered rendering.
func (w *Window) SwapBuffers() {
	w.w.SwapBuffers()
}

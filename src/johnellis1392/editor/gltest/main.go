package main

// "github.com/go-gl/gl/v4.1-core/gl"
// "github.com/go-gl/glfw/v3.2/glfw"

import (
	"fmt"
	"johnellis1392/game/gl"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 300
	height = 300
	title  = "Example App"
)

const textboxVshader = `
#version 410
in vec3 vp;
void main(void) {
	gl_Position = vec4(vp, 1.0);
}
` + "\x00"

const textboxFshader = `
#version 410
out vec4 frag_colour;
void main(void) {
	frag_color = vec4(1, 1, 1, 1);
}
` + "\x00"

const divVshader = `
#version 410
in vec3 vp;
void main(void) {
	gl_Position = vec4(vp, 1.0);
}
` + "\x00"

const divFshader = `
#version 410
out vec4 frag_colour;
void main(void) {
	frag_color = vec4(1, 1, 1, 1);
}
` + "\x00"

type Window struct {
	w *glfw.Window
}

func (w *Window) ShouldClose() bool {
	if w == nil || w.w == nil {
		return true
	}
	return w.w.ShouldClose()
}

func (w *Window) Init() (err error) {
	if err = gl.Init(); err != nil {
		return
	}

	return
}

func CreateWindow() (w *Window, err error) {
	if err = glfw.Init(); err != nil {
		return
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var window Window
	if window.w, err = glfw.CreateWindow(width, height, title, nil, nil); err != nil {
		return
	}

	window.w.MakeContextCurrent()
	if err = window.Init(); err != nil {
		return
	}

	w = &window
	return
}

type GLContext interface {
	Bind(Buffer)
}

type createContext struct{}

var _ GLContext = (*createContext)(nil)

func (ctx *createContext) Bind(b Buffer) {}

type updateContext struct{}

var _ GLContext = (*updateContext)(nil)

func (ctx *updateContext) Bind(b Buffer) {}

type Buffer struct {
	id uint32
}

type Component interface {
	Init(GLContext)
	Draw(GLContext)
}

type div struct{}

func (d *div) Init(ctx GLContext) {}
func (d *div) Draw(ctx GLContext) {}

var _ Component = (*div)(nil)

type textbox struct{}

func (d *textbox) Init(ctx GLContext) {}
func (d *textbox) Draw(ctx GLContext) {}

var _ Component = (*textbox)(nil)

func main() {
	fmt.Println("Running...")
	w, err := CreateWindow()
	if err != nil {
		panic(err)
	}

	gos := []Component{
		&div{},
		&textbox{},
	}

	var ctx GLContext
	ctx = &createContext{}
	for _, o := range gos {
		o.Init(ctx)
	}

	const fps = 10
	for !w.ShouldClose() {
		t := time.Now()

		ctx = &updateContext{}
		for _, o := range gos {
			o.Draw(ctx)
		}

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

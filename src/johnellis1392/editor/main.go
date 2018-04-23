package main

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	defaultWidth  uint = 640
	defaultHeight uint = 480
	defaultTitle       = "Test App"
)

type cursor struct {
	X, Y uint
}

type window struct {
	w *glfw.Window
}

func newWindow(c *config) (*window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glwin, err := glfw.CreateWindow(int(c.Width), int(c.Height), c.Title, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := gl.Init(); err != nil {
		return nil, err
	}

	return &window{
		w: glwin,
	}, nil
}

type buffer struct {
	// viewport
}

type app struct {
	config  *config
	window  *window
	buffers []*buffer
}

func (a *app) run() {
	for {
		// TODO: Start Event Loop
		// TODO: Handle Rendering
	}
}

func newApp(c *config) (*app, error) {
	w, err := newWindow(c)
	if err != nil {
		return nil, err
	}

	return &app{
		// cursor: c,
		window: w,
	}, nil
}

type config struct {
	Width  uint
	Height uint
	Title  string
}

func getenvOrElse(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func envConfig() *config {
	return &config{
		// Width:  getenvOrElse("WIDTH", defaultWidth),
		// Height: getenvOrElse("HEIGHT", defaultHeight),
		Width:  defaultWidth,
		Height: defaultHeight,
		Title:  defaultTitle,
	}
}

// TODO: Get filename from cli
func main() {
	// fmt.Println("Hello, World!")
	// c := envConfig()
	// app, err := newApp(c)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// app.run()

	fmt.Println("Testing GoCUI")
	GoCUITest()
}

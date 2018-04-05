package main

// This is an example usage of Faiface's Pixel and PixelGL
// UI and window management packages. More information on
// these utilities is available at their github page here:
// https://github.com/faiface/pixel
import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Function for running pixel code; creates
// a new window and updates it until it is
// closed. The pixelgl.Run function will
// handle the rest of the startup / shutdown
// process.
func run() {
	var err error

	config := pixelgl.WindowConfig{
		Title:  "Hello, World!",          // Window Title
		Bounds: pixel.R(0, 0, 1024, 768), // Bounding Box
		VSync:  true,                     // Enable 60hz refresh rate
	}

	win, err := pixelgl.NewWindow(config)
	if err != nil {
		panic(err)
	}

	// Clear window with specified color
	win.Clear(colornames.Skyblue)

	// While the window isn't closed, draw window UI.
	for !win.Closed() {
		// Draw window frame; This doesn't clear the previous state,
		// it just draws over what was previously there.
		win.Update()
	}
}

// Set pixelgl window manager running. This will
// set up the pixel gl context we need and handle
// the initialization and shutdown of all the various
// GL and window components we need to run the GUI,
// leaving us the customize the UX.
func main() {
	pixelgl.Run(run)
}

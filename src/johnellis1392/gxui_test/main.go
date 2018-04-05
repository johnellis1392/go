package main

// This library is a golang wrapper for OpenGL's full
// libraries. Every version of OpenGL is available here.
// "github.com/go-gl/gl"
//
// This is the library used to generate the Go->C
// OpenGL bindings used by github.com/go-gl/gl above.
// "github.com/go-gl/glow"

import (
	"time"

	"github.com/google/gxui"
	gldriver "github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/gxfont"
	"github.com/google/gxui/math"
	"github.com/google/gxui/samples/flags"

	goxjs "github.com/goxjs/gl"
	goxjsglfw "github.com/goxjs/glfw"
)

func appMain(driver gxui.Driver) {
	theme := flags.CreateTheme(driver)

	font, err := driver.CreateFont(gxfont.Default, 75)
	if err != nil {
		panic(err)
	}

	window := theme.CreateWindow(380, 100, "Hi")
	window.SetBackgroundBrush(gxui.CreateBrush(gxui.Gray50))

	label := theme.CreateLabel()
	label.SetFont(font)
	label.SetText("Hello world")

	window.AddChild(label)

	ticker := time.NewTicker(time.Millisecond * 30)
	go func() {
		phase := float32(0)
		for _ = range ticker.C {
			c := gxui.Color{
				R: 0.75 + 0.25*math.Cosf((phase+0.000)*math.TwoPi),
				G: 0.75 + 0.25*math.Cosf((phase+0.333)*math.TwoPi),
				B: 0.75 + 0.25*math.Cosf((phase+0.666)*math.TwoPi),
				A: 0.50 + 0.50*math.Cosf(phase*10),
			}
			phase += 0.01
			driver.Call(func() {
				label.SetColor(c)
			})
		}
	}()

	window.OnClose(ticker.Stop)
	window.OnClose(driver.Terminate)
}

func gxuiExample() {
	gldriver.StartDriver(appMain)
}

func loop() {
	// for {
	goxjsglfw.WaitEvents()
	// }
}

func goxjsExample() {
	var err error

	if err = goxjsglfw.Init(goxjs.ContextWatcher); err != nil {
		panic(err)
	}
	defer goxjsglfw.Terminate()

	goxjs.ClearColor(1.0, 0.0, 0.0, 0.5)
	goxjs.Clear(goxjs.COLOR_BUFFER_BIT)

	loop()

	// This is an example of how to use your own custom
	// ContextWatcher element:
	// window.MakeContextCurrent()
	// goxjs.ContextWatcher.OnMakeCurrent(nil)
	//
	// goxjsglfw.DetachCurrentContext()
	// goxjs.ContextWatcher.OnDetach()
}

func main() {
	// gxuiExample()
	goxjsExample()
}

// import (
// 	"fmt"
// 	"syscall"
// 	"unsafe"
// )
//
// type (
// 	HANDLE uintptr
// 	HWND   HANDLE
// )
//
// // Note: This may all be windows stuff
// var (
// 	// Library
// 	libuser32 uintptr
//
// 	// Functions
// 	setParent uintptr
// )
//
// func MustLoadLibrary(libname string) uintptr {
// 	return uintptr(0)
// }
//
// func MustGetProcAddress(lib uintptr, fn string) uintptr {
// 	return uintptr(0)
// }
//
// func init() {
// 	is64bit := unsafe.Sizeof(uintptr(0)) == 8
// 	libuser32 = MustLoadLibrary("user32.dll")
// 	setParent = MustGetProcAddress(libuser32, "SetParent")
// }
//
// // SetParent - Copy of function from github.com/lxn/win
// func SetParent(hWnd HWND, parentHWnd HWND) HWND {
// 	ret, _, _ := syscall.Syscall(setParent, 2, uintptr(hWnd), uintptr(parentHWnd), 0)
// 	return HWND(ret)
// }
//
// func main() {
// 	fmt.Println("Hello, World!")
//
// 	// var handle1 uintptr
//
// }

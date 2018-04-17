package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
)

// #include <stdlib.h>
import "C"

func CreateProgram() Program {
	pid := gl.CreateProgram()
	p := Program{pid, []Shader{}}
	return p
}

func CreateShader(shaderType uint32, filename string) Shader {
	sid := gl.CreateShader(shaderType)
	shader := Shader{sid, filename}
	return shader
}

type Program struct {
	ID      uint32
	Shaders []Shader
}

func (p Program) Use() {
	gl.UseProgram(p.ID)
}

func (p Program) Validate() error {
	if p.ID == 0 {
		return fmt.Errorf("error: gl returned invalid null pointer when creating program: Program.ID = 0")
	}

	gl.ValidateProgram(p.ID)

	var status int32
	gl.GetProgramiv(p.ID, gl.VALIDATE_STATUS, &status)
	if status != gl.TRUE {
		return fmt.Errorf("%s", p.InfoLog())
	}

	return nil
}

func (p Program) AddShader(shaderType uint32, src string) {
	s := CreateShader(shaderType, src)
	p.Shaders = append(p.Shaders, s)
}

func (p Program) InfoLog() string {
	var logLength int32
	gl.GetProgramiv(p.ID, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetProgramInfoLog(p.ID, logLength, nil, &logBuffer[0])
	return GoString(&logBuffer[0])
}

func (p Program) Delete() {
	for _, s := range p.Shaders {
		s.Delete()
	}
	gl.DeleteProgram(p.ID)
}

func (p Program) link() error {
	gl.LinkProgram(p.ID)

	var programi int32
	gl.GetProgramiv(p.ID, gl.LINK_STATUS, &programi)
	if programi == 0 {
		// defer gl.DeleteProgram(p)
		defer p.Delete()
		return fmt.Errorf("failed to link program: %s", p.InfoLog())
	}

	return nil
}

func (p Program) Compile() error {
	var err error
	for _, s := range p.Shaders {
		if err = s.Compile(); err != nil {
			defer p.Delete()
			return err
		}
	}
	return p.link()
}

type Shader struct {
	ID   uint32
	file string
}

func (s Shader) Compile() error {
	// Load Shader Source from File
	var src []byte
	var err error
	if src, err = ioutil.ReadFile(s.file); err != nil {
		return err
	}

	glsrc, free := gl.Strs(string(src) + "\x00")
	gl.ShaderSource(s.ID, 1, glsrc, nil)
	free()

	// Compile Shader
	gl.CompileShader(s.ID)

	// Check Shader Compile Status
	var status int32
	gl.GetShaderiv(s.ID, gl.COMPILE_STATUS, &status)
	fmt.Printf("Received Shader Compile Status: %v, %q", status, status)
	if status == 0 {
		return fmt.Errorf("%s", s.InfoLog())
	}

	return nil
}

func (s Shader) InfoLog() string {
	var logLength int32
	gl.GetShaderiv(s.ID, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetShaderInfoLog(s.ID, logLength, nil, &logBuffer[0])
	return GoString(&logBuffer[0])
}

func (s Shader) Delete() {
	gl.DeleteShader(s.ID)
}

type GLWindow struct{}

type GLCanvas struct{}

// Conversion Functions

func Str(str string) *uint8 {
	if !strings.HasSuffix(str, "\x00") {
		panic("str argument missing null terminator: " + str)
	}
	header := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return (*uint8)(unsafe.Pointer(header.Data))
}

func GoString(cstr *uint8) string {
	return C.GoString((*C.char)(unsafe.Pointer(cstr)))
}

func GetString(n uint32) string {
	return GoString(gl.GetString(n))
}

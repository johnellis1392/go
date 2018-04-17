package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
)

func CreateProgram() Program {
	pid := gl.CreateProgram()
	p := Program{pid, []Shader{}}
	return p
}

func CreateShader(shaderType uint32, src string) (Shader, error) {
	sid := gl.CreateShader(shaderType)

	f, err := os.Open(src)
	if err != nil {
		return Shader{}, err
	}

	var s string
	if s, err := ioutil.ReadFile(src); err != nil {
		return Shader{}, err
	}

	glsrc, free := gl.Strs(src + "\x00")
	gl.ShaderSource(sid, 1, glsrc, nil)
	free()

	shader := Shader{sid, src}
	return shader, nil
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
		return fmt.Errorf("failed to create program: Program.ID = %v", p.ID)
	}

	gl.ValidateProgram(p.ID)

	var status int32
	gl.GetProgramiv(p.ID, gl.VALIDATE_STATUS, &status)
	if status != gl.TRUE {
		return fmt.Errorf("%s", p.InfoLog())
	}

	return nil
}

func (p Program) AddShader(shaderType uint32, src string) error {
	s, err := CreateShader(shaderType, src)
	if err != nil {
		return err
	}

	p.Shaders = append(p.Shaders, s)
	return nil
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

type Shader struct {
	ID   uint32
	file string
}

func (s Shader) Compile() error {
	gl.CompileShader(s.ID)

	var status int32
	gl.GetShaderiv(s.ID, gl.COMPILE_STATUS, &status)
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

type GLWindow struct{}

type GLCanvas struct{}

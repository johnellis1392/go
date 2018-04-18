package gl

import (
	"fmt"
	"io/ioutil"

	"github.com/go-gl/gl/v2.1/gl"
)

// Enum is an abstraction over GL's enormous collection of Enum values.
type Enum uint32

// Implementation values for the above Enum type.
const (
	Vendor                 Enum = gl.VENDOR
	Renderer                    = gl.RENDERER
	Version                     = gl.VERSION
	ShadingLanguageVersion      = gl.SHADING_LANGUAGE_VERSION
)

// Unbox converts a GL Enum value back into the expected int32 type.
// Golang doesn't automatically cast alias types, so this function
// fluentizes the type casting process for these values.
func (e Enum) Unbox() uint32 {
	return uint32(e)
}

// String returns a readable string representation of the given enum value.
func (e Enum) String() string {
	switch e {
	case Vendor:
		return "Vendor"
	case Renderer:
		return "Renderer"
	case Version:
		return "Version"
	case ShadingLanguageVersion:
		return "ShadingLanguageVersion"
	default:
		return "[unknown]"
	}
}

// CreateProgram queries OpenGL to create a new program for us to work with,
// and returns a new Program object for us to work with.
func CreateProgram() Program {
	pid := gl.CreateProgram()
	p := Program{pid, []Shader{}}
	return p
}

// CreateShader queries OpenGL to create a new shader object and returns a
// new Shader struct for us to work with.
func CreateShader(shaderType uint32, filename string) Shader {
	sid := gl.CreateShader(shaderType)
	shader := Shader{sid, filename}
	return shader
}

// Program is an abstraction of GL's program and shader interfaces.
type Program struct {
	ID      uint32
	Shaders []Shader
}

// Use sets the current program as the active program in GL.
func (p *Program) Use() {
	gl.UseProgram(p.ID)
}

// AddShader creates a new shader from the given source file, and saves it to
// to the program's shader list.
func (p *Program) AddShader(shaderType uint32, src string) {
	s := CreateShader(shaderType, src)
	p.Shaders = append(p.Shaders, s)
}

// InfoLog retrieves the GL program info log for the given program.
func (p *Program) InfoLog() string {
	var logLength int32
	gl.GetProgramiv(p.ID, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetProgramInfoLog(p.ID, logLength, nil, &logBuffer[0])
	return GoString(&logBuffer[0])
}

// Delete deletes the current program and all its associated shaders.
func (p *Program) Delete() {
	for _, s := range p.Shaders {
		s.Delete()
	}
	gl.DeleteProgram(p.ID)
}

// validate validates the current program. GL can fail to create, compile, or
// link programs, and other errors can potentially occur inside the GPU, so this
// function queries GL to validate the given program.
func (p *Program) validate() error {
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

// link links the given program in GL.
func (p *Program) link() error {
	gl.LinkProgram(p.ID)

	var programi int32
	gl.GetProgramiv(p.ID, gl.LINK_STATUS, &programi)
	if programi == 0 {
		defer p.Delete()
		return fmt.Errorf("failed to link program: %s", p.InfoLog())
	}

	return nil
}

// Compile compiles the given program and all of its associated shaders in GL.
// Each shader that has been associated with this program must be compiled and
// attached to the program in order to run correctly. As well, the program must
// be linked in order to run properly, and validated to verify no errors occurred.
func (p *Program) Compile() error {
	var err error
	fmt.Println("Compiling:", p.Shaders)
	for _, s := range p.Shaders {
		fmt.Println("Compiling Shader:", s.ID, s.file)
		if err = s.Compile(); err != nil {
			defer p.Delete()
			return err
		}
		gl.AttachShader(p.ID, s.ID)
	}

	if err = p.link(); err != nil {
		return err
	}

	if err = p.validate(); err != nil {
		return err
	}

	return nil
}

// Shader is an abstraction over GL's Shader objects.
type Shader struct {
	ID   uint32
	file string
}

// Compile reads the shader's source from the file specified by s.file, loads
// the shader's source into GL, and compiles the shader. If an error occurs
// either during this process or during GL's compilation process on the GPU,
// the error is returned to the caller.
func (s *Shader) Compile() error {
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

// InfoLog gets the Shader info log for the given shader.
func (s *Shader) InfoLog() string {
	var logLength int32
	gl.GetShaderiv(s.ID, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetShaderInfoLog(s.ID, logLength, nil, &logBuffer[0])
	return GoString(&logBuffer[0])
}

// Delete deletes the given shader.
func (s *Shader) Delete() {
	gl.DeleteShader(s.ID)
}

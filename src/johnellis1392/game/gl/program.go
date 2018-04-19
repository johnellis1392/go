package gl

import (
	"fmt"
	"io/ioutil"

	"github.com/go-gl/gl/v2.1/gl"
)

// Program is an abstraction of GL's program and shader interfaces.
type Program struct {
	ID             Ptr
	VertexShader   *Shader
	FragmentShader *Shader
}

// Use sets the current program as the active program in GL.
func (p *Program) Use() {
	gl.UseProgram(p.ID.Unbox())
}

// AddShader creates a new shader from the given source file, and saves it to
// to the program's shader list.
func (p *Program) AddShader(shaderType Enum, src string) {
	var s Shader
	switch shaderType {
	case VertexShader:
		s = CreateShader(shaderType, src)
		p.VertexShader = &s
	case FragmentShader:
		s = CreateShader(shaderType, src)
		p.FragmentShader = &s
	default:
		// Invalid Shader Type; Pass
		break
	}
}

// InfoLog retrieves the GL program info log for the given program.
func (p *Program) InfoLog() string {
	var logLength int32
	gl.GetProgramiv(p.ID.Unbox(), InfoLogLength.Unbox(), &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetProgramInfoLog(p.ID.Unbox(), logLength, nil, &logBuffer[0])
	return GoString(&logBuffer[0])
}

// Delete deletes the current program and all its associated shaders.
func (p *Program) Delete() {
	for _, s := range p.Shaders() {
		if s != nil {
			s.Delete()
		}
	}

	gl.DeleteProgram(p.ID.Unbox())
}

// Shaders collects the given Program's Shaders into a single list and returns it.
func (p *Program) Shaders() []*Shader {
	return []*Shader{
		p.VertexShader,
		p.FragmentShader,
	}
}

// validate validates the current program. GL can fail to create, compile, or
// link programs, and other errors can potentially occur inside the GPU, so this
// function queries GL to validate the given program.
func (p *Program) validate() error {
	if p.ID == 0 {
		return fmt.Errorf("error: gl returned invalid null pointer when creating program: Program.ID = 0")
	}

	gl.ValidateProgram(p.ID.Unbox())

	var status int32
	gl.GetProgramiv(p.ID.Unbox(), ValidateStatus.Unbox(), &status)
	if status != gl.TRUE {
		return fmt.Errorf("%s", p.InfoLog())
	}

	return nil
}

// link links the given program in GL.
func (p *Program) link() error {
	gl.LinkProgram(p.ID.Unbox())

	var programi int32
	gl.GetProgramiv(p.ID.Unbox(), LinkStatus.Unbox(), &programi)
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
	shaders := p.Shaders()
	fmt.Println("Compiling:", shaders)
	for _, s := range shaders {
		fmt.Println("Compiling Shader:", s.ID, s.file)
		if err = s.Compile(); err != nil {
			defer p.Delete()
			return err
		}
		gl.AttachShader(p.ID.Unbox(), s.ID.Unbox())
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
	ID   Ptr
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
	gl.ShaderSource(s.ID.Unbox(), 1, glsrc, nil)
	free()

	// Compile Shader
	gl.CompileShader(s.ID.Unbox())

	// Check Shader Compile Status
	var status int32
	gl.GetShaderiv(s.ID.Unbox(), CompileStatus.Unbox(), &status)
	fmt.Printf("Received Shader Compile Status: %v, %q", status, status)
	if status == 0 {
		return fmt.Errorf("%s", s.InfoLog())
	}

	return nil
}

// InfoLog gets the Shader info log for the given shader.
func (s *Shader) InfoLog() string {
	var logLength int32
	gl.GetShaderiv(s.ID.Unbox(), InfoLogLength.Unbox(), &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetShaderInfoLog(s.ID.Unbox(), logLength, nil, &logBuffer[0])
	return GoString(&logBuffer[0])
}

// Delete deletes the given shader.
func (s *Shader) Delete() {
	gl.DeleteShader(s.ID.Unbox())
}

// ///////////////// //
// Utility Functions //
// ///////////////// //

// CreateProgram queries OpenGL to create a new program for us to work with,
// and returns a new Program object for us to work with.
func CreateProgram() Program {
	pid := gl.CreateProgram()
	p := Program{Ptr(pid), nil, nil}
	return p
}

// CreateShader queries OpenGL to create a new shader object and returns a
// new Shader struct for us to work with.
func CreateShader(shaderType Enum, filename string) Shader {
	sid := gl.CreateShader(shaderType.Unbox())
	shader := Shader{Ptr(sid), filename}
	return shader
}

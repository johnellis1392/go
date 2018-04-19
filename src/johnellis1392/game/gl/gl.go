package gl

import (
	"johnellis1392/game/math"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Init locks the current OS thread (mandatory for GL to work; GL must run
// on main application thread), and initializes GLFW.
func Init() error {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		return err
	}

	return nil
}

// Terminate unlocks the current OS thread and deletes the current GLFW context.
func Terminate() {
	glfw.Terminate()
	runtime.UnlockOSThread()
}

// Ptr represents an Identifier for a GL Object.
type Ptr uint32

// Unbox converts the pointer object into a valid uint32 value.
func (p Ptr) Unbox() uint32 {
	return uint32(p)
}

// Buffer is a GL Buffer object.
type Buffer struct {
	ID   Ptr
	Type Enum
}

// Bind binds the current buffer in GL.
func (b Buffer) Bind() {
	gl.BindBuffer(b.Type.Unbox(), b.ID.Unbox())
}

// Set loads the given list of vertices into the current buffer.
func (b Buffer) Set(vs []math.Vec3) {
	gl.BufferData(b.Type.Unbox(), int(len(vs)), gl.Ptr(&vs[0]), StaticDraw.Unbox())
}

func (b Buffer) delete() {
	// TODO
}

// Attrib is a GL Vertex Attribute Object.
type Attrib struct {
	ID   Ptr
	Name string
}

// Enable enables the VertexAttributeArray in GL.
func (a Attrib) Enable() {
	gl.EnableVertexAttribArray(a.ID.Unbox())
}

// Bind associates the current Vertex Attribute object to the buffer currently
// bound to gl.ARRAY_BUFFER
//
// From the Mozilla Developer Network's WebGL Documentation:
// (available here: https://developer.mozilla.org/en-US/docs/Web/API/WebGLRenderingContext/vertexAttribPointer)
// "[gl.VertexAttribPointer(...)] binds the buffer currently bound to gl.ARRAY_BUFFER
// to a generic vertex attribute of the of the current vertex buffer object and
// specifies its layout."
func (a Attrib) Bind(size int) {
	gl.VertexAttribPointer(a.ID.Unbox(), int32(size), Float.Unbox(), false, 0, nil)
}

func (a Attrib) delete() {
	// TODO
}

// Uniform represents a GL Uniform variable.
type Uniform struct {
	ID   Ptr
	Name string
}

// Set loads the given matrix's data into the given Uniform.
func (u Uniform) Set(m mgl32.Mat4) {
	gl.UniformMatrix4fv(int32(u.ID.Unbox()), int32(len(m)/(4*4)), false, &m[0])
}

func (u Uniform) delete() {
	// TODO
}

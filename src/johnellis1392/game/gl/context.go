package gl

import "github.com/go-gl/gl/v2.1/gl"

// Context is a wrapper around OpenGL's native functionality that consolidates
// the state of the running application and provides some extra useful features
// for graphical development.
type Context interface {
	Draw(int)
	CreateBuffer(Enum) Buffer
	GetAttrib(string) Attrib
	GetUniform(string) Uniform
}

type context struct {
	program  Program
	buffers  []Buffer
	attribs  map[string]Attrib
	uniforms map[string]Uniform
}

func (c context) Draw(n int) {
	gl.DrawArrays(Triangles.Unbox(), int32(0), int32(n))
}

func (c context) CreateBuffer(typ Enum) Buffer {
	var bufptr uint32
	gl.GenBuffers(1, &bufptr)

	buffer := Buffer{Ptr(bufptr), typ}
	c.buffers = append(c.buffers, buffer)
	return buffer
}

func (c context) GetAttrib(name string) Attrib {
	if a, ok := c.attribs[name]; ok {
		return a
	}

	s := AsCString(name)
	attrid := gl.GetAttribLocation(c.program.ID.Unbox(), s.val)
	a := Attrib{Ptr(attrid), name}

	c.attribs[name] = a
	return a
}

func (c context) GetUniform(name string) Uniform {
	if u, ok := c.uniforms[name]; ok {
		return u
	}

	s := AsCString(name)
	// s := Str(name + "\x00")
	uniformptr := gl.GetUniformLocation(c.program.ID.Unbox(), s.val)
	u := Uniform{Ptr(uniformptr), name}

	c.uniforms[name] = u
	return u
}

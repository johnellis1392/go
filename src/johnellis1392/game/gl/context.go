package gl

// Context Context
type Context interface {
	// Draw()
	// CreateBuffer() uint32
	// GetAttrib(name string) uint32
	// GetUniform() uint32
}

type context struct {
	program Program
	// buffers  []uint32
	// attribs  []uint32
	// uniforms []uint32
}

// func (c context) Draw() {
// 	// TODO: Drawing Logic Goes Here
// }

// func (c context) CreateBuffer() uint32 {
// 	var buffer uint32
// 	gl.GenBuffers(1, &buffer)
// 	c.buffers = append(c.buffers, buffer)
// 	return buffer
// }

// func (c context) GetAttrib(name string) uint32 {
// 	return uint32(gl.GetAttribLocation(c.program.ID, gl.Str(name+"\x00")))
// }

// func (c context) GetUniform() uint32 {
//
// }

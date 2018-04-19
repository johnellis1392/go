package gl

import (
	"github.com/go-gl/gl/v2.1/gl"
)

// Enum is an abstraction over GL's enormous collection of Enum values.
type Enum uint32

// Implementation values for the above Enum type.
const (
	ArrayBuffer            Enum = gl.ARRAY_BUFFER
	CompileStatus          Enum = gl.COMPILE_STATUS
	Float                  Enum = gl.FLOAT
	FragmentShader         Enum = gl.FRAGMENT_SHADER
	InfoLogLength          Enum = gl.INFO_LOG_LENGTH
	LinkStatus             Enum = gl.LINK_STATUS
	Renderer               Enum = gl.RENDERER
	ShadingLanguageVersion Enum = gl.SHADING_LANGUAGE_VERSION
	StaticDraw             Enum = gl.STATIC_DRAW
	Triangles              Enum = gl.TRIANGLES
	ValidateStatus         Enum = gl.VALIDATE_STATUS
	Vendor                 Enum = gl.VENDOR
	Version                Enum = gl.VERSION
	VertexShader           Enum = gl.VERTEX_SHADER
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

package math

import "fmt"

type Color struct {
	R, G, B, A float32
}

type Point struct {
	X, Y float64
}

type Dim struct {
	W, H int
}

type Vec3 [3]float32

func (v Vec3) Dot(v2 Vec3) float32 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2]
}

func (v Vec3) Plus(v2 Vec3) Vec3 {
	return Vec3{
		v[0] + v2[0],
		v[1] + v2[1],
		v[2] + v2[2],
	}
}

type Mat4 [4][4]float32

func Mat4Id() Mat4 {
	return Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (m Mat4) Mul(m2 Mat4) Mat4 {
	var mres Mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			var s float32
			for k := 0; k < 4; k++ {
				s += m[i][k] * m2[k][j]
			}
			mres[i][j] = s
		}
	}
	return mres
}

func (m Mat4) Vec() [16]float32 {
	var res [16]float32
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			res[4*i+j] = m[i][j]
		}
	}
	return res
}

func (m Mat4) String() string {
	const format = `[
	[ %v, %v, %v, %v ],
	[ %v, %v, %v, %v ],
	[ %v, %v, %v, %v ],
	[ %v, %v, %v, %v ],
]
`
	var vs = make([]interface{}, 16)
	for i, f := range m.Vec() {
		vs[i] = f
	}
	return fmt.Sprintf(format, vs...)
}

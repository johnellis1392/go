package main

import "fmt"

type vec3 [3]float32

func (v vec3) dot(v2 vec3) float32 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2]
}

func (v vec3) plus(v2 vec3) vec3 {
	return vec3{
		v[0] + v2[0],
		v[1] + v2[1],
		v[2] + v2[2],
	}
}

type mat4 [4][4]float32

func mat4Id() mat4 {
	return mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (m mat4) mul(m2 mat4) mat4 {
	var mres mat4
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

func (m mat4) vec() [16]float32 {
	var res [16]float32
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			res[4*i+j] = m[i][j]
		}
	}
	return res
}

func (m mat4) String() string {
	const format = `[
	[ %v, %v, %v, %v ],
	[ %v, %v, %v, %v ],
	[ %v, %v, %v, %v ],
	[ %v, %v, %v, %v ],
]
`
	var vs = make([]interface{}, 16)
	for i, f := range m.vec() {
		vs[i] = f
	}
	return fmt.Sprintf(format, vs...)
}

//func test2() {
//	m1 := mat4{
//		{1.0, 0.0, 0.0},
//		{0.0, 1.0, 0.0},
//		{0.0, 0.0, 1.0},
//	}
//
//	m2 := mat4{
//		{1.0, 0.0, 0.0},
//		{0.0, 1.0, 0.0},
//		{0.0, 0.0, 1.0},
//	}
//
//	m3 := m1.mul(m2)
//	fmt.Printf("m1 * m2 = %v", m3)
//}

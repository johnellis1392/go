package main

import (
	"time"

	"github.com/go-gl/gl/v2.1/gl"
)

type GameState struct {
	objects []GameObject
	time    time.Time
}

type GameContext interface{}

type GameObject interface {
	Init() error
	Update(c GameContext, t time.Time)
	Render(c GLContext)
}

type GLPos struct {
	Pos, Rot, Scale mat4
}

type Cube struct {
	GLPos
	Vertices []vec3
}

type GLContext interface {
	Draw()
	CreateBuffer() uint32
	GetAttrib(name string) uint32
	// GetUniform() uint32
}

type context struct {
	program  Program
	buffers  []uint32
	attribs  []uint32
	uniforms []uint32
}

func (c context) Draw() {
	// TODO: Drawing Logic Goes Here
}

func (c context) CreateBuffer() uint32 {
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	c.buffers = append(c.buffers, buffer)
	return buffer
}

func (c context) GetAttrib(name string) uint32 {
	return uint32(gl.GetAttribLocation(c.program.ID, gl.Str(name+"\x00")))
}

// func (c context) GetUniform() uint32 {
//
// }

type stateFn func(context, GameObject) stateFn

type monster struct {
	pos, rot, scale mat4
	state           stateFn
}

func monsterAttack(c context, g GameObject) stateFn {
	return nil
}

func monsterIdle(c context, g GameObject) stateFn {
	return nil
}

func monsterInit(c context, g GameObject) stateFn {
	// switch m := g.(monster); {
	// case m.Senses(c.GetPlayer()):
	// 	return monsterAttack
	// default:
	// 	return monsterIdle
	// }
	return nil
}

func (m monster) Update(t time.Time) {
	// m.state = m.state(nil, m)
}

type player struct {
	pos, rot, scale mat4 // Position, Rotation & Scale
	state           stateFn
}

func playerInit(c context, p GameObject) stateFn {
	return nil
}

func (p player) Update(t time.Time) {
	if p.state == nil {
		// Player is Dead; Game Over
		// TODO: Signal End of Game
		return
	}

	// p.state = p.state(c, p)
}

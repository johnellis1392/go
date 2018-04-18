package main

import (
	"johnellis1392/game/gl"
	"johnellis1392/game/math"
	"time"
)

type GameObject interface {
	Init() error
	Update(t time.Time)
	Render(c gl.Context)
	Destroy() error
}

type GLPos struct {
	Pos, Rot, Scale math.Mat4
}

type stateFn func(gl.Context, GameObject) stateFn

type monster struct {
	pos, rot, scale math.Mat4
	state           stateFn
}

func monsterAttack(c gl.Context, g GameObject) stateFn {
	return nil
}

func monsterIdle(c gl.Context, g GameObject) stateFn {
	return nil
}

func monsterInit(c gl.Context, g GameObject) stateFn {
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
	pos, rot, scale math.Mat4 // Position, Rotation & Scale
	state           stateFn
}

func playerInit(c gl.Context, p GameObject) stateFn {
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

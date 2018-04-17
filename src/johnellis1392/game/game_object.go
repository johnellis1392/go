package main

import (
	"time"
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

type GLContext interface {
	Draw()
}

type context struct {
	// ...
}

func (c context) Draw() {
	// TODO: Drawing Logic Goes Here
}

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

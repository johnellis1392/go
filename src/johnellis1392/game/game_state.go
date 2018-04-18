package main

import "time"

type GameState struct {
	objects []GameObject
	time    time.Time
}

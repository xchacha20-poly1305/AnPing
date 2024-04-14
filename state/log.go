package state

import (
	"time"
)

type Logger interface {
	OnStart(address string, state *State)
	OnRecv(address string, state *State, t time.Duration)
	OnLost(address string, state *State, errMessage string)
	OnFinish(address string, state *State)
}

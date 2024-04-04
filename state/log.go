package state

import (
	"time"
)

type Logger interface {
	OnStart(address string)
	OnRecv(address string, t time.Duration)
	OnLost(address string, errMessage string)
	OnFinish(address string, probed, lost, succeed, min, avg, max, mdev uint64)
}

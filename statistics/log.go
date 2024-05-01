package statistics

import (
	"time"
)

type Logger interface {
	OnStart(address string, statistics StatisticsGetter)
	OnRecv(address string, statistics StatisticsGetter, t time.Duration)
	OnLost(address string, statistics StatisticsGetter, errMessage string, t time.Duration)
	OnFinish(address string, statistics StatisticsGetter)
}

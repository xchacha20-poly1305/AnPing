package statistics

import (
	"time"

	M "github.com/sagernet/sing/common/metadata"
)

type Logger interface {
	OnStart(address M.Socksaddr, statistics StatisticsGetter)
	OnRecv(address M.Socksaddr, statistics StatisticsGetter, t time.Duration)
	OnLost(address M.Socksaddr, statistics StatisticsGetter, errMessage string, t time.Duration)
	OnFinish(address M.Socksaddr, statistics StatisticsGetter)
}

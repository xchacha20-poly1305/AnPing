package state

type Logger interface {
	OnStart(address string)
	OnRecv(address string, t int)
	OnLost(address string, errMessage string)
	OnFinish(address string, probed, lost, succeed, min, avg, max, mdev uint64)
}

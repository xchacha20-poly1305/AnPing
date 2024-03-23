package anping

type Logger interface {
	OnStart(opts *Options)
	OnRecv(opts *Options, t int)
	OnLost(opts *Options, errMessage string)
	OnFinish(opts *Options)
}

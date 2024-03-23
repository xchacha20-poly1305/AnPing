package anping

var _ Logger = (*LoggerNotNil)(nil)

type LoggerNotNil struct {
	L Logger
}

func (l *LoggerNotNil) OnStart(opts *Options) {
	if l.L != nil {
		l.L.OnStart(opts)
	}
}

func (l *LoggerNotNil) OnRecv(opts *Options, t int) {
	if l.L != nil {
		l.L.OnRecv(opts, t)
	}
}

func (l *LoggerNotNil) OnLost(opts *Options, errMessage string) {
	if l.L != nil {
		l.L.OnLost(opts, errMessage)
	}
}

func (l *LoggerNotNil) OnFinish(opts *Options) {
	if l.L != nil {
		l.L.OnFinish(opts)
	}
}

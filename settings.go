package anping

import (
	"time"

	"github.com/sagernet/sing/common/control"
	M "github.com/sagernet/sing/common/metadata"
)

// Settings allows setting common settings.
type Settings interface {
	Logger() Logger
	SetLogger(logger Logger)

	Target() M.Socksaddr
	SetTarget(target M.Socksaddr)

	Interval() time.Duration
	SetInterval(interval time.Duration)
	Timeout() time.Duration
	SetTimeout(timeout time.Duration)

	Count() int64
	SetCount(count int64)

	Interface() string
	SetInterface(interfaceName string)
}

var _ Settings = (*DefaultSetting)(nil)

type DefaultSetting struct {
	logger          Logger
	target          M.Socksaddr
	interval        time.Duration
	timeout         time.Duration
	count           int64
	InterfaceFinder control.InterfaceFinder
	interfaceName   string
}

func NewSettings(finder control.InterfaceFinder) DefaultSetting {
	if finder == nil {
		finder = control.NewDefaultInterfaceFinder()
	}
	return DefaultSetting{
		InterfaceFinder: finder,
	}
}

func (d DefaultSetting) Logger() Logger {
	return d.logger
}

func (d *DefaultSetting) SetLogger(logger Logger) {
	d.logger = logger
}

func (d DefaultSetting) Target() M.Socksaddr {
	return d.target
}

func (d *DefaultSetting) SetTarget(target M.Socksaddr) {
	d.target = target
}

func (d DefaultSetting) Interval() time.Duration {
	return d.interval
}

func (d *DefaultSetting) SetInterval(interval time.Duration) {
	d.interval = interval
}

func (d DefaultSetting) Timeout() time.Duration {
	return d.timeout
}

func (d *DefaultSetting) SetTimeout(timeout time.Duration) {
	d.timeout = timeout
}

func (d DefaultSetting) Count() int64 {
	return d.count
}

func (d *DefaultSetting) SetCount(count int64) {
	d.count = count
}

func (d DefaultSetting) Interface() string {
	return d.interfaceName
}

func (d *DefaultSetting) SetInterface(interfaceName string) {
	d.interfaceName = interfaceName
}

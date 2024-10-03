package tcpping

import (
	"github.com/xchacha20-poly1305/anping"
)

var _ anping.Settings = (*TcpSettings)(nil)

type TcpSettings struct {
	anping.DefaultSetting
}

func NewSettings(commonSettings anping.DefaultSetting) TcpSettings {
	return TcpSettings{
		DefaultSetting: commonSettings,
	}
}

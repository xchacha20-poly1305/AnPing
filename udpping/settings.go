package udpping

import (
	"crypto/rand"

	"github.com/xchacha20-poly1305/anping"
)

const (
	PayloadLength = 40
)

var _ anping.Settings = (*UdpSettings)(nil)

type UdpSettings struct {
	anping.DefaultSetting
	payload func() []byte
}

func NewSettings(commonSettings anping.DefaultSetting, payload func() []byte) UdpSettings {
	if payload() == nil {
		randomPayload := make([]byte, PayloadLength)
		_, _ = rand.Read(randomPayload)
		payload = func() []byte {
			return randomPayload
		}
	}
	return UdpSettings{
		DefaultSetting: commonSettings,
		payload:        payload,
	}
}

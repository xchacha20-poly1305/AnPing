package icmpping

import (
	"crypto/rand"

	"github.com/xchacha20-poly1305/anping"
)

const (
	PayloadLength = 40
)

var _ anping.Settings = (*IcmpSettings)(nil)

type IcmpSettings struct {
	anping.DefaultSetting
	payload func() []byte
}

func NewSettings(commonSettings anping.DefaultSetting, payload func() []byte) IcmpSettings {
	if payload() == nil {
		randomPayload := make([]byte, PayloadLength)
		_, _ = rand.Read(randomPayload)
		payload = func() []byte {
			return randomPayload
		}
	}
	return IcmpSettings{
		DefaultSetting: commonSettings,
		payload:        payload,
	}
}

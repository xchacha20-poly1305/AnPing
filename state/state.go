package state

import (
	"math"

	"github.com/sagernet/sing/common/atomic"
)

type State struct {
	// probed is the time that probed.
	probed atomic.Uint64
	// succeed is the time that probing succeed.
	succeed atomic.Uint64
	// lost is the time that packet lost.
	lost atomic.Uint64
	// min is the minimum of probing time.
	min atomic.Uint64
	// max is the maximum of probing time.
	max atomic.Uint64

	avg  atomic.Uint64
	mdev atomic.Uint64
}

func NewState() *State {
	opts := &State{}
	opts.min.Store(math.MaxUint64)
	return opts
}

func (s *State) Add(t int, success bool) {
	s.probed.Add(1)
	if !success {
		s.lost.Add(1)
		return
	}

	uintTime := uint64(t)
	s.succeed.Add(1)
	if s.min.Load() > uintTime {
		s.min.Store(uintTime)
	}
	if s.max.Load() < uintTime {
		s.max.Store(uintTime)
	}

	avg := s.avg.Load()
	if avg == 0 {
		s.avg.Store(uintTime)
	} else {
		s.avg.Store((avg + uintTime) / 2)
	}

	adev := s.mdev.Load()
	abs := diffAbs(adev, avg)
	if adev == 0 {
		s.mdev.Store(abs)
	} else {
		s.mdev.Store((adev + abs) / 2)
	}
}

func (s *State) Probed() uint64 {
	return s.probed.Load()
}

func (s *State) Succeed() uint64 {
	return s.succeed.Load()
}

func (s *State) Lost() uint64 {
	return s.lost.Load()
}

func (s *State) Min() uint64 {
	return s.max.Load()
}

func (s *State) Max() uint64 {
	return s.max.Load()
}

func (s *State) Avg() uint64 {
	return s.avg.Load()
}

func (s *State) Mdev() uint64 {
	return s.mdev.Load()
}

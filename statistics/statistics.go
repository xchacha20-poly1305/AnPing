// Package statistics use to statistical ping.
package statistics

import (
	"math"

	"github.com/sagernet/sing/common/atomic"
)

type StatisticsGetter interface {
	Probed() uint64
	Succeed() uint64
	Lost() uint64
	Min() uint64
	Max() uint64
	Avg() uint64
	Mdev() uint64
}

var (
	_ StatisticsGetter = (*Statistics)(nil)
	_ StatisticsGetter = (*statisticsInternal)(nil)
)

type Statistics struct {
	sta statisticsInternal
}

func (s *Statistics) Getter() StatisticsGetter {
	return &s.sta
}

func (s *Statistics) Probed() uint64 {
	return s.sta.Probed()
}

func (s *Statistics) Succeed() uint64 {
	return s.sta.Succeed()
}

func (s *Statistics) Lost() uint64 {
	return s.sta.Lost()
}

func (s *Statistics) Min() uint64 {
	return s.sta.Min()
}

func (s *Statistics) Max() uint64 {
	return s.sta.Max()
}

func (s *Statistics) Avg() uint64 {
	return s.sta.Avg()
}

func (s *Statistics) Mdev() uint64 {
	return s.sta.Mdev()
}

type statisticsInternal struct {
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

func NewStatistics() *Statistics {
	s := &Statistics{}
	s.sta.min.Store(math.MaxUint64)
	return s
}

func (s *Statistics) Add(t uint64, success bool) {
	s.sta.probed.Add(1)
	if !success {
		s.sta.lost.Add(1)
		return
	}

	s.sta.succeed.Add(1)
	if s.sta.min.Load() > t {
		s.sta.min.Store(t)
	}
	if s.sta.max.Load() < t {
		s.sta.max.Store(t)
	}

	avg := s.sta.avg.Load()
	if avg == 0 {
		s.sta.avg.Store(t)
	} else {
		s.sta.avg.Store((avg + t) / 2)
	}

	adev := s.sta.mdev.Load()
	abs := diffAbs(adev, avg)
	if adev == 0 {
		s.sta.mdev.Store(abs)
	} else {
		s.sta.mdev.Store((adev + abs) / 2)
	}
}

func (s *statisticsInternal) Probed() uint64 {
	return s.probed.Load()
}

func (s *statisticsInternal) Succeed() uint64 {
	return s.succeed.Load()
}

func (s *statisticsInternal) Lost() uint64 {
	return s.lost.Load()
}

func (s *statisticsInternal) Min() uint64 {
	return s.max.Load()
}

func (s *statisticsInternal) Max() uint64 {
	return s.max.Load()
}

func (s *statisticsInternal) Avg() uint64 {
	return s.avg.Load()
}

func (s *statisticsInternal) Mdev() uint64 {
	return s.mdev.Load()
}

package libtcping

func newStatistics() *statistics {
	return &statistics{
		packetsSent: 0,
		packetsRecv: 0,
		packetsLoss: 0,

		minRtt: -1,
		maxRtt: -1,

		avgRtt:    -1,
		stddevRtt: -1,
	}
}

type statistics struct {
	packetsSent int
	packetsRecv int
	packetsLoss int

	minRtt int
	maxRtt int

	avgRtt    int
	stddevRtt int
}

func (s *statistics) freshRtt(rtt int) {
	// minRtt < 0: not set
	// rtt < minRtt: fresh
	if s.minRtt < 0 || rtt < s.minRtt {
		s.minRtt = rtt
	}

	if s.maxRtt < 0 || rtt > s.maxRtt {
		s.maxRtt = rtt
	}

	// init avg
	if s.avgRtt < 0 {
		s.avgRtt = rtt
	} else {
		s.avgRtt = (s.avgRtt + rtt) / 2
	}

	s.stddevRtt = s.maxRtt - s.minRtt
}

func (s *statistics) GetSent() int {
	return s.packetsSent
}

func (s *statistics) GetRecv() int {
	return s.packetsRecv
}

func (s *statistics) GetLoss() int {
	return s.packetsLoss
}

func (s *statistics) GetMinRtt() int {
	return s.minRtt
}

func (s *statistics) GetMaxRtt() int {
	return s.maxRtt
}

func (s *statistics) GetAvgRtt() int {
	return s.avgRtt
}

func (s *statistics) GetStddevRtt() int {
	return s.stddevRtt
}

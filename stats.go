package avchecker

type stats struct {
	TryCount     int     `json:"try_count"`
	SuccessCount int     `json:"success_count"`
	SuccessRatio float64 `json:"success_ratio"`
}

func (b *stats) reset() {
	b.TryCount = 0
	b.SuccessCount = 0
}

func (s *stats) compute() {
	s.SuccessRatio = float64(s.SuccessCount) / float64(s.TryCount)
}

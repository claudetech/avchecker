package avchecker

type stats struct {
	TryCount     int     `json:"try_count"`
	SuccessCount int     `json:"success_count"`
	SuccessRatio float64 `json:"success_ratio"`
	totalTime    int64
	AverageTime  float64 `json:"average_time"`
	extraFields  map[string]interface{}
}

func (b *stats) reset() {
	b.TryCount = 0
	b.SuccessCount = 0
	b.totalTime = 0
	b.AverageTime = 0
}

func (s *stats) toMap() map[string]interface{} {
	m := map[string]interface{}{
		"try_count":     s.TryCount,
		"success_count": s.SuccessCount,
		"success_ratio": s.SuccessRatio,
	}
	if s.AverageTime > 0 {
		m["average_time"] = s.AverageTime
	}
	for k, v := range s.extraFields {
		m[k] = v
	}
	return m
}

func (s *stats) compute() {
	s.SuccessRatio = float64(s.SuccessCount) / float64(s.TryCount)
	if s.SuccessCount > 0 {
		s.AverageTime = float64(s.totalTime) / float64(s.SuccessCount) / 1000000.0
	}
}

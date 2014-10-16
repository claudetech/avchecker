package avchecker

type stats struct {
	TryCount     int     `json:"try_count"`
	SuccessCount int     `json:"success_count"`
	SuccessRatio float64 `json:"success_ratio"`
	extraFields  map[string]interface{}
}

func (b *stats) reset() {
	b.TryCount = 0
	b.SuccessCount = 0
}

func (s *stats) toMap() map[string]interface{} {
	m := map[string]interface{}{
		"try_count":     s.TryCount,
		"success_count": s.SuccessCount,
		"success_ratio": s.SuccessRatio,
	}
	for k, v := range s.extraFields {
		m[k] = v
	}
	return m
}

func (s *stats) compute() {
	s.SuccessRatio = float64(s.SuccessCount) / float64(s.TryCount)
}

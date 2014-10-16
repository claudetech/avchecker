package avchecker

import (
	"encoding/json"
)

type Formatter interface {
	Format(stats *stats) ([]byte, error)
}

type JsonFormatter struct{}

func (f *JsonFormatter) Format(s *stats) ([]byte, error) {
	return json.Marshal(s.toMap())
}

package avchecker

import (
	"time"
)

type Reporter interface {
	SendStats() error
}

type Stats interface {
	String() string
}

type BaseStats struct {
	TryCount     int
	SuccessCount int
	Average      float32
}

func baseStats() {

}

type JsonStats struct {
	baseStats
}

type AvailabilityChecker struct {
	Url            string
	CheckInterval  time.Duration
	ReportInterval time.Duration
	Reporter       Reporter
	stats          baseStats
	lastReport     time.Time
}

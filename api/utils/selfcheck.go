package utils

import (
	"runtime"
	"time"

	"github.com/runar-rkmedia/gabyoall/logger"
)

type SelfCheckLimit struct {
	MemoryMB   uint64
	GoRoutines int
	Streaks    int
	Interval   time.Duration
}

// SelfCheck will watch the applications metrics, and crash it if it is considered to be unhealthy.
// This is only for use outside of kubernetes/docker
func SelfCheck(limits SelfCheckLimit, log logger.AppLogger, unhealhtyStreak int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	memoryMB := bToMb(m.HeapAlloc)
	goroutines := runtime.NumGoroutine()
	healthy := true
	if goroutines > limits.GoRoutines {
		healthy = false
	}
	if memoryMB > limits.MemoryMB {
		healthy = false
	}
	if healthy {

		unhealhtyStreak = 0
	} else {
		unhealhtyStreak += 1
	}
	l := logger.With(log.With().
		Int("unhealhtyStreak", unhealhtyStreak).
		Int("goroutines", goroutines).
		Uint64("memobryMB", memoryMB).
		Bool("healthy", healthy).
		Interface("limits", limits).
		Logger())
	if unhealhtyStreak > limits.Streaks {
		l.Fatal().Msg("Self-check reported unhealthy too many times in a row.")
	} else if !healthy {
		l.Warn().Msg("Unhealthy Self-check-result")

	} else {
		if l.HasDebug() {
			l.Debug().Msg("Self-check-result")
		}
	}
	time.Sleep(limits.Interval)
	SelfCheck(limits, log, unhealhtyStreak)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

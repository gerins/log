package log

import (
	"fmt"
	"time"
)

const (
	durationCallerName = "DURATION"
)

// processData is a data model for holding process information
type processData struct {
	request   *request  // Request data
	name      string    // Process name
	timeStart time.Time // Process start time
}

// Stop the total duration a process could take
func (p processData) Stop() {
	duration := float64(time.Since(p.timeStart).Nanoseconds()) / 1e6
	msg := fmt.Sprintf("[%.3fms] %s", duration, p.name)
	p.request.subLogs = append(p.request.subLogs, subLog{Level: GetCaller(durationCallerName, subLogSkipLevel), Message: msg})
}

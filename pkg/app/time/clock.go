package time

import "time"

type Clock interface {
	GetCurrentTime() int64
}

type SystemClock struct{}

func (clock SystemClock) GetCurrentTime() int64 {
	return time.Now().UnixNano()
}

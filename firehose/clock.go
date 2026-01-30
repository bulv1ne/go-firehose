package firehose

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (rc RealClock) Now() time.Time {
	return time.Now()
}

type MockClock struct {
	CurrentTime time.Time
}

func (mc *MockClock) Now() time.Time {
	return mc.CurrentTime
}

func (mc *MockClock) Advance(duration time.Duration) {
	mc.CurrentTime = mc.CurrentTime.Add(duration)
}

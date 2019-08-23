package mocks

import "github.com/stretchr/testify/mock"

type MockClock struct {
	mock.Mock
}

func (clock *MockClock) GetCurrentTime() int64 {
	args := clock.Called()
	return args.Get(0).(int64)
}

func (clock *MockClock) GivenGetCurrentTimeReturns(time int64) {
	clock.ExpectedCalls = []*mock.Call{}
	clock.On("GetCurrentTime").Return(time)
}

func CreateMockClock() *MockClock {
	mockClock := new(MockClock)
	mockClock.GivenGetCurrentTimeReturns(0)
	return mockClock
}

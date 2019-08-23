package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockLoggerSink mock of a Sink
type MockLoggerSink struct {
	mock.Mock
}

// Info mock
func (sink *MockLoggerSink) Info(message string, fields map[string]interface{}) {
	sink.Called(message, fields)
}

// Error mock
func (sink *MockLoggerSink) Error(err error, fields map[string]interface{}) {
	sink.Called(err, fields)
}

func CreateMockLoggerSink() *MockLoggerSink {
	sink := new(MockLoggerSink)
	sink.On("Info", mock.Anything, mock.Anything).Return()
	sink.On("Error", mock.Anything, mock.Anything).Return()
	return sink
}

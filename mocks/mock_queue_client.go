package mocks

import (
	"kube-job-runner/pkg/app/queue"

	"github.com/stretchr/testify/mock"
)

type MockQueueClient struct {
	mock.Mock
}

func (queueClient *MockQueueClient) SendMessage(message string) error {
	args := queueClient.Called(message)
	return args.Error(0)
}

func (queueClient *MockQueueClient) ReceiveMessages() ([]queue.Message, error) {
	args := queueClient.Called()
	return args.Get(0).([]queue.Message), args.Error(1)
}

func (queueClient *MockQueueClient) GivenSendMessageSucceeds() {
	queueClient.ExpectedCalls = []*mock.Call{}
	queueClient.On("SendMessage", mock.Anything).Return(nil)
}

func (queueClient *MockQueueClient) GivenSendMessageFails(err error) {
	queueClient.ExpectedCalls = []*mock.Call{}
	queueClient.On("SendMessage", mock.Anything).Return(err)
}

func (queueClient *MockQueueClient) GivenReceiveMessagesReturns(messages []queue.Message) {
	queueClient.ExpectedCalls = []*mock.Call{}
	queueClient.On("ReceiveMessages").Return(messages, nil)
}

func CreateMockQueueClient() *MockQueueClient {
	mockQueueClient := new(MockQueueClient)
	mockQueueClient.GivenSendMessageSucceeds()
	return mockQueueClient
}

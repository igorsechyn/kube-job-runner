package mocks

import (
	"context"

	"kube-job-runner/pkg/app/job"

	"github.com/stretchr/testify/mock"
)

type MockJobClient struct {
	mock.Mock
}

func (jobClient *MockJobClient) SubmitJob(job job.Job) (string, error) {
	args := jobClient.Called(job)
	return args.String(0), args.Error(1)
}

func (jobClient *MockJobClient) DeleteJob(jobID string) error {
	args := jobClient.Called(jobID)
	return args.Error(0)
}

func (jobClient *MockJobClient) GetJobIDForPod(podID string) (string, error) {
	args := jobClient.Called(podID)
	return args.String(0), args.Error(1)
}

func (jobClient *MockJobClient) DeleteEvent(eventID string) error {
	args := jobClient.Called(eventID)
	return args.Error(0)
}

func (jobClient *MockJobClient) WatchJobs(ctx context.Context) {
	jobClient.Called(ctx)
}

func (jobClient *MockJobClient) WatchEvents(ctx context.Context) {
	jobClient.Called(ctx)
}

func (jobClient *MockJobClient) AddJobStatusListener(listener job.StatusListener) {
	jobClient.Called(listener)
}

func (jobClient *MockJobClient) AddPodEventsListener(listener job.PodEventListener) {
	jobClient.Called(listener)
}

func (jobClient *MockJobClient) GivenSubmitJobSucceeds() {
	jobClient.On("SubmitJob", mock.AnythingOfType("job.Job")).Return("default-job-name", nil)
}

func (jobClient *MockJobClient) GivenGetJobIDForPodReturns(jobID string) {
	jobClient.On("GetJobIDForPod", mock.Anything).Return(jobID, nil)
}

func (jobClient *MockJobClient) GivenGetJobIDForPodFails(err error) {
	jobClient.On("GetJobIDForPod", mock.Anything).Return("", err)
}

func (jobClient *MockJobClient) GivenDeleteJobSucceeds() {
	jobClient.On("DeleteJob", mock.Anything).Return(nil)
}

func (jobClient *MockJobClient) GivenDeleteEventSucceeds() {
	jobClient.On("DeleteEvent", mock.Anything).Return(nil)
}

func (jobClient *MockJobClient) GivenDeleteJobFailed(err error) {
	jobClient.On("DeleteJob", mock.Anything).Return(err)
}

func (jobClient *MockJobClient) GivenDeleteEventFailed(err error) {
	jobClient.On("DeleteEvent", mock.Anything).Return(err)
}

func (jobClient *MockJobClient) GivenSubmitJobFails(err error) {
	jobClient.On("SubmitJob", mock.AnythingOfType("job.Job")).Return("", err)
}

func CreateMockJobClient() *MockJobClient {
	mockJobClient := new(MockJobClient)
	return mockJobClient
}

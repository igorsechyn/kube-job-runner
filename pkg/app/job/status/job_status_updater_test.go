// +build unit

package status_test

import (
	"fmt"
	"testing"

	"kube-job-runner/mocks"
	"kube-job-runner/pkg/app/data"
	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/job/status"
	"kube-job-runner/pkg/app/queue"
	"kube-job-runner/pkg/app/reporter"
)

func whenAStatusUpdateIsReceived(jobStatus job.Status, allMocks mocks.AllMocks) {
	logger := reporter.New(allMocks.MockLoggerSink)
	queueManager := queue.NewManager(allMocks.MockQueueClient, logger)
	jobStatusUpdater := status.JobStatusUpdater{
		QueueManager: queueManager,
	}

	jobStatusUpdater.Process(jobStatus)
}

type testCase struct {
	description    string
	status         job.Status
	expectedStatus string
}

func TestJobStatusUpdate(t *testing.T) {
	testCases := []testCase{
		{
			description:    "it should send an acknowledged status update messages, if no containers have started yet",
			status:         job.Status{JobID: "some-id", FailedContainers: 0, RunningContainers: 0, SucceededContainers: 0},
			expectedStatus: data.JobAcknowledged,
		},
		{
			description:    "it should send an in progress status update messages, if all containers are running",
			status:         job.Status{JobID: "some-id", FailedContainers: 0, RunningContainers: 2, SucceededContainers: 0},
			expectedStatus: data.JobInProgress,
		},
		{
			description:    "it should send an in progress status update messages, if at least 1 container is running",
			status:         job.Status{JobID: "some-id", FailedContainers: 1, RunningContainers: 1, SucceededContainers: 1},
			expectedStatus: data.JobInProgress,
		},
		{
			description:    "it should send a succeeded status update messages, if all containers succeeded",
			status:         job.Status{JobID: "some-id", FailedContainers: 0, RunningContainers: 0, SucceededContainers: 2},
			expectedStatus: data.JobSucceeded,
		},
		{
			description:    "it should send a failed status update messages, if all containers failed",
			status:         job.Status{JobID: "some-id", FailedContainers: 2, RunningContainers: 0, SucceededContainers: 0},
			expectedStatus: data.JobFailed,
		},
		{
			description:    "it should send a failed status update messages, if some containers failed",
			status:         job.Status{JobID: "some-id", FailedContainers: 1, RunningContainers: 0, SucceededContainers: 1},
			expectedStatus: data.JobFailed,
		},
	}

	for _, testCase := range testCases {
		allMocks := mocks.InitMocks()

		t.Run(testCase.description, func(t *testing.T) {
			allMocks.MockQueueClient.GivenSendMessageSucceeds()

			whenAStatusUpdateIsReceived(testCase.status, allMocks)

			allMocks.MockQueueClient.AssertCalled(t, "SendMessage", fmt.Sprintf(`{"jobID":"some-id","status":"%v","type":"JOB_STATUS_UPDATE"}`, testCase.expectedStatus))
		})
	}
}

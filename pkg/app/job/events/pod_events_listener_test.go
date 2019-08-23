// +build unit

package events_test

import (
	"fmt"
	"testing"

	"kube-job-runner/mocks"
	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/job/events"
	"kube-job-runner/pkg/app/queue"
	"kube-job-runner/pkg/app/reporter"

	"github.com/stretchr/testify/mock"
)

func whenAnEventIsReceived(event job.PodEvent, allMocks mocks.AllMocks) {
	logger := reporter.New(allMocks.MockLoggerSink)
	queueManager := queue.NewManager(allMocks.MockQueueClient, logger)
	podEventListener := events.Listener{
		JobClient:    allMocks.MockJobClient,
		QueueManager: queueManager,
		Reporter:     logger,
	}

	podEventListener.Process(event)
}

func TestListener_Process(t *testing.T) {
	t.Run("it should get job id of the failed pod", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenGetJobIDForPodReturns("job-id")
		allMocks.MockJobClient.GivenDeleteEventSucceeds()

		whenAnEventIsReceived(job.PodEvent{PodID: "pod-id", Status: "Failed"}, allMocks)

		allMocks.MockJobClient.AssertCalled(t, "GetJobIDForPod", "pod-id")
	})

	t.Run("it should send a update job status message, if the status is failed", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenDeleteEventSucceeds()
		allMocks.MockJobClient.GivenGetJobIDForPodReturns("job-id")

		whenAnEventIsReceived(job.PodEvent{PodID: "pod-id", Status: "Failed", Reason: "Error"}, allMocks)

		allMocks.MockQueueClient.AssertCalled(t, "SendMessage", queue.MessageBody(`{"jobID":"job-id","status":"Failed","type":"JOB_STATUS_UPDATE","reason":"Error"}`))
	})

	t.Run("it should not send a message, if the status is not failed", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenGetJobIDForPodReturns("job-id")

		whenAnEventIsReceived(job.PodEvent{PodID: "pod-id", Status: "Started"}, allMocks)

		allMocks.MockQueueClient.AssertNotCalled(t, "SendMessage", mock.Anything)
		allMocks.MockJobClient.AssertNotCalled(t, "GetJobIDForPod", mock.Anything)
		allMocks.MockJobClient.AssertNotCalled(t, "DeleteEvent", mock.Anything)
	})

	t.Run("it should return, if getting the job id fails", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenGetJobIDForPodFails(fmt.Errorf("some error"))

		whenAnEventIsReceived(job.PodEvent{PodID: "pod-id", Status: "Failed"}, allMocks)

		allMocks.MockQueueClient.AssertNotCalled(t, "SendMessage", mock.Anything)
		allMocks.MockJobClient.AssertNotCalled(t, "DeleteEvent", mock.Anything)
	})

	t.Run("it should log an error, if getting the job id fails", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenGetJobIDForPodFails(fmt.Errorf("some error"))

		whenAnEventIsReceived(job.PodEvent{PodID: "pod-id", Status: "Failed"}, allMocks)

		allMocks.MockLoggerSink.AssertCalled(t, "Error", fmt.Errorf("some error"), map[string]interface{}{"code": "job.get.id.error"})
	})

	t.Run("it should delete event, if the message was sent successfully", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenGetJobIDForPodReturns("job-id")
		allMocks.MockJobClient.GivenDeleteEventSucceeds()

		whenAnEventIsReceived(job.PodEvent{ID: "event-id", Status: "Failed"}, allMocks)

		allMocks.MockJobClient.AssertCalled(t, "DeleteEvent", "event-id")
	})

	t.Run("it should not delete event, if the message was not sent", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenGetJobIDForPodReturns("job-id")
		allMocks.MockQueueClient.GivenSendMessageFails(fmt.Errorf("some error"))

		whenAnEventIsReceived(job.PodEvent{ID: "event-id", Status: "Failed"}, allMocks)

		allMocks.MockJobClient.AssertNotCalled(t, "DeleteEvent", mock.Anything)
	})

	t.Run("it should log an error, if sending message failed", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenGetJobIDForPodReturns("job-id")
		allMocks.MockQueueClient.GivenSendMessageFails(fmt.Errorf("some error"))

		whenAnEventIsReceived(job.PodEvent{ID: "event-id", Status: "Failed"}, allMocks)

		allMocks.MockLoggerSink.AssertCalled(t, "Error", fmt.Errorf("some error"), map[string]interface{}{"code": "send.event.message.error"})
	})
}

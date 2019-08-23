// +build unit

package messages_test

import (
	"fmt"
	"testing"

	"kube-job-runner/mocks"
	"kube-job-runner/pkg/app/data"
	"kube-job-runner/pkg/app/queue"

	"github.com/stretchr/testify/mock"
)

func TestJobUpdateMessageProcessing(t *testing.T) {
	t.Run("it should ignore message, if it is not update job staus message", func(t *testing.T) {
		allMocks := mocks.InitMocks()

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"wrongField":"uuid","type":"WRONG_TYPE"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockStore.AssertNotCalled(t, "GetDocuments", mock.Anything)
		allMocks.MockStore.AssertNotCalled(t, "PutDocument", mock.Anything)
	})

	t.Run("it should retrieve job details from data store", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{{Status: "InProgress"}})
		allMocks.MockStore.GivenPutDocumentSucceeds()

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"status": "InProgress", "jobID":"uuid","type":"JOB_STATUS_UPDATE"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockStore.AssertCalled(t, "GetDocuments", "uuid")
	})

	t.Run("it should update job status, if it changed", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockClock.GivenGetCurrentTimeReturns(100)
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{{Image: "image", Tag: "tag", Status: "InProgress", Timestamp: 10}})
		allMocks.MockStore.GivenPutDocumentSucceeds()

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"status": "Succeeded", "jobID":"uuid","type":"JOB_STATUS_UPDATE"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockStore.AssertCalled(t, "PutDocument", data.Document{Image: "image", Tag: "tag", Status: "Succeeded", Timestamp: 100})
	})

	t.Run("it should not update job status, if it did not changed", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{{Image: "image", Tag: "tag", Status: "InProgress", Timestamp: 10}})

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"status": "InProgress", "jobID":"uuid","type":"JOB_STATUS_UPDATE"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockStore.AssertNotCalled(t, "PutDocument", mock.Anything)
	})

	t.Run("it should not update job status, if retrieving job fails", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenGetDocumentsFailed(fmt.Errorf("get document failed"))

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"status": "InProgress", "jobID":"uuid","type":"JOB_STATUS_UPDATE"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockStore.AssertNotCalled(t, "PutDocument", mock.Anything)
	})

	t.Run("it should not update job status to InProgress, if it already succeeded", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenPutDocumentSucceeds()
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{{Image: "image", Tag: "tag", Status: "Succeeded", Timestamp: 10}})

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"status": "InProgress", "jobID":"uuid","type":"JOB_STATUS_UPDATE"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockStore.AssertNotCalled(t, "PutDocument", mock.Anything)
	})

	t.Run("it should not update job status to InProgress, if it already failed", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenPutDocumentSucceeds()
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{{Image: "image", Tag: "tag", Status: "Failed", Timestamp: 10}})

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"status": "InProgress", "jobID":"uuid","type":"JOB_STATUS_UPDATE"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockStore.AssertNotCalled(t, "PutDocument", mock.Anything)
	})

	t.Run("it should send job cleanup message, when job finished successfully", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockClock.GivenGetCurrentTimeReturns(100)
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{{Status: "InProgress", Timestamp: 10}})
		allMocks.MockStore.GivenPutDocumentSucceeds()

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"status": "Succeeded", "jobID":"uuid","type":"JOB_STATUS_UPDATE"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockQueueClient.AssertCalled(t, "SendMessage", queue.MessageBody(`{"jobID":"uuid","type":"CLEANUP_JOB"}`))
	})

	t.Run("it should send job cleanup message, when job failed", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockClock.GivenGetCurrentTimeReturns(100)
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{{Status: "InProgress", Timestamp: 10}})
		allMocks.MockStore.GivenPutDocumentSucceeds()

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"status": "Failed", "jobID":"uuid","type":"JOB_STATUS_UPDATE"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockQueueClient.AssertCalled(t, "SendMessage", queue.MessageBody(`{"jobID":"uuid","type":"CLEANUP_JOB"}`))
	})
}

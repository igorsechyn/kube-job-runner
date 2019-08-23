// +build unit

package messages_test

import (
	"fmt"
	"testing"

	"kube-job-runner/mocks"
	"kube-job-runner/pkg/app/data"
	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/queue"

	"github.com/stretchr/testify/mock"
)

func TestJobCreationMessageProcessing(t *testing.T) {
	t.Run("it should ignore message, if it is not create job request message", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"wrongField":"uuid","type":"CREATE_JOB"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockStore.AssertNotCalled(t, "GetDocuments", mock.Anything)
		allMocks.MockJobClient.AssertNotCalled(t, "SubmitJob", mock.Anything)
	})

	t.Run("it should retrieve job details from data store", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{{}})
		allMocks.MockJobClient.GivenSubmitJobSucceeds()

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"jobID":"uuid","type":"CREATE_JOB"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockStore.AssertCalled(t, "GetDocuments", "uuid")
	})

	t.Run("it should not call job client, if retrieving job fails", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenSubmitJobSucceeds()
		allMocks.MockStore.GivenGetDocumentsFailed(fmt.Errorf("get document failed"))

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"jobID":"uuid","type":"CREATE_JOB"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockJobClient.AssertNotCalled(t, "SubmitJob")
	})

	t.Run("it should create a new job", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{{JobID: "job-id", Tag: "some-tag", Image: "some-image"}})
		allMocks.MockJobClient.GivenSubmitJobSucceeds()

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"jobID":"uuid","type":"CREATE_JOB"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockJobClient.AssertCalled(t, "SubmitJob", job.Job{JobName: "job-id", Tag: "some-tag", Image: "some-image"})
	})

	t.Run("it logs an error, if creating job fails", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenPutDocumentSucceeds()
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{{}})
		allMocks.MockJobClient.GivenSubmitJobFails(fmt.Errorf("creating job failed"))

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"jobID":"uuid","type":"CREATE_JOB"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockLoggerSink.AssertCalled(t, "Error", fmt.Errorf("creating job failed"), mock.Anything)
	})

	t.Run("it should update job status to failed, if creating job fails", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{{}})
		allMocks.MockStore.GivenPutDocumentSucceeds()
		allMocks.MockJobClient.GivenSubmitJobFails(fmt.Errorf("creating job failed"))

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   queue.MessageBody(`{"jobID":"uuid","type":"CREATE_JOB"}`),
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockStore.AssertCalled(t, "PutDocument", data.Document{Status: "Failed"})
	})
}

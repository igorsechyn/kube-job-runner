// +build unit

package job_test

import (
	"fmt"
	"testing"

	"kube-job-runner/mocks"
	"kube-job-runner/pkg/app/data"
	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/queue"
	"kube-job-runner/pkg/app/reporter"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getJobService(allMocks mocks.AllMocks) job.Service {
	reporter := reporter.New(allMocks.MockLoggerSink)
	jobService := job.Service{
		IDGenerator: allMocks.MockIdGenerator,
		Reporter:    reporter,
		Clock:       allMocks.MockClock,
		DataStore:   allMocks.MockStore,
		Queue:       queue.NewManager(allMocks.MockQueueClient, reporter),
	}

	return jobService
}

func whenJobRequestIsSubmitted(allMocks mocks.AllMocks) (string, error) {
	return whenJobRequestIsSubmittedWithDetails(allMocks, "default-image", "default-tag", 0)
}

func whenJobRequestIsSubmittedWithDetails(allMocks mocks.AllMocks, image, tag string, timeout int64) (string, error) {
	jobService := getJobService(allMocks)

	return jobService.SubmitJobCreationRequest(job.Request{Image: image, Tag: tag, Timeout: timeout})
}

func TestService_SubmitJobRequest(t *testing.T) {
	t.Run("it should call uuid generator for the job id", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenPutDocumentSucceeds()

		_, err := whenJobRequestIsSubmitted(allMocks)

		require.NoError(t, err)
		allMocks.MockIdGenerator.AssertCalled(t, "Generate")
	})

	t.Run("it should get the current time stamp", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenPutDocumentSucceeds()

		_, err := whenJobRequestIsSubmitted(allMocks)

		require.NoError(t, err)
		allMocks.MockClock.AssertCalled(t, "GetCurrentTime")
	})

	t.Run("it should add job details to event store", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenPutDocumentSucceeds()
		allMocks.MockClock.GivenGetCurrentTimeReturns(100)
		allMocks.MockIdGenerator.GivenGeneratorReturnsID("uuid")

		_, err := whenJobRequestIsSubmittedWithDetails(allMocks, "some-image", "some-tag", 100)

		require.NoError(t, err)
		allMocks.MockStore.AssertCalled(t, "PutDocument", data.Document{
			Image:     "some-image",
			Tag:       "some-tag",
			JobID:     "uuid",
			Status:    "Acknowledged",
			Timestamp: 100,
			Timeout:   100,
		})
	})

	t.Run("it should return an error, if saving job details fails", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		expectedError := fmt.Errorf("failed to save job details")
		allMocks.MockStore.GivenPutDocumentFails(expectedError)

		_, err := whenJobRequestIsSubmitted(allMocks)

		require.EqualError(t, err, "failed to save job details")
	})

	t.Run("it should put a message on a queue with job id", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockIdGenerator.GivenGeneratorReturnsID("job-uuid")
		allMocks.MockQueueClient.GivenSendMessageSucceeds()
		allMocks.MockStore.GivenPutDocumentSucceeds()

		_, err := whenJobRequestIsSubmitted(allMocks)

		require.NoError(t, err)
		allMocks.MockQueueClient.AssertCalled(t, "SendMessage", queue.MessageBody(`{"jobID":"job-uuid","type":"CREATE_JOB"}`))
	})

	t.Run("it should return an error, if sending a message fails", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockQueueClient.GivenSendMessageFails(fmt.Errorf("failed to send message"))
		allMocks.MockStore.GivenPutDocumentSucceeds()

		_, err := whenJobRequestIsSubmitted(allMocks)

		require.EqualError(t, err, "failed to send message")
	})
}

func whenGetJobDetailsIsCalled(allMocks mocks.AllMocks, jobId string) (job.Details, error) {
	jobService := getJobService(allMocks)
	return jobService.GetJobDetails(jobId)
}

func TestService_GetJobDetails(t *testing.T) {
	t.Run("it should return the latest job status details, when job id exists", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenGetDocumentsReturns([]data.Document{
			{Image: "image", Tag: "tag", JobID: "some-id", Status: "Succeeded", Timestamp: 3, Timeout: 100},
			{Image: "image", Tag: "tag", JobID: "some-id", Status: "InProgress", Timestamp: 2, Timeout: 100},
			{Image: "image", Tag: "tag", JobID: "some-id", Status: "Acknowledged", Timestamp: 1, Timeout: 100},
		})

		details, err := whenGetJobDetailsIsCalled(allMocks, "some-id")

		require.NoError(t, err)
		assert.Equal(t, job.Details{Image: "image", Tag: "tag", Name: "some-id", Status: "Succeeded", Timeout: 100}, details)
	})

	t.Run("it should return an error, when data store returns an error", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockStore.GivenGetDocumentsFailed(fmt.Errorf("no document for id"))

		_, err := whenGetJobDetailsIsCalled(allMocks, "some-id")

		require.EqualError(t, err, "no document for id")
	})
}

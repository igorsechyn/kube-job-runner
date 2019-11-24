// +build unit

package messages_test

import (
	"fmt"
	"testing"

	"kube-job-runner/mocks"
	"kube-job-runner/pkg/app/queue"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJobCleanupMessageProcessing(t *testing.T) {
	t.Run("it should ignore message, if it is not cleanup job request message", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   `{"wrongField":"uuid","type":"WRONG_TYPE_JOB"}`,
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockJobClient.AssertNotCalled(t, "CleanupJob", mock.Anything)
	})

	t.Run("it should delete job", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenDeleteJobSucceeds()

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   `{"jobID":"uuid","type":"CLEANUP_JOB"}`,
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockJobClient.AssertCalled(t, "DeleteJob", "uuid")
	})

	t.Run("it logs an error, if deleting job fails", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenDeleteJobFailed(fmt.Errorf("failed delete job"))

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   `{"jobID":"uuid","type":"CLEANUP_JOB"}`,
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockLoggerSink.AssertCalled(t, "Error", fmt.Errorf("failed delete job"), mock.Anything)
	})

	t.Run("it ignores an error, if the job does not exist any more", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenDeleteJobFailed(fmt.Errorf("job not found"))

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   `{"jobID":"uuid","type":"CLEANUP_JOB"}`,
				Delete: func() {},
			},
		}, allMocks)

		allMocks.MockLoggerSink.AssertNotCalled(t, "Error", mock.Anything, mock.Anything)
	})

	t.Run("it should delete message after processing it", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenDeleteJobSucceeds()
		deleted := false
		deleteFunction := func() { deleted = true }

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   `{"jobID":"uuid","type":"CLEANUP_JOB"}`,
				Delete: deleteFunction,
			},
		}, allMocks)

		assert.True(t, deleted, "message was not deleted")
	})

	t.Run("it should not delete message if processing failed", func(t *testing.T) {
		allMocks := mocks.InitMocks()
		allMocks.MockJobClient.GivenDeleteJobFailed(fmt.Errorf(("some error")))
		deleted := false
		deleteFunction := func() { deleted = true }

		whenMessagesFromQueueAreProcessed([]queue.Message{
			{
				Body:   `{"jobID":"uuid","type":"CLEANUP_JOB"}`,
				Delete: deleteFunction,
			},
		}, allMocks)

		assert.False(t, deleted, "message was deleted")
	})
}

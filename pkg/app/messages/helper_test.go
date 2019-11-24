// +build unit

package messages_test

import (
	"kube-job-runner/mocks"
	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/messages"
	"kube-job-runner/pkg/app/queue"
	"kube-job-runner/pkg/app/reporter"
)

func whenMessagesFromQueueAreProcessed(queueMessages []queue.Message, allMocks mocks.AllMocks) {
	logger := reporter.New(allMocks.MockLoggerSink)
	jobService := job.Service{
		JobClient: allMocks.MockJobClient,
		DataStore: allMocks.MockStore,
		Reporter:  logger,
		Clock:     allMocks.MockClock,
	}
	queueManager := queue.NewManager(allMocks.MockQueueClient, logger)

	createJobHandler := &messages.CreateJobHandler{
		JobService: &jobService,
		Reporter:   logger,
	}
	updateJobStatusHandler := &messages.UpdateJobStatusHandler{
		JobService:   &jobService,
		Reporter:     logger,
		QueueManager: queueManager,
	}
	cleanupJobHandler := &messages.CleanupJobHandler{
		JobService: &jobService,
		Reporter:   logger,
	}

	queueManager.AddListener(createJobHandler)
	queueManager.AddListener(updateJobStatusHandler)
	queueManager.AddListener(cleanupJobHandler)

	mockDeleteFunction(queueMessages)
	queueManager.ProcessMessages(queueMessages)
}

func mockDeleteFunction(queueMessages []queue.Message) {
	for _, message := range queueMessages {
		if message.Delete == nil {
			message.Delete = func() {}
		}
	}
}

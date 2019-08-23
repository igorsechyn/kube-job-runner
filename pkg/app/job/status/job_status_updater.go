package status

import (
	"kube-job-runner/pkg/app/data"
	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/queue"
)

type JobStatusUpdater struct {
	QueueManager *queue.Manager
}

func (updater *JobStatusUpdater) Process(status job.Status) {
	jobStatus := getJobStatus(status)
	updater.QueueManager.SendJobStatusUpdateMessage(status.JobID, jobStatus, "")
}

func getJobStatus(status job.Status) string {
	if status.RunningContainers > 0 {
		return data.JobInProgress
	}

	if status.RunningContainers == 0 && status.SucceededContainers > 0 && status.FailedContainers == 0 {
		return data.JobSucceeded
	}

	if status.RunningContainers == 0 && status.FailedContainers > 0 {
		return data.JobFailed
	}

	return data.JobAcknowledged
}

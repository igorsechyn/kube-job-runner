package messages

import (
	"bytes"
	"encoding/json"

	"kube-job-runner/pkg/app/data"
	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/queue"
	"kube-job-runner/pkg/app/reporter"
)

type UpdateJobStatusHandler struct {
	JobService   *job.Service
	Reporter     *reporter.Reporter
	QueueManager *queue.Manager
}

func (handler *UpdateJobStatusHandler) Accept(messageType string) bool {
	return messageType == queue.JobStatusUpdateType
}

func (handler *UpdateJobStatusHandler) Process(message string) error {
	jobStatusUpdate, err := parseUpdateJobStatusMessage(message)
	if err != nil {
		return err
	}

	err = handler.updateJobStatus(jobStatusUpdate)
	if err != nil {
		return err
	}

	err = handler.cleanupJob(jobStatusUpdate)

	return err
}

func (handler *UpdateJobStatusHandler) updateJobStatus(jobStatusUpdate queue.JobStatusUpdateMessage) error {
	details, err := handler.JobService.GetJobDetails(jobStatusUpdate.JobID)
	if err != nil {
		return err
	}

	if shouldUpdateStatus(details.Status, jobStatusUpdate.Status) {
		_, err = handler.JobService.SaveJobDetails(job.Details{
			Image:  details.Image,
			Name:   details.Name,
			Tag:    details.Tag,
			Status: jobStatusUpdate.Status,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func shouldUpdateStatus(current, new string) bool {
	if current == data.JobSucceeded || current == data.JobFailed {
		return false
	}

	return current != new
}

func (handler *UpdateJobStatusHandler) cleanupJob(jobStatusUpdate queue.JobStatusUpdateMessage) error {
	if jobStatusUpdate.Status == data.JobSucceeded || jobStatusUpdate.Status == data.JobFailed {
		err := handler.QueueManager.SendJobCleanupMessage(jobStatusUpdate.JobID)
		if err != nil {
			handler.Reporter.Error("send.cleanup.job.error", err, map[string]interface{}{"jobId": jobStatusUpdate.JobID})
			return err
		}
	}

	return nil
}

func parseUpdateJobStatusMessage(message string) (queue.JobStatusUpdateMessage, error) {
	var updateJobMessage queue.JobStatusUpdateMessage
	dec := json.NewDecoder(bytes.NewReader([]byte(message)))
	dec.DisallowUnknownFields()
	err := dec.Decode(&updateJobMessage)

	return updateJobMessage, err
}

package messages

import (
	"bytes"
	"encoding/json"

	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/queue"
	"kube-job-runner/pkg/app/reporter"
)

type CleanupJobHandler struct {
	JobService *job.Service
	Reporter   *reporter.Reporter
}

func (handler *CleanupJobHandler) Accept(messageType string) bool {
	return messageType == queue.JobCleanupType
}

func (handler *CleanupJobHandler) Process(message string) error {
	cleanupJobMessage, err := parseCleanupJobMessage(message)
	if err != nil {
		handler.Reporter.Error("cleanup.job.message.wrong.format", err, map[string]interface{}{"message": message})
		return err
	}

	err = handler.JobService.DeleteJob(cleanupJobMessage.JobID)
	return err
}

func parseCleanupJobMessage(message string) (queue.JobCleanupMessage, error) {
	var jobCleanupMessage queue.JobCleanupMessage
	dec := json.NewDecoder(bytes.NewReader([]byte(message)))
	dec.DisallowUnknownFields()
	err := dec.Decode(&jobCleanupMessage)

	return jobCleanupMessage, err
}

package messages

import (
	"bytes"
	"encoding/json"

	"kube-job-runner/pkg/app/data"
	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/queue"
	"kube-job-runner/pkg/app/reporter"
)

type CreateJobHandler struct {
	JobService *job.Service
	Reporter   *reporter.Reporter
}

func (handler *CreateJobHandler) Accept(messageType string) bool {
	return messageType == queue.CreateMessageType
}

func (handler *CreateJobHandler) Process(message queue.MessageBody) error {
	createJobMessage, err := parseCreateJobMessage(message)
	if err != nil {
		handler.Reporter.Error("create.job.message.wrong.format", err, map[string]interface{}{"message": message})
		return err
	}

	jobDetails, err := handler.JobService.GetJobDetails(createJobMessage.JobID)

	if err != nil {
		return err
	}

	err = handler.JobService.CreateJob(job.Job{JobName: jobDetails.Name, Image: jobDetails.Image, Tag: jobDetails.Tag})
	if err != nil {
		handler.failJob(jobDetails)
	}
	return err
}

func (handler *CreateJobHandler) failJob(details job.Details) {
	handler.JobService.SaveJobDetails(job.Details{
		Status: data.JobFailed,
		Image:  details.Image,
		Tag:    details.Tag,
		Name:   details.Name,
	})
}

func parseCreateJobMessage(message queue.MessageBody) (queue.CreateJobMessage, error) {
	var createJobMessage queue.CreateJobMessage
	dec := json.NewDecoder(bytes.NewReader([]byte(message)))
	dec.DisallowUnknownFields()
	err := dec.Decode(&createJobMessage)

	return createJobMessage, err
}

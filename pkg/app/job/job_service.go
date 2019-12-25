package job

import (
	"context"
	"fmt"
	"strings"

	"kube-job-runner/pkg/app/data"
	"kube-job-runner/pkg/app/id"
	"kube-job-runner/pkg/app/queue"
	"kube-job-runner/pkg/app/reporter"
	"kube-job-runner/pkg/app/time"
)

type Request struct {
	Image   string `json:"image"`
	Tag     string `json:"tag"`
	Timeout int64  `json:"timeout,omitempty"`
}

type Details struct {
	Image   string
	Tag     string
	Name    string
	Status  string
	Timeout int64
}

type Service struct {
	IDGenerator id.IDGenerator
	Reporter    *reporter.Reporter
	JobClient   Client
	Clock       time.Clock
	DataStore   data.Store
	Queue       *queue.Manager
}

func (jobService *Service) SubmitJobCreationRequest(request Request) (string, error) {
	details := Details{
		Image:   request.Image,
		Tag:     request.Tag,
		Name:    jobService.IDGenerator.Generate(),
		Status:  data.JobAcknowledged,
		Timeout: request.Timeout,
	}
	jobId, err := jobService.SaveJobDetails(details)
	if err != nil {
		return jobId, err
	}
	err = jobService.sendSubmitJobRequestMessage(jobId)
	return jobId, err
}

func (jobService *Service) GetJobDetails(id string) (Details, error) {
	documents, err := jobService.DataStore.GetDocuments(id)

	if err != nil {
		jobService.Reporter.Error("job.get.details.error", err, map[string]interface{}{"jobId": id})
		return Details{}, err
	}

	document := documents[0]

	return Details{
		Name:    document.JobID,
		Tag:     document.Tag,
		Image:   document.Image,
		Status:  document.Status,
		Timeout: document.Timeout,
	}, nil
}

func (jobService *Service) CreateJob(job Job) error {
	_, err := jobService.JobClient.SubmitJob(job)
	if err != nil {
		jobService.Reporter.Error("job.create.error", err, map[string]interface{}{"jobId": job.JobName})
		return err
	}

	jobService.Reporter.Info("job.created", "Successfully created job", map[string]interface{}{"jobId": job.JobName})
	return nil
}

func (jobService *Service) DeleteJob(jobID string) error {
	err := jobService.JobClient.DeleteJob(jobID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}

		jobService.Reporter.Error("delete.job.error", err, map[string]interface{}{"jobId": jobID})
		return err
	}

	jobService.Reporter.Info("delete.job", fmt.Sprintf("Successfully deleted job %v", jobID), map[string]interface{}{"jobId": jobID})
	return nil
}

func (jobService *Service) SaveJobDetails(details Details) (string, error) {
	timestamp := jobService.Clock.GetCurrentTime()
	err := jobService.DataStore.PutDocument(data.Document{
		Image:     details.Image,
		Tag:       details.Tag,
		JobID:     details.Name,
		Status:    details.Status,
		Timestamp: timestamp,
		Timeout:   details.Timeout,
	})

	if err != nil {
		jobService.Reporter.Error(
			"job.request.save.error", err, map[string]interface{}{"jobId": details.Name},
		)
		return "", err
	}

	jobService.Reporter.Info(
		"job.request.saved",
		fmt.Sprintf("Successfully saved job details %v, status: %v", details.Name, details.Status),
		map[string]interface{}{"jobId": details.Name},
	)

	return details.Name, nil
}

func (jobService *Service) sendSubmitJobRequestMessage(jobId string) error {
	err := jobService.Queue.SendCreateJobMessage(jobId)
	if err != nil {
		jobService.Reporter.Error(
			"job.request.message.send.error", err, map[string]interface{}{"jobId": jobId},
		)
		return err
	}

	jobService.Reporter.Info(
		"job.request.message.sent",
		fmt.Sprintf("Successfully sent job request message for %v", jobId),
		map[string]interface{}{"jobId": jobId},
	)

	return nil
}

func (jobService *Service) WatchJobs(ctx context.Context, listener StatusListener) {
	jobService.JobClient.AddJobStatusListener(listener)
	jobService.JobClient.WatchJobs(ctx)
}

func (jobService *Service) WatchEvents(ctx context.Context, listener PodEventListener) {
	jobService.JobClient.AddPodEventsListener(listener)
	jobService.JobClient.WatchEvents(ctx)
}

package job

import (
	"context"
)

type Job struct {
	JobName string
	Image   string
	Tag     string
}

type Status struct {
	JobID               string
	RunningContainers   int32
	FailedContainers    int32
	SucceededContainers int32
}

type StatusListener interface {
	Process(status Status)
}

type PodEvent struct {
	ID     string
	PodID  string
	Status string
	Reason string
}

type PodEventListener interface {
	Process(event PodEvent)
}

type Client interface {
	SubmitJob(job Job) (string, error)
	WatchJobs(ctx context.Context)
	DeleteJob(jobID string) error
	DeleteEvent(eventID string) error
	WatchEvents(ctx context.Context)
	GetJobIDForPod(podID string) (string, error)
	AddJobStatusListener(listener StatusListener)
	AddPodEventsListener(listener PodEventListener)
}

package worker

import (
	"context"

	"kube-job-runner/pkg/app"
	"kube-job-runner/pkg/app/job/status"
)

type JobStatusProcessor struct{}

func (processor *JobStatusProcessor) Run(ctx context.Context, app *app.App) {
	listener := &status.JobStatusUpdater{QueueManager: app.Queue}

	app.JobService.WatchJobs(ctx, listener)
}

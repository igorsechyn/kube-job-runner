package worker

import (
	"context"

	"kube-job-runner/pkg/app"
	"kube-job-runner/pkg/app/job/events"
)

type PodEventsProcessor struct{}

func (processor *PodEventsProcessor) Run(ctx context.Context, app *app.App) {
	listener := &events.Listener{
		Reporter:     app.Reporter,
		JobClient:    app.JobClient,
		QueueManager: app.Queue,
	}

	app.JobService.WatchEvents(listener)
}

package worker

import (
	"context"
	"time"

	"kube-job-runner/pkg/app"
	"kube-job-runner/pkg/app/messages"
)

type MessagesProcessor struct{}

func (worker *MessagesProcessor) Run(ctx context.Context, app *app.App) {
	createJobListener := &messages.CreateJobHandler{
		Reporter:   app.Reporter,
		JobService: app.JobService,
	}

	updateJobStatusListener := &messages.UpdateJobStatusHandler{
		Reporter:     app.Reporter,
		JobService:   app.JobService,
		QueueManager: app.Queue,
	}

	cleanupJobListener := &messages.CleanupJobHandler{
		Reporter:   app.Reporter,
		JobService: app.JobService,
	}

	app.Queue.AddListener(createJobListener)
	app.Queue.AddListener(updateJobStatusListener)
	app.Queue.AddListener(cleanupJobListener)

	worker.PollMessages(ctx, app)
}

func (worker *MessagesProcessor) PollMessages(ctx context.Context, app *app.App) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			queueMessages := app.Queue.ReceiveMessages()
			go app.Queue.ProcessMessages(queueMessages)
		}
	}
}

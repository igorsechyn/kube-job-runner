package app

import (
	"context"

	"kube-job-runner/pkg/app/config"
	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/queue"
	"kube-job-runner/pkg/app/reporter"
)

type App struct {
	Reporter          *reporter.Reporter
	JobClient         job.Client
	JobService        *job.Service
	Queue             *queue.Manager
	Config            config.Config
}

type Component interface {
	Run(ctx context.Context, app *App)
}

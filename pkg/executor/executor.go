package executor

import (
	"context"

	"kube-job-runner/pkg/app"
	"kube-job-runner/pkg/app/config"
	"kube-job-runner/pkg/executor/web"
	"kube-job-runner/pkg/executor/worker"
)

const (
	WebServerComponent          = "web-server"
	MessagesProcessorComponent  = "messages-processor"
	JobStatusProcessorComponent = "job-status-processor"
	PodEventsProcessorComponent = "pod-events-processor"
)

func RunWithContext(ctx context.Context, app *app.App) {
	components := GetComponents(app.Config)
	for _, component := range components {
		go component.Run(ctx, app)
	}

	<-ctx.Done()
}

func GetComponents(config config.Config) []app.Component {
	if len(config.Components) == 0 {
		return []app.Component{&web.Server{}, &worker.MessagesProcessor{}, &worker.JobStatusProcessor{}, &worker.PodEventsProcessor{}}
	}

	components := make([]app.Component, len(config.Components), len(config.Components))
	for index, component := range config.Components {
		switch component {
		case WebServerComponent:
			components[index] = &web.Server{}
			break
		case MessagesProcessorComponent:
			components[index] = &worker.MessagesProcessor{}
			break
		case JobStatusProcessorComponent:
			components[index] = &worker.JobStatusProcessor{}
			break
		case PodEventsProcessorComponent:
			components[index] = &worker.PodEventsProcessor{}
			break
		}
	}
	return components
}

package runner

import (
	"context"

	"kube-job-runner/pkg/app"
	"kube-job-runner/pkg/app/config"
	"kube-job-runner/pkg/runner/web"
	"kube-job-runner/pkg/runner/worker"
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

	components := make([]app.Component, len(config.Components))
	for index, component := range config.Components {
		switch component {
		case WebServerComponent:
			components[index] = &web.Server{}
		case MessagesProcessorComponent:
			components[index] = &worker.MessagesProcessor{}
		case JobStatusProcessorComponent:
			components[index] = &worker.JobStatusProcessor{}
		case PodEventsProcessorComponent:
			components[index] = &worker.PodEventsProcessor{}
		}
	}
	return components
}

// +build unit

package runner_test

import (
	"testing"

	"kube-job-runner/pkg/app"
	"kube-job-runner/pkg/app/config"
	"kube-job-runner/pkg/runner"
	"kube-job-runner/pkg/runner/web"
	"kube-job-runner/pkg/runner/worker"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	description        string
	config             config.Config
	expectedComponents []app.Component
}

func Test_GetComponents(t *testing.T) {
	testCases := []testCase{
		{
			description:        "it should return all components, the config value is empty",
			config:             config.Config{Components: []string{}},
			expectedComponents: []app.Component{&web.Server{}, &worker.MessagesProcessor{}, &worker.JobStatusProcessor{}, &worker.PodEventsProcessor{}},
		},
		{
			description:        "it should return web server component",
			config:             config.Config{Components: []string{runner.WebServerComponent}},
			expectedComponents: []app.Component{&web.Server{}},
		},
		{
			description:        "it should return job status processor component",
			config:             config.Config{Components: []string{runner.JobStatusProcessorComponent}},
			expectedComponents: []app.Component{&worker.JobStatusProcessor{}},
		},
		{
			description:        "it should return messages processor component",
			config:             config.Config{Components: []string{runner.MessagesProcessorComponent}},
			expectedComponents: []app.Component{&worker.MessagesProcessor{}},
		},
		{
			description:        "it should return messages processor component",
			config:             config.Config{Components: []string{runner.PodEventsProcessorComponent}},
			expectedComponents: []app.Component{&worker.PodEventsProcessor{}},
		},
		{
			description:        "it should return multiple components, if they are specified in the config",
			config:             config.Config{Components: []string{runner.WebServerComponent, runner.JobStatusProcessorComponent}},
			expectedComponents: []app.Component{&web.Server{}, &worker.JobStatusProcessor{}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			components := runner.GetComponents(testCase.config)

			require.Equal(t, testCase.expectedComponents, components)
		})

	}
}

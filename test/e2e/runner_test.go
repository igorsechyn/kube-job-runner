// +build integration

package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	Image   = "igorsechyn/samplejob"
	Tag     = "1.0.0"
	Timeout = 10 * 60 * 1000
)

func whenJobIsFinished(image, tag string, timeout int64) (string, error) {
	client, err := NewRunnerClient()
	if err != nil {
		return "", err
	}
	jobUrl, err := client.SubmitJob(image, tag, timeout)

	if err != nil {
		return "", err
	}

	return getJobStatus(client, jobUrl, 1*time.Minute)
}

func getJobStatus(client *RunnerClient, jobURL string, timeout time.Duration) (string, error) {
	timeoutTicker := time.NewTimer(timeout)
	pollTicker := time.NewTicker(5000 * time.Millisecond)
	for {
		select {
		case <-pollTicker.C:
			status, err := client.GetJobStatus(jobURL)
			if err != nil {
				continue
			}

			if status == "InProgress" || status == "Acknowledged" {
				continue
			}

			return status, nil
		case <-timeoutTicker.C:
			return "", fmt.Errorf("job did not finish in time")
		}
	}
}

func TestRunner(t *testing.T) {
	t.Run("it should run a successful job", func(t *testing.T) {
		status, err := whenJobIsFinished(Image, Tag, Timeout)
		require.NoError(t, err)

		assert.Equal(t, "Succeeded", status)
	})

	t.Run("it should execute a job with missing docker image", func(t *testing.T) {
		status, err := whenJobIsFinished("igorsechyn/not-exist", Tag, Timeout)
		require.NoError(t, err)

		assert.Equal(t, "Failed", status)
	})
}

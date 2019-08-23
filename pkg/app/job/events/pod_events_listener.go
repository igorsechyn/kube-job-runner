package events

import (
	"kube-job-runner/pkg/app/data"
	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/queue"
	"kube-job-runner/pkg/app/reporter"
)

type Listener struct {
	JobClient    job.Client
	Reporter     *reporter.Reporter
	QueueManager *queue.Manager
}

func (processor *Listener) Process(event job.PodEvent) {
	if event.Status == "Failed" {
		jobID, err := processor.JobClient.GetJobIDForPod(event.PodID)
		if err != nil {
			processor.Reporter.Error("job.get.id.error", err, map[string]interface{}{})
			return
		}
		err = processor.QueueManager.SendJobStatusUpdateMessage(jobID, data.JobFailed, event.Reason)
		if err != nil {
			processor.Reporter.Error("send.event.message.error", err, map[string]interface{}{})
			return
		}
		err = processor.JobClient.DeleteEvent(event.ID)
		if err != nil {
			processor.Reporter.Error("delete.event.error", err, map[string]interface{}{})
		}
	}
}

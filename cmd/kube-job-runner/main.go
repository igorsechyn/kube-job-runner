package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"kube-job-runner/pkg/app"
	"kube-job-runner/pkg/app/config/viper"
	"kube-job-runner/pkg/app/data/pq"
	"kube-job-runner/pkg/app/id/uuid"
	"kube-job-runner/pkg/app/job"
	"kube-job-runner/pkg/app/job/k8s"
	"kube-job-runner/pkg/app/queue"
	"kube-job-runner/pkg/app/queue/sqs"
	"kube-job-runner/pkg/app/reporter"
	"kube-job-runner/pkg/app/reporter/zerolog"
	"kube-job-runner/pkg/app/time"
	"kube-job-runner/pkg/executor"
)

func main() {
	appConfig := viper.LoadConfig()
	appReporter := reporter.New(zerolog.NewJSONLogger())
	k8sClient, _ := k8s.NewClient(appConfig.Namespace)
	queueClient, err := sqs.NewClient(appConfig)
	if err != nil {
		appReporter.Error("queue.client.create.error", err, map[string]interface{}{})
		os.Exit(1)
		return
	}
	queueManager := queue.NewManager(queueClient, appReporter)
	// queueManager := queue.NewManager(queue.NewInMemoryQueueClient(), appReporter)
	store, err := pq.NewStore(appConfig)
	if err != nil {
		appReporter.Error("pq.client.create.error", err, map[string]interface{}{})
		os.Exit(1)
		return
	}
	// store := data.NewInMemoryStore()
	jobService := &job.Service{
		Reporter:    appReporter,
		JobClient:   k8sClient,
		IDGenerator: uuid.Generator{},
		Clock:       time.SystemClock{},
		DataStore:   store,
		Queue:       queueManager,
	}
	application := &app.App{
		Reporter:          appReporter,
		JobService:        jobService,
		Queue:             queueManager,
		Config:            appConfig,
		JobClient:         k8sClient,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cancelOnInterrupt(ctx, cancel)
	executor.RunWithContext(ctx, application)
}

func cancelOnInterrupt(ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
		case <-c:
			cancel()
		}
	}()
}

package k8s

import (
	"context"
	"fmt"
	"sync"
	"time"

	"kube-job-runner/pkg/app/job"

	"k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Client struct {
	clientset          *kubernetes.Clientset
	namespace          string
	statusListeners    []job.StatusListener
	podEventsListeners []job.PodEventListener
	mux                sync.Mutex
}

func NewClient(namespace string) (*Client, error) {
	clientset, err := getK8SClientset()
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset:          clientset,
		namespace:          namespace,
		statusListeners:    make([]job.StatusListener, 0),
		podEventsListeners: make([]job.PodEventListener, 0),
	}, nil
}

func getK8SClientset() (*kubernetes.Clientset, error) {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func (jobClient *Client) SubmitJob(job job.Job) (string, error) {
	k8sJob := v1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: job.JobName,
		},
		Spec: v1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: "Never",
					Containers: []corev1.Container{
						{
							Name:            "job-container",
							Image:           fmt.Sprintf("%v:%v", job.Image, job.Tag),
							ImagePullPolicy: "IfNotPresent",
						},
					},
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: job.JobName,
				},
			},
		},
	}
	_, err := jobClient.clientset.BatchV1().Jobs(jobClient.namespace).Create(&k8sJob)
	return job.JobName, err
}

func (jobClient *Client) DeleteJob(jobID string) error {
	policy := metav1.DeletePropagationForeground
	err := jobClient.clientset.BatchV1().Jobs(jobClient.namespace).Delete(jobID, &metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	return err
}

func (jobClient *Client) GetJobIDForPod(podID string) (string, error) {
	pod, err := jobClient.clientset.CoreV1().Pods(jobClient.namespace).Get(podID, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return pod.Labels["job-name"], nil
}

func (jobClient *Client) DeleteEvent(eventID string) error {
	err := jobClient.clientset.CoreV1().Events(jobClient.namespace).Delete(eventID, &metav1.DeleteOptions{})
	return err
}

func (jobClient *Client) AddJobStatusListener(listener job.StatusListener) {
	jobClient.mux.Lock()
	defer jobClient.mux.Unlock()
	jobClient.statusListeners = append(jobClient.statusListeners, listener)
}

func (jobClient *Client) AddPodEventsListener(listener job.PodEventListener) {
	jobClient.mux.Lock()
	defer jobClient.mux.Unlock()
	jobClient.podEventsListeners = append(jobClient.podEventsListeners, listener)
}

func (jobClient *Client) WatchJobs(ctx context.Context) {
	factory := informers.NewSharedInformerFactoryWithOptions(
		jobClient.clientset,
		5*time.Second,
		informers.WithNamespace(jobClient.namespace),
	)
	informer := factory.Batch().V1().Jobs().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			jobObject := obj.(*v1.Job)
			jobClient.callStatusListeners(jobObject)
		},
		DeleteFunc: func(obj interface{}) {
			jobObject := obj.(*v1.Job)
			jobClient.callStatusListeners(jobObject)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			jobObject := newObj.(*v1.Job)
			jobClient.callStatusListeners(jobObject)
		},
	})
	informer.Run(ctx.Done())
}

func (jobClient *Client) WatchEvents(ctx context.Context) {
	factory := informers.NewSharedInformerFactoryWithOptions(
		jobClient.clientset,
		5*time.Second,
		informers.WithNamespace(jobClient.namespace),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = "involvedObject.kind=Pod,reason=Failed"
		}),
	)

	informer := factory.Core().V1().Events().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			eventObject := obj.(*corev1.Event)
			jobClient.callEventListeners(eventObject)
		},
		DeleteFunc: func(obj interface{}) {
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
		},
	})
	informer.Run(ctx.Done())
}

func (jobClient *Client) callStatusListeners(k8sJob *v1.Job) {
	jobStatus := job.Status{
		JobID:               k8sJob.Name,
		FailedContainers:    k8sJob.Status.Failed,
		SucceededContainers: k8sJob.Status.Succeeded,
		RunningContainers:   k8sJob.Status.Active,
	}
	for _, listener := range jobClient.statusListeners {
		listener.Process(jobStatus)
	}
}

func (jobClient *Client) callEventListeners(k8sEvent *corev1.Event) {
	event := job.PodEvent{
		ID:     k8sEvent.Name,
		Status: k8sEvent.Reason,
		Reason: k8sEvent.Message,
		PodID:  k8sEvent.InvolvedObject.Name,
	}
	for _, listener := range jobClient.podEventsListeners {
		listener.Process(event)
	}
}

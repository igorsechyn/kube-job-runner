package queue

import (
	"encoding/json"
	"sync"

	"kube-job-runner/pkg/app/reporter"
)

type CreateJobMessage struct {
	JobID string `json:"jobID"`
	Type  string `json:"type"`
}

type JobStatusUpdateMessage struct {
	JobID  string `json:"jobID"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Reason string `json:"reason,omitempty"`
}

type JobCleanupMessage struct {
	JobID string `json:"jobID"`
	Type  string `json:"type"`
}

type MessageListener interface {
	Accept(messageType string) bool
	Process(message string) error
}

type Manager struct {
	client    Client
	listeners []MessageListener
	mux       sync.Mutex
	reporter  *reporter.Reporter
}

func NewManager(client Client, reporter *reporter.Reporter) *Manager {
	return &Manager{
		client:    client,
		listeners: make([]MessageListener, 0),
		reporter:  reporter,
	}
}

func (queue *Manager) AddListener(listener MessageListener) {
	queue.mux.Lock()
	defer queue.mux.Unlock()
	queue.listeners = append(queue.listeners, listener)
}

func (queue *Manager) SendCreateJobMessage(jobId string) error {
	request := CreateJobMessage{
		JobID: jobId,
		Type:  CreateMessageType,
	}

	return queue.sendMessage(request)
}

func (queue *Manager) SendJobStatusUpdateMessage(jobID string, status string, reason string) error {
	update := JobStatusUpdateMessage{
		JobID:  jobID,
		Status: status,
		Type:   JobStatusUpdateType,
		Reason: reason,
	}

	return queue.sendMessage(update)
}

func (queue *Manager) SendJobCleanupMessage(jobID string) error {
	update := JobCleanupMessage{
		JobID: jobID,
		Type:  JobCleanupType,
	}

	return queue.sendMessage(update)
}

func (queue *Manager) sendMessage(body interface{}) error {
	message, err := json.Marshal(body)

	if err != nil {
		return err
	}

	err = queue.client.SendMessage(string(message))
	return err
}

func (queue *Manager) ProcessMessages(messages []Message) {
	for _, message := range messages {
		messageType, err := GetMessageType(message.Body)
		if err != nil {
			continue
		}

		for _, listener := range queue.listeners {
			if listener.Accept(messageType) {
				err := listener.Process(message.Body)
				if err == nil {
					message.Delete()
				}
			}
		}
	}
}

func (queue *Manager) ReceiveMessages() []Message {
	messages, err := queue.client.ReceiveMessages()
	if err != nil {
		queue.reporter.Error("receive.messages", err, map[string]interface{}{})
	}
	return messages
}

func GetMessageType(message string) (string, error) {
	var messageType struct {
		Type string `json:"type"`
	}
	err := json.Unmarshal([]byte(message), &messageType)
	if err != nil {
		return "", err
	}
	return messageType.Type, nil
}

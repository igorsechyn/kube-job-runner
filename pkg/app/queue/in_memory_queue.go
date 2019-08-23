package queue

import (
	"sync"
)

func NewInMemoryQueueClient() Client {
	return &InMemoryQueueClient{
		messagesList: []Message{},
	}
}

type InMemoryQueueClient struct {
	messagesList []Message
	mux          sync.Mutex
}

func (store *InMemoryQueueClient) SendMessage(message MessageBody) error {
	store.mux.Lock()
	defer store.mux.Unlock()
	store.messagesList = append(store.messagesList, Message{Body: message, Delete: func() {

	}})
	return nil
}

func (store *InMemoryQueueClient) ReceiveMessages() ([]Message, error) {
	store.mux.Lock()
	defer store.mux.Unlock()
	allMessages := store.messagesList
	store.messagesList = make([]Message, 0)
	return allMessages, nil
}

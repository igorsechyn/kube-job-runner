package queue

type MessageBody string

const (
	CreateMessageType   = "CREATE_JOB"
	JobStatusUpdateType = "JOB_STATUS_UPDATE"
	JobCleanupType      = "CLEANUP_JOB"
	PodEventType        = "POD_EVENT"
)

type Message struct {
	Body   MessageBody
	Delete func()
}

type Client interface {
	SendMessage(message MessageBody) error
	ReceiveMessages() ([]Message, error)
}

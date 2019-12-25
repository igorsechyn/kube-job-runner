package queue

type MessageBody string

const (
	CreateMessageType   = "CREATE_JOB"
	JobStatusUpdateType = "JOB_STATUS_UPDATE"
	JobCleanupType      = "CLEANUP_JOB"
)

type Message struct {
	Body   string
	Delete func()
}

type Client interface {
	SendMessage(message string) error
	ReceiveMessages() ([]Message, error)
}

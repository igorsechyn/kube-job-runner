package data

const (
	JobAcknowledged = "Acknowledged"
	JobInProgress   = "InProgress"
	JobSucceeded    = "Succeeded"
	JobFailed       = "Failed"
)

type Document struct {
	Image     string
	Tag       string
	JobID     string
	Status    string
	Timestamp int64
	Timeout   int64
}

type Store interface {
	PutDocument(document Document) error
	GetDocuments(id string) ([]Document, error)
}

package config

type Config struct {
	Environment        string
	PgHost             string
	PgUser             string
	PgPassword         string
	PgDatabase         string
	PgPort             int
	SQSQueueUrl        string
	SQSQueueRegion     string
	SQSQueueName       string
	SQSWaitTimeSeconds int
	Namespace          string
	Components         []string
}

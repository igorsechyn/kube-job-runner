package sqs

import (
	"kube-job-runner/pkg/app/config"
	"kube-job-runner/pkg/app/queue"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Client struct {
	queue           *sqs.SQS
	queueUrl        *string
	waitTimeSeconds int
}

func NewClient(config config.Config) (*Client, error) {
	awsSession, err := session.NewSession(getAWSConfig(config))
	if err != nil {
		return nil, err
	}
	sqsQueue := sqs.New(awsSession)
	return &Client{queue: sqsQueue, queueUrl: aws.String(config.SQSQueueUrl), waitTimeSeconds: config.SQSWaitTimeSeconds}, nil
}

func getAWSConfig(config config.Config) *aws.Config {
	awsConfig := &aws.Config{
		Region: aws.String(config.SQSQueueRegion),
	}

	if config.Environment == "local" {
		awsConfig.Credentials = credentials.NewStaticCredentialsFromCreds(credentials.Value{AccessKeyID: "x", SecretAccessKey: "x"})
		awsConfig.Endpoint = aws.String(config.SQSQueueUrl)
		awsConfig.DisableSSL = aws.Bool(true)
	}

	return awsConfig
}

func (client *Client) SendMessage(message string) error {
	sqsMessage := &sqs.SendMessageInput{
		MessageBody: aws.String(string(message)),
		QueueUrl:    client.queueUrl,
	}

	_, err := client.queue.SendMessage(sqsMessage)
	return err
}

func (client *Client) ReceiveMessages() ([]queue.Message, error) {
	receiveMessagesRequest := &sqs.ReceiveMessageInput{
		QueueUrl:            client.queueUrl,
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(int64(client.waitTimeSeconds)),
	}
	response, err := client.queue.ReceiveMessage(receiveMessagesRequest)
	if err != nil {
		return []queue.Message{}, err
	}

	messages := client.toMessages(response.Messages)
	return messages, nil
}

func (client *Client) toMessages(sqsMessages []*sqs.Message) []queue.Message {
	messages := make([]queue.Message, len(sqsMessages))
	for index, sqsMessage := range sqsMessages {
		deleteFunction := client.getDeleteFunction(sqsMessage)
		message := queue.Message{
			Body:   *sqsMessage.Body,
			Delete: deleteFunction,
		}
		messages[index] = message
	}
	return messages
}

func (client *Client) getDeleteFunction(message *sqs.Message) func() {
	return func() {
		client.queue.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      client.queueUrl,
			ReceiptHandle: message.ReceiptHandle,
		})
	}
}

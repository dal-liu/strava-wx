package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSClient struct {
	svc *sqs.Client
}

func (c SQSClient) Send(ctx context.Context, messageBody, queueUrl string) error {
	_, err := c.svc.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    aws.String(queueUrl),
	})
	return err
}

func CreateClient(ctx context.Context) (SQSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return SQSClient{}, err
	}
	return SQSClient{sqs.NewFromConfig(cfg)}, nil
}

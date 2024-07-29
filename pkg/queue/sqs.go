package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var svc *sqs.Client

func CreateClient(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	svc = sqs.NewFromConfig(cfg)
	return nil
}

func Send(ctx context.Context, messageBody, queueUrl string) error {
	_, err := svc.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    aws.String(queueUrl),
	})
	if err != nil {
		return err
	}
	return nil
}

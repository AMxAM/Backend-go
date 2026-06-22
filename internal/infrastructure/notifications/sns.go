package notifications

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSService struct {
	client   *sns.Client
	topicArn string
}

func NewSNSService() (*SNSService, error) {

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
	)
	if err != nil {
		return nil, err
	}

	return &SNSService{
		client: sns.NewFromConfig(cfg),
		topicArn: os.Getenv(
			"SNS_TOPIC_ARN",
		),
	}, nil
}

func (s *SNSService) Publish(
	message string,
) error {

	_, err := s.client.Publish(
		context.Background(),
		&sns.PublishInput{
			TopicArn: &s.topicArn,
			Message:  &message,
		},
	)

	return err
}
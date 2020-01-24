package ddd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func (c *SnsClient) PublishMessage(message string) (*sns.PublishOutput, error) {
	return c.Publish(&sns.PublishInput{
		Message:  &message,
		Subject:  &c.subject,
		TopicArn: &c.topic,
	})
}

type SnsClient struct {
	*sns.SNS
	subject string
	topic   string
}

func NewSnsClient(topic, subject string) (*SnsClient, error) {
	snsClient, err := snsClient("us-west-2")
	if err != nil {
		return nil, err
	}
	return &SnsClient{
		SNS:     snsClient,
		subject: subject,
		topic:   topic,
	}, nil
}

func snsClient(region string) (*sns.SNS, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	return sns.New(sess), nil
}

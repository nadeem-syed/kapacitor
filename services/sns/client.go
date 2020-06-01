package sns

import (
	"fmt"
	//"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SNSConfig struct {
	Region    string
	AccessKey string
	SecretKey string
}

type Client interface {
	Version() (string, error)
	Publish(*sns.PublishInput) error
}

type snsStruct struct {
	client *sns.SNS
}

func new(c SNSConfig) (Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	svc := sns.New(sess)
	return &snsStruct{
		client: svc,
	}, nil
}

func (svc *snsStruct) Version() (string, error) {
	return "2013-10-15", nil
}

func (svc *snsStruct) Publish(message *sns.PublishInput) error {
	return nil
}

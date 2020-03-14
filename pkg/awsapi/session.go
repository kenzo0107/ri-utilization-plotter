package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// NewSession : new session
func NewSession(c *sts.Credentials, region string) *session.Session {
	config := aws.Config{
		Region: aws.String(region),
	}
	return session.Must(session.NewSession(&config))
}

package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// NewSession : new session
func NewSession(c *sts.Credentials, region string) *session.Session {
	var config aws.Config
	if c != nil {
		// credentials exist ... except AWS Account which assume role doesn't belong to
		creds := credentials.NewStaticCredentials(
			*c.AccessKeyId,
			*c.SecretAccessKey,
			*c.SessionToken,
		)
		config = aws.Config{
			Region:      aws.String(region),
			Credentials: creds,
		}
	} else {
		// only AWS Account which assume role belongs to
		config = aws.Config{
			Region: aws.String(region),
		}
	}
	return session.Must(session.NewSession(&config))
}

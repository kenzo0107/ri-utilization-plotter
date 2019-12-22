package awsapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

const (
	durationSeconds    = 900
	fmtRoleArn         = "arn:aws:iam::%s:role/stsMonitor"
	fmtRoleSessionName = "RIof%s"
)

// STSIface : sts interface
type STSIface interface {
	GetAssumeRoleCredentials(accountID string) (*sts.Credentials, error)
}

// STSInstance : sts instance
type STSInstance struct {
	client stsiface.STSAPI
}

// NewSTSAssumeRoler ... generate a new sts client
func NewSTSAssumeRoler(client stsiface.STSAPI) STSIface {
	return &STSInstance{client: client}
}

// GetAssumeRoleCredentials ... get assume role's credentilas
func (d *STSInstance) GetAssumeRoleCredentials(accountID string) (*sts.Credentials, error) {
	roleArn := fmt.Sprintf(fmtRoleArn, accountID)
	roleSessionName := fmt.Sprintf(fmtRoleSessionName, accountID)
	p := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(durationSeconds),
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(roleSessionName),
	}
	r, err := d.client.AssumeRole(p)
	return r.Credentials, err
}

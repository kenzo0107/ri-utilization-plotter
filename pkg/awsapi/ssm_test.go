package awsapi

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/google/go-cmp/cmp"
)

type mockSSMClient struct {
	ssmiface.SSMAPI

	Output *ssm.GetParametersOutput
	Error  error
}

func (m *mockSSMClient) GetParameters(*ssm.GetParametersInput) (*ssm.GetParametersOutput, error) {
	return m.Output, m.Error
}

func TestGetSSMParameters(t *testing.T) {
	m := NewSSMClient(&mockSSMClient{
		Output: &ssm.GetParametersOutput{
			Parameters: []*ssm.Parameter{
				&ssm.Parameter{
					Name:             aws.String("datadog_api_key"),
					Type:             aws.String("SecureString"),
					Value:            aws.String("hogehoge"),
					Version:          aws.Int64(1),
					LastModifiedDate: aws.Time(time.Now()),
					ARN:              aws.String("arn:aws:ssm:ap-northeast-1:123456789012:parameter/datadog_api_key"),
				},
				&ssm.Parameter{
					Name:             aws.String("datadog_app_key"),
					Type:             aws.String("SecureString"),
					Value:            aws.String("mogemoge"),
					Version:          aws.Int64(1),
					LastModifiedDate: aws.Time(time.Now()),
					ARN:              aws.String("arn:aws:ssm:ap-northeast-1:123456789012:parameter/datadog_app_key"),
				},
			},
			InvalidParameters: []*string{},
		},
		Error: nil,
	})
	keys := []string{"datadog_api_key", "datadog_app_key"}
	s, err := m.GetSSMParameters(keys)
	if err != nil {
		t.Error(err)
	}
	datadogAPIKey := s["datadog_api_key"]
	datadogAppKey := s["datadog_app_key"]

	if diff := cmp.Diff(datadogAPIKey, "hogehoge"); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if diff := cmp.Diff(datadogAppKey, "mogemoge"); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
}

func TestGetSSMParametersFailed(t *testing.T) {
	m := NewSSMClient(&mockSSMClient{
		Output: &ssm.GetParametersOutput{
			Parameters: []*ssm.Parameter{
				&ssm.Parameter{
					Name:             aws.String("datadog_api_key"),
					Type:             aws.String("SecureString"),
					Value:            aws.String("hogehoge"),
					Version:          aws.Int64(1),
					LastModifiedDate: aws.Time(time.Now()),
					ARN:              aws.String("arn:aws:ssm:ap-northeast-1:123456789012:parameter/datadog_api_key"),
				},
			},
			InvalidParameters: []*string{},
		},
		Error: errors.New("error occured"),
	})
	keys := []string{"datadog_api_key", "datadog_app_key"}
	_, err := m.GetSSMParameters(keys)
	if err == nil {
		t.Errorf("wrong result : err is nil : %s", err.Error())
	}
}

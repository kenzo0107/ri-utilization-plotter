package awsapi

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/google/go-cmp/cmp"
)

type mockSTSClient struct {
	stsiface.STSAPI
}

func (m *mockSTSClient) AssumeRole(*sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	return &sts.AssumeRoleOutput{
		AssumedRoleUser: &sts.AssumedRoleUser{
			Arn:           aws.String("arn:aws:iam::123456789012:role/hoge"),
			AssumedRoleId: aws.String("1234567"),
		},
		Credentials: &sts.Credentials{
			AccessKeyId:     aws.String("hogehoge"),
			SecretAccessKey: aws.String("mogemoge"),
		},
	}, nil
}

func TestGetAssumeRoleCredentials(t *testing.T) {
	m := NewSTSAssumeRoler(&mockSTSClient{})
	actual, err := m.GetAssumeRoleCredentials("223456789012")
	if err != nil {
		t.Error(err)
	}

	expected := &sts.Credentials{
		AccessKeyId:     aws.String("hogehoge"),
		SecretAccessKey: aws.String("mogemoge"),
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
}

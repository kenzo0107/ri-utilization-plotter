package awsapi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/google/go-cmp/cmp"
)

func TestNewSession(t *testing.T) {
	region := "ap-northeast-1"
	sess := NewSession(&sts.Credentials{
		AccessKeyId:     aws.String("hogehoge"),
		SecretAccessKey: aws.String("mogemoge"),
		SessionToken:    aws.String("barbar"),
	}, region)

	actual, err := sess.Config.Credentials.Get()
	if err != nil {
		t.Error(err)
	}

	expected := credentials.Value{
		AccessKeyID:     "hogehoge",
		SecretAccessKey: "mogemoge",
		SessionToken:    "barbar",
		ProviderName:    "StaticProvider",
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
}

func TestNewSessionWithoutCredentials(t *testing.T) {
	credentialsPath := filepath.Join("..", "..", "testdata", "credentials")
	if err := os.Setenv("AWS_SHARED_CREDENTIALS_FILE", credentialsPath); err != nil {
		t.Error("error occured in os.Setenv(\"AWS_SHARED_CREDENTIALS_FILE\")")
	}
	region := "ap-northeast-1"
	sess := NewSession(nil, region)
	_, err := sess.Config.Credentials.Get()
	if err == nil {
		t.Errorf("wrong result : err is nil. \n%s", err.Error())
	}
}

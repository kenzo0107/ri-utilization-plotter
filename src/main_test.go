package main

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadConfigureAWSAccount(t *testing.T) {
	credentialsPath := filepath.Join("..", "testdata", "awsaccount.yml")
	awsAccounts, err := loadConfigureAWSAccount(credentialsPath)
	if err != nil {
		t.Error(err)
	}

	expected := []awsAccount{
		awsAccount{
			ID:      "123456789012",
			Profile: "original",
			Default: true,
		},
		awsAccount{
			ID:      "923456789012",
			Profile: "hoge",
			Default: false,
		},
		awsAccount{
			ID:      "823456789012",
			Profile: "moge",
			Default: false,
		},
	}

	if diff := cmp.Diff(expected, awsAccounts); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
}

func TestLoadConfigureAWSAccountFailed(t *testing.T) {
	_, err := loadConfigureAWSAccount("")
	if err == nil {
		t.Errorf("wrong result : err is nil. \n%s", err.Error())
	}
}

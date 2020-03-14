package awsapi

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

func TestNewSessionWithoutCredentials(t *testing.T) {
	sess := NewSession(nil, endpoints.ApNortheast1RegionID)
	_, err := sess.Config.Credentials.Get()
	if err != nil {
		t.Error(err)
	}
}

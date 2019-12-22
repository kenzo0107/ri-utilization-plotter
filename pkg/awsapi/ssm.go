package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// SSMIface : -
type SSMIface interface {
	GetSSMParameters(keys []string) (map[string]string, error)
}

// SSMInstance : ssm instance
type SSMInstance struct {
	client ssmiface.SSMAPI
}

// NewSSMClient ... generate a new ssm client
func NewSSMClient(client ssmiface.SSMAPI) SSMIface {
	return &SSMInstance{
		client: client,
	}
}

// GetSSMParameters ... get values from ssm parameter store
func (d *SSMInstance) GetSSMParameters(keys []string) (map[string]string, error) {
	names := []*string{}
	for _, k := range keys {
		names = append(names, aws.String(k))
	}

	ssmParameters := &ssm.GetParametersInput{
		Names:          names,
		WithDecryption: aws.Bool(true),
	}

	r, err := d.client.GetParameters(ssmParameters)
	if err != nil {
		return nil, err
	}

	s := make(map[string]string)
	for _, p := range r.Parameters {
		s[*p.Name] = *p.Value
	}
	return s, nil
}

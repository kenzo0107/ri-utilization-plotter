package configs

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/caarlos0/env"
	"github.com/kenzo0107/ri-utilization-plotter/pkg/awsapi"
	"github.com/pkg/errors"
)

// Envs : environment values
var Envs envParameters

type envParameters struct {
	DatadogAPIKeyName string `env:"DD_API_KEY_NAME" envDefault:"datadog_api_key"`
	DatadogAppKeyName string `env:"DD_APP_KEY_NAME" envDefault:"datadog_app_key"`
	TagKey            string `env:"TAG_KEY" envDefault:"account"`
	TagVal            string `env:"TAG_VAL" envDefault:"yourproject"`
	AWSRegionID       string `env:"AWS_REGION"`
}

// Session : session
var Session *session.Session

// Secrets : secrets
var Secrets SecretParameters

// SecretParameters : secret parameters
type SecretParameters struct {
	DatadogAPIKey string
	DatadogAppKey string
}

func init() {
	cfg := envParameters{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed on env.Parse(&cfg)"))
	}
	Envs = cfg

	Session = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(Envs.AWSRegionID),
	}))

	ssmClient := awsapi.NewSSMClient(ssm.New(Session))
	s, err := ssmClient.GetSSMParameters([]string{
		Envs.DatadogAPIKeyName,
		Envs.DatadogAppKeyName,
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed on ssmClient.GetSSMParameters"))
	}
	ddAPIKey := s[Envs.DatadogAPIKeyName]
	ddAPPKey := s[Envs.DatadogAppKeyName]

	Secrets.DatadogAPIKey = ddAPIKey
	Secrets.DatadogAppKey = ddAPPKey
}

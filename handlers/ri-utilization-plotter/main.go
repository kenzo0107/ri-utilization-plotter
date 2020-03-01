package main

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	"github.com/zorkian/go-datadog-api"
	"gopkg.in/yaml.v2"

	"github.com/kenzo0107/ri-utilization-plotter/pkg/awsapi"
	"github.com/kenzo0107/ri-utilization-plotter/pkg/utility"
)

const (
	region = endpoints.ApNortheast1RegionID

	// expectedCountOfAWSAccounts : この値以下になるだろうと予測される AWS Account 数
	//   configs/awsaccount.yml から slice に載せる際の初期容量設定時に利用します。
	//   容量拡張のコストは然程気にしなくても良いかとは思いますが、念の為、
	//   容量の数の検討がつく場合は設定しておく基本に沿って以下設定。
	expectedCountOfAWSAccounts = 80
)

type awsAccount struct {
	ID      string
	Profile string
	Default bool
}

var (
	services = []string{
		"Amazon Elastic Compute Cloud - Compute",
		"Amazon Relational Database Service",
		"Amazon ElastiCache",
		"Amazon Redshift",
		"Amazon Elasticsearch Service",
	}
	unixTime             float64
	startDay             string
	endDay               string
	datadogClient        *datadog.Client
	ssmDatadogAPIKeyName string
	ssmDatadogAPPKeyName string
	configAWSAccount     string = filepath.Join("configs", "awsaccount.yml")
)

func init() {
	ssmDatadogAPIKeyName = os.Getenv("DD_API_KEY_NAME")
	ssmDatadogAPPKeyName = os.Getenv("DD_APP_KEY_NAME")

	now := time.Now()
	unixTime = float64(now.Unix())
	endDay = now.Format("2006-01-02")
	// GetReservationUtilization 呼び出し時に最低でも 2 日前を指定する必要がある
	startDay = now.AddDate(0, 0, -2).Format("2006-01-02")
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}))

	ssmClient := awsapi.NewSSMClient(ssm.New(sess))
	s, err := ssmClient.GetSSMParameters([]string{ssmDatadogAPIKeyName, ssmDatadogAPPKeyName})
	if err != nil {
		return errors.Wrap(err, "failed ssmClient.GetSSMParameters")
	}
	ddAPIKey := s[ssmDatadogAPIKeyName]
	ddAPPKey := s[ssmDatadogAPPKeyName]

	datadogClient = datadog.NewClient(ddAPIKey, ddAPPKey)

	stsClient := awsapi.NewSTSAssumeRoler(sts.New(sess))

	awsAccounts, err := loadConfigureAWSAccount(configAWSAccount)
	if err != nil {
		return errors.Wrap(err, "loadConfigureAWSAccount")
	}

	for _, account := range awsAccounts {
		if !account.Default {
			creds, err := stsClient.GetAssumeRoleCredentials(account.ID)
			if err != nil {
				return errors.Wrap(err, "failed stsClient.GetAssumeRoleCredentials(account.ID)")
			}

			// temporary credentials for a role
			sess = awsapi.NewSession(creds, region)
		}

		costexplorerClient := awsapi.NewCostexplorer(costexplorer.New(sess))

		for _, service := range services {
			// RI Utilization Coverage
			utilPct, err := costexplorerClient.FetchRIUtilizationPercentage(service, startDay, endDay)
			if err != nil {
				return errors.Wrap(err, "on costexplorerClient.FetchRIUtilizationPercentage.")
			}

			if utilPct == "" {
				// utilPct == "" means that you do not use the service
				continue
			}

			utilPercentage, _ := strconv.ParseFloat(utilPct, 64)
			if err := postMetricRIUtil(service, utilPercentage, account); err != nil {
				return errors.Wrap(err, "on postMetricRIUtil.")
			}

			// RI Cost Coverage
			coveragePcts, err := costexplorerClient.FetchRICoveragePercentage(service, startDay, endDay)
			if err != nil {
				return errors.Wrap(err, "on costexplorerClient.FetchRICoveragePercentage.")
			}

			if len(coveragePcts) == 0 {
				continue
			}

			for _, g := range coveragePcts {
				// post metric of RI coverage to Datadog
				if err := postMetricRICoverage(service, g, account); err != nil {
					return errors.Wrap(err, "on postMetricRICoverage.")
				}
			}
		}
	}
	return nil
}

// postMetricRIUtil : post metric of RI utilization to Datadog
func postMetricRIUtil(service string, utilPercentage float64, account awsAccount) error {
	metric := "aws.ri.utilization"
	typeDatadog := "guage"

	tags := []string{
		utility.CombineStrings([]string{"account:", account.Profile}),
		account.Profile,
		utility.CombineStrings([]string{"service:", service}),
	}

	series := []datadog.Metric{
		{
			Metric: &metric,
			Points: []datadog.DataPoint{
				{&unixTime, &utilPercentage},
			},
			Type: &typeDatadog,
			Host: &account.Profile,
			Tags: tags,
		},
	}
	return datadogClient.PostMetrics(series)
}

// postMetricRICoverage : post metric of RI utilization to Datadog
func postMetricRICoverage(service string, g *costexplorer.ReservationCoverageGroup, account awsAccount) error {
	metric := "aws.ri.coverage"
	typeDatadog := "guage"

	// string to float64
	pct, _ := strconv.ParseFloat(*g.Coverage.CoverageHours.CoverageHoursPercentage, 64)

	tags := []string{
		utility.CombineStrings([]string{"instance_type:", *g.Attributes["instanceType"]}),
		utility.CombineStrings([]string{"region:", *g.Attributes["region"]}),
		utility.CombineStrings([]string{"account:", account.Profile}),
		account.Profile,
		utility.CombineStrings([]string{"service:", service}),
	}

	series := []datadog.Metric{
		{
			Metric: &metric,
			Points: []datadog.DataPoint{
				{&unixTime, &pct},
			},
			Type: &typeDatadog,
			Host: &account.Profile,
			Tags: tags,
		},
	}

	// post custom metric to Datadog
	return datadogClient.PostMetrics(series)
}

func readOnStruct(fileBuffer []byte) ([]awsAccount, error) {
	data := make([]awsAccount, 0, expectedCountOfAWSAccounts)
	err := yaml.Unmarshal(fileBuffer, &data)
	return data, err
}

func loadConfigureAWSAccount(file string) ([]awsAccount, error) {
	buf, err := ioutil.ReadFile(filepath.Clean(file))
	if err != nil {
		return nil, err
	}

	awsAccounts, err := readOnStruct(buf)
	return awsAccounts, err
}

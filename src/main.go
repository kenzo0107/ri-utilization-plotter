package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/zorkian/go-datadog-api"
	"gopkg.in/yaml.v2"

	"github.com/kenzo0107/ri-utilization-plotter/pkg/awsapi"
	"github.com/kenzo0107/ri-utilization-plotter/pkg/utility"
)

const region = "ap-northeast-1"

type awsAccount struct {
	ID      string
	Profile string
	Default bool
}

var (
	services = [5]string{
		"Amazon Elastic Compute Cloud - Compute",
		"Amazon Relational Database Service",
		"Amazon ElastiCache",
		"Amazon Redshift",
		"Amazon Elasticsearch Service",
	}
	unixTime         float64
	startDay         string
	endDay           string
	datadogClient    *datadog.Client
	ssmDatadogAPIKey string
	ssmDatadogAPPKey string
)

func init() {
	ssmDatadogAPIKey = os.Getenv("DD_API_KEY")
	ssmDatadogAPPKey = os.Getenv("DD_APP_KEY")

	now := time.Now()
	unixTime = float64(now.Unix())
	endDay = now.Format("2006-01-02")
	startDay = now.AddDate(0, 0, -2).Format("2006-01-02")
}

func readOnStruct(fileBuffer []byte) ([]awsAccount, error) {
	data := make([]awsAccount, 30)
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

func handler(ctx context.Context) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	ssmClient := awsapi.NewSSMClient(ssm.New(sess))
	s, err := ssmClient.GetSSMParameters([]string{ssmDatadogAPIKey, ssmDatadogAPPKey})
	if err != nil {
		return "ng", err
	}
	ddAPIKey := s[ssmDatadogAPIKey]
	ddAPPKey := s[ssmDatadogAPPKey]

	datadogClient = datadog.NewClient(ddAPIKey, ddAPPKey)

	stsClient := awsapi.NewSTSAssumeRoler(sts.New(sess))

	configAWSAccount := filepath.Join("configs", "awsaccount.yml")
	awsAccounts, err := loadConfigureAWSAccount(configAWSAccount)
	if err != nil {
		return "ng", err
	}

	for _, account := range awsAccounts {
		if !account.Default {
			creds, err := stsClient.GetAssumeRoleCredentials(account.ID)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			// temporary credentials for a role
			sess = awsapi.NewSession(creds, region)
		}

		costexplorerClient := awsapi.NewCostexplorer(costexplorer.New(sess))

		for _, service := range services {
			// RI Utilization Coverage
			if utilPercentage, err := costexplorerClient.FetchRIUtilizationPercentage(service, startDay, endDay); err == nil {
				log.Println("utilPercentage: ", utilPercentage)
				postMetricRIUtil(service, utilPercentage, account)
			}

			// RI Cost Coverage
			if costCoverages, err := costexplorerClient.FetchRICoveragePercentage(service, startDay, endDay); err == nil {
				if len(costCoverages) == 0 {
					continue
				}

				// post metric of RI coverage to Datadog
				postMetricRICoverage(service, costCoverages, account)
			}
		}
	}
	return "ok", nil
}

func postMetricRIUtil(service string, utilPercentage float64, account awsAccount) {
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
	err := datadogClient.PostMetrics(series)
	if err != nil {
		log.Println(err)
	}
}

func postMetricRICoverage(service string, costCoverages []awsapi.CostCoverage, account awsAccount) {
	metric := "aws.ri.coverage"
	typeDatadog := "guage"

	for _, v := range costCoverages {
		tags := []string{
			utility.CombineStrings([]string{"instance_type:", v.InstanceType}),
			utility.CombineStrings([]string{"region:", v.Region}),
			utility.CombineStrings([]string{"account:", account.Profile}),
			account.Profile,
			utility.CombineStrings([]string{"service:", service}),
		}

		series := []datadog.Metric{
			{
				Metric: &metric,
				Points: []datadog.DataPoint{
					{&unixTime, &v.CoverageHoursPercentage},
				},
				Type: &typeDatadog,
				Host: &account.Profile,
				Tags: tags,
			},
		}

		// post custom metric to Datadog
		err := datadogClient.PostMetrics(series)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	lambda.Start(handler)
}

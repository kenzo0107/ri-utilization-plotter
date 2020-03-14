package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/pkg/errors"
	"github.com/zorkian/go-datadog-api"

	"github.com/kenzo0107/ri-utilization-plotter/configs"
	"github.com/kenzo0107/ri-utilization-plotter/pkg/awsapi"
	"github.com/kenzo0107/ri-utilization-plotter/pkg/utility"
)

var (
	services = []string{
		"Amazon Elastic Compute Cloud - Compute",
		"Amazon Relational Database Service",
		"Amazon ElastiCache",
		"Amazon Redshift",
		"Amazon Elasticsearch Service",
	}
	unixTime      float64
	startDay      string
	endDay        string
	datadogClient *datadog.Client
)

func init() {
	now := time.Now()
	unixTime = float64(now.Unix())
	endDay = now.Format("2006-01-02")
	// GetReservationUtilization 呼び出し時に最低でも 2 日前を指定する必要がある
	startDay = now.AddDate(0, 0, -2).Format("2006-01-02")

	datadogClient = datadog.NewClient(
		configs.Secrets.DatadogAPIKey,
		configs.Secrets.DatadogAppKey,
	)
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(configs.Envs.AWSRegionID),
	}))

	costexplorerClient := awsapi.NewCostexplorer(costexplorer.New(sess))

	for _, service := range services {
		// RI Utilization
		utilPct, errRIUtil := costexplorerClient.FetchRIUtilizationPercentage(service, startDay, endDay)
		if errRIUtil != nil {
			return errors.Wrap(
				errRIUtil,
				fmt.Sprintf("service: %s on costexplorerClient.FetchRIUtilizationPercentage", service),
			)
		}

		if utilPct != "" {
			// utilPct == "" means that you do not use the service
			utilPercentage, _ := strconv.ParseFloat(utilPct, 64)
			if err := postMetricRIUtil(service, utilPercentage, configs.Envs.TagKey, configs.Envs.TagVal); err != nil {
				return errors.Wrap(err, "on postMetricRIUtil.")
			}
		}

		// RI Coverage
		coveragePcts, errRICov := costexplorerClient.FetchRICoveragePercentage(service, startDay, endDay)
		if errRICov != nil {
			return errors.Wrap(
				errRICov,
				fmt.Sprintf("service: %s on costexplorerClient.FetchRICoveragePercentage", service),
			)

		}

		for _, g := range coveragePcts {
			// post metric of RI coverage to Datadog
			if err := postMetricRICoverage(service, g, configs.Envs.TagKey, configs.Envs.TagVal); err != nil {
				return errors.Wrap(err, "on postMetricRICoverage.")
			}
		}
	}
	return nil
}

// postMetricRIUtil : post metric of RI utilization to Datadog
func postMetricRIUtil(service string, utilPercentage float64, tagKey, tagVal string) error {
	metric := "aws.ri.utilization"
	typeDatadog := "guage"

	tags := []string{
		utility.CombineStrings([]string{tagKey, ":", tagVal}),
		tagVal,
		utility.CombineStrings([]string{"service:", service}),
	}

	series := []datadog.Metric{
		{
			Metric: &metric,
			Points: []datadog.DataPoint{
				{&unixTime, &utilPercentage},
			},
			Type: &typeDatadog,
			Host: &tagVal,
			Tags: tags,
		},
	}
	return datadogClient.PostMetrics(series)
}

// postMetricRICoverage : post metric of RI utilization to Datadog
func postMetricRICoverage(service string, g *costexplorer.ReservationCoverageGroup, tagKey, tagVal string) error {
	metric := "aws.ri.coverage"
	typeDatadog := "guage"

	// string to float64
	pct, _ := strconv.ParseFloat(*g.Coverage.CoverageHours.CoverageHoursPercentage, 64)

	tags := []string{
		utility.CombineStrings([]string{"instance_type:", *g.Attributes["instanceType"]}),
		utility.CombineStrings([]string{"region:", *g.Attributes["region"]}),
		utility.CombineStrings([]string{tagKey, ":", tagVal}),
		tagVal,
		utility.CombineStrings([]string{"service:", service}),
	}

	series := []datadog.Metric{
		{
			Metric: &metric,
			Points: []datadog.DataPoint{
				{&unixTime, &pct},
			},
			Type: &typeDatadog,
			Host: &tagVal,
			Tags: tags,
		},
	}
	return datadogClient.PostMetrics(series)
}

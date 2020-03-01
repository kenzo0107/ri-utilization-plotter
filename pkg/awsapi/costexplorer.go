package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/costexplorer/costexploreriface"
)

// CostexplorerIface : costexplorer interface
type CostexplorerIface interface {
	FetchRIUtilizationPercentage(service, startDay, endDay string) (string, error)
	FetchRICoveragePercentage(service, startDay, endDay string) ([]*costexplorer.ReservationCoverageGroup, error)
}

// CostexplorerInstance : costexplorer instance
type CostexplorerInstance struct {
	client costexploreriface.CostExplorerAPI
}

// NewCostexplorer ... generate new costexplorer client
func NewCostexplorer(client costexploreriface.CostExplorerAPI) CostexplorerIface {
	return &CostexplorerInstance{
		client: client,
	}
}

// FetchRIUtilizationPercentage ... fetch RI Utilization Percentage
func (c *CostexplorerInstance) FetchRIUtilizationPercentage(service, startDay, endDay string) (riUtilPct string, err error) {
	input := &costexplorer.GetReservationUtilizationInput{
		Granularity: aws.String("DAILY"),
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(startDay),
			End:   aws.String(endDay),
		},
		Filter: &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key: aws.String("SERVICE"),
				Values: []*string{
					aws.String(service),
				},
			},
		},
	}
	r, err := c.client.GetReservationUtilization(input)
	if err != nil {
		return riUtilPct, err
	}

	// You do not use this service
	if len(r.UtilizationsByTime) == 0 {
		return
	}

	return *r.UtilizationsByTime[0].Total.UtilizationPercentage, nil
}

// FetchRICoveragePercentage ... fetch RI Coverage Percentage
func (c *CostexplorerInstance) FetchRICoveragePercentage(service, startDay, endDay string) ([]*costexplorer.ReservationCoverageGroup, error) {
	input := &costexplorer.GetReservationCoverageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(startDay),
			End:   aws.String(endDay),
		},
		Filter: &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key: aws.String("SERVICE"),
				Values: []*string{
					aws.String(service),
				},
			},
		},
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("REGION"),
			},
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("INSTANCE_TYPE"),
			},
		},
	}

	r, err := c.client.GetReservationCoverage(input)

	if err != nil {
		return []*costexplorer.ReservationCoverageGroup{}, err
	}

	return r.CoveragesByTime[0].Groups, nil
}

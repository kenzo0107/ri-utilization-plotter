package awsapi

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/costexplorer/costexploreriface"
)

// CostexplorerIface : costexplorer interface
type CostexplorerIface interface {
	FetchRIUtilizationPercentage(service, startDay, endDay string) (float64, error)
	FetchRICoveragePercentage(service, startDay, endDay string) ([]CostCoverage, error)
}

// CostexplorerInstance : costexplorer instance
type CostexplorerInstance struct {
	client costexploreriface.CostExplorerAPI
}

// CostCoverage ... include cost coverage datas.
type CostCoverage struct {
	InstanceType            string
	Region                  string
	CoverageHoursPercentage float64
}

// NewCostexplorer ... generate new costexplorer client
func NewCostexplorer(client costexploreriface.CostExplorerAPI) CostexplorerIface {
	return &CostexplorerInstance{
		client: client,
	}
}

// FetchRIUtilizationPercentage ... fetch RI Utilization Percentage
func (c *CostexplorerInstance) FetchRIUtilizationPercentage(service, startDay, endDay string) (float64, error) {
	var riUtilizationPercentage float64
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
		return riUtilizationPercentage, err
	}

	riUtilizationPercentage, _ = strconv.ParseFloat(*r.Total.UtilizationPercentage, 64)

	return riUtilizationPercentage, nil
}

// FetchRICoveragePercentage ... fetch RI Coverage Percentage
func (c *CostexplorerInstance) FetchRICoveragePercentage(service, startDay, endDay string) ([]CostCoverage, error) {
	cc := []CostCoverage{}

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
		return cc, err
	}

	for _, c := range r.CoveragesByTime {
		for _, g := range c.Groups {
			chp, _ := strconv.ParseFloat(*g.Coverage.CoverageHours.CoverageHoursPercentage, 64)
			cc = append(cc, CostCoverage{
				InstanceType:            *g.Attributes["instanceType"],
				Region:                  *g.Attributes["region"],
				CoverageHoursPercentage: chp,
			})
		}
	}

	return cc, nil
}

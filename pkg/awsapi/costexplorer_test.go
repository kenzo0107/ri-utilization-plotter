package awsapi

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/costexplorer/costexploreriface"
	"github.com/google/go-cmp/cmp"
)

type mockCostExplorerClient struct {
	costexploreriface.CostExplorerAPI

	reservationUtilizationOutput *costexplorer.GetReservationUtilizationOutput
	reservationCoverageOutput    *costexplorer.GetReservationCoverageOutput
	Error                        error
}

func (m *mockCostExplorerClient) GetReservationUtilization(*costexplorer.GetReservationUtilizationInput) (*costexplorer.GetReservationUtilizationOutput, error) {
	return m.reservationUtilizationOutput, m.Error
}

func (m *mockCostExplorerClient) GetReservationCoverage(*costexplorer.GetReservationCoverageInput) (*costexplorer.GetReservationCoverageOutput, error) {
	return m.reservationCoverageOutput, m.Error
}

func TestFetchRIUtilizationPercentage(t *testing.T) {
	m := NewCostexplorer(&mockCostExplorerClient{
		reservationUtilizationOutput: &costexplorer.GetReservationUtilizationOutput{
			Total: &costexplorer.ReservationAggregates{
				UtilizationPercentage:     aws.String("100"),
				PurchasedHours:            aws.String("48.0"),
				TotalActualHours:          aws.String("48.0"),
				UnusedHours:               aws.String("0.0"),
				OnDemandCostOfRIHoursUsed: aws.String("0.3264"),
				NetRISavings:              aws.String("0.11669041096774194"),
				TotalPotentialRISavings:   aws.String("0.11669041096774194"),
				AmortizedUpfrontFee:       aws.String("0.10410958903225806"),
				AmortizedRecurringFee:     aws.String("0.1056"),
				TotalAmortizedFee:         aws.String("0.20970958903225806"),
			},
			UtilizationsByTime: []*costexplorer.UtilizationByTime{
				&costexplorer.UtilizationByTime{
					TimePeriod: &costexplorer.DateInterval{
						Start: aws.String("2019-12-20"),
						End:   aws.String("2019-12-21"),
					},
					Total: &costexplorer.ReservationAggregates{
						UtilizationPercentage:     aws.String("100"),
						PurchasedHours:            aws.String("48"),
						TotalActualHours:          aws.String("48"),
						UnusedHours:               aws.String("0"),
						OnDemandCostOfRIHoursUsed: aws.String("0.326"),
						NetRISavings:              aws.String("0.116690411"),
						TotalPotentialRISavings:   aws.String("0.11669"),
						AmortizedUpfrontFee:       aws.String("0.104109589"),
						AmortizedRecurringFee:     aws.String("0.1056"),
						TotalAmortizedFee:         aws.String("0.209709589"),
					},
				},
			},
		},
		Error: nil,
	})

	service := "Amazon Elastic Compute Cloud - Compute"

	now := time.Now()
	startDay := now.AddDate(0, 0, -2).Format("2006-01-02")
	endDay := now.Format("2006-01-02")

	utilPercentage, err := m.FetchRIUtilizationPercentage(service, startDay, endDay)
	if err != nil {
		t.Error(err)
	}

	expected := float64(100)
	if diff := cmp.Diff(expected, utilPercentage); diff != "" {
		t.Errorf("wront result : %s", diff)
	}
}

func TestFetchRIUtilizationPercentageFailed(t *testing.T) {
	m := NewCostexplorer(&mockCostExplorerClient{
		reservationUtilizationOutput: &costexplorer.GetReservationUtilizationOutput{},
		Error:                        errors.New("error occured"),
	})

	service := "Amazon Elastic Compute Cloud - Compute"

	now := time.Now()
	startDay := now.AddDate(0, 0, -2).Format("2006-01-02")
	endDay := now.Format("2006-01-02")

	_, err := m.FetchRIUtilizationPercentage(service, startDay, endDay)
	if err == nil {
		t.Error("wrong result : err is null")
	}
}

func TestFetchRICoveragePercentage(t *testing.T) {
	m := NewCostexplorer(&mockCostExplorerClient{
		reservationCoverageOutput: &costexplorer.GetReservationCoverageOutput{
			CoveragesByTime: []*costexplorer.CoverageByTime{
				&costexplorer.CoverageByTime{
					TimePeriod: &costexplorer.DateInterval{
						Start: aws.String("2019-12-20"),
						End:   aws.String("2019-12-21"),
					},
					Groups: []*costexplorer.ReservationCoverageGroup{
						&costexplorer.ReservationCoverageGroup{
							Attributes: map[string]*string{
								"instanceType": aws.String("t3.nano"),
								"region":       aws.String("ap-northeast-1"),
							},
							Coverage: &costexplorer.Coverage{
								CoverageHours: &costexplorer.CoverageHours{
									OnDemandHours:           aws.String("24"),
									ReservedHours:           aws.String("0"),
									TotalRunningHours:       aws.String("24"),
									CoverageHoursPercentage: aws.String("0"),
								},
							},
						},
						&costexplorer.ReservationCoverageGroup{
							Attributes: map[string]*string{
								"instanceType": aws.String("t2.micro"),
								"region":       aws.String("ap-northeast-3"),
							},
							Coverage: &costexplorer.Coverage{
								CoverageHours: &costexplorer.CoverageHours{
									OnDemandHours:           aws.String("24"),
									ReservedHours:           aws.String("24"),
									TotalRunningHours:       aws.String("48"),
									CoverageHoursPercentage: aws.String("50"),
								},
							},
						},
					},
				},
			},
			Total: &costexplorer.Coverage{
				CoverageHours: &costexplorer.CoverageHours{
					OnDemandHours:           aws.String("48"),
					ReservedHours:           aws.String("96"),
					TotalRunningHours:       aws.String("144"),
					CoverageHoursPercentage: aws.String("66.6666666667"),
				},
			},
		},
		Error: nil,
	})

	service := "Amazon Elastic Compute Cloud - Compute"

	now := time.Now()
	startDay := now.AddDate(0, 0, -2).Format("2006-01-02")
	endDay := now.Format("2006-01-02")
	costCoverages, err := m.FetchRICoveragePercentage(service, startDay, endDay)
	if err != nil {
		t.Error(err)
	}

	expected := []CostCoverage{
		CostCoverage{
			InstanceType:            "t3.nano",
			Region:                  "ap-northeast-1",
			CoverageHoursPercentage: 0,
		},
		CostCoverage{
			InstanceType:            "t2.micro",
			Region:                  "ap-northeast-3",
			CoverageHoursPercentage: 50,
		},
	}

	if diff := cmp.Diff(expected, costCoverages); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
}

func TestFetchRICoveragePercentageFailed(t *testing.T) {
	m := NewCostexplorer(&mockCostExplorerClient{
		reservationCoverageOutput: &costexplorer.GetReservationCoverageOutput{},
		Error:                     errors.New("error occured"),
	})

	service := "Amazon Elastic Compute Cloud - Compute"

	now := time.Now()
	startDay := now.AddDate(0, 0, -2).Format("2006-01-02")
	endDay := now.Format("2006-01-02")
	_, err := m.FetchRICoveragePercentage(service, startDay, endDay)
	if err == nil {
		t.Error("wrong result : err is nil")
	}
}

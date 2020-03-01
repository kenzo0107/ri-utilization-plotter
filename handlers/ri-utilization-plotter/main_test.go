package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	services = []string{
		"Amazon Elastic Compute Cloud - Compute",
		"Amazon Redshift",
	}
	type args struct {
		ctx                  context.Context
		ssmDatadogAPIKeyName string
		ssmDatadogAPPKeyName string
		configAWSAccount     string
		startDay, endDay     string
	}

	now := time.Now()

	ctx := context.Background()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successfully",
			args: args{
				ctx:                  ctx,
				ssmDatadogAPIKeyName: "datadog_api_key",
				ssmDatadogAPPKeyName: "datadog_app_key",
				configAWSAccount:     filepath.Join("..", "..", "testdata", "awsaccount.yml"),
				startDay:             now.AddDate(0, 0, -3).Format("2006-01-02"),
				endDay:               now.Format("2006-01-02"),
			},
			wantErr: false,
		},
		{
			name: "not set Datadog API key name as environment values",
			args: args{
				ctx:                  ctx,
				ssmDatadogAPIKeyName: "",
				ssmDatadogAPPKeyName: "",
				configAWSAccount:     "",
				startDay:             now.AddDate(0, 0, -2).Format("2006-01-02"),
				endDay:               now.Format("2006-01-02"),
			},
			wantErr: true,
		},
		{
			name: "not found aws account configure file",
			args: args{
				ctx:                  ctx,
				ssmDatadogAPIKeyName: "datadog_api_key",
				ssmDatadogAPPKeyName: "datadog_app_key",
				configAWSAccount:     filepath.Join("..", "..", "testdata", "failedAWSAccount.yml"),
				startDay:             now.AddDate(0, 0, -2).Format("2006-01-02"),
				endDay:               now.Format("2006-01-02"),
			},
			wantErr: true,
		},
		{
			name: "set incorrect datadog api key",
			args: args{
				ctx:                  ctx,
				ssmDatadogAPIKeyName: "incorrect_datadog_api_key",
				ssmDatadogAPPKeyName: "incorrect_datadog_app_key",
				configAWSAccount:     filepath.Join("..", "..", "testdata", "awsaccount.yml"),
				startDay:             now.AddDate(0, 0, -2).Format("2006-01-02"),
				endDay:               now.Format("2006-01-02"),
			},
			wantErr: true,
		},
		{
			name: "not default account in awsaccount.yml",
			args: args{
				ctx:                  ctx,
				ssmDatadogAPIKeyName: "datadog_api_key",
				ssmDatadogAPPKeyName: "datadog_app_key",
				configAWSAccount:     filepath.Join("..", "..", "testdata", "awsaccount2.yml"),
				startDay:             now.AddDate(0, 0, -2).Format("2006-01-02"),
				endDay:               now.Format("2006-01-02"),
			},
			wantErr: true,
		},
		{
			name: "start date cannot be after 2 days ago",
			args: args{
				ctx:                  ctx,
				ssmDatadogAPIKeyName: "datadog_api_key",
				ssmDatadogAPPKeyName: "datadog_app_key",
				configAWSAccount:     filepath.Join("..", "..", "testdata", "awsaccount.yml"),
				startDay:             now.Format("2006-01-02"),
				endDay:               now.Format("2006-01-02"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ssmDatadogAPIKeyName = tt.args.ssmDatadogAPIKeyName
			ssmDatadogAPPKeyName = tt.args.ssmDatadogAPPKeyName
			configAWSAccount = tt.args.configAWSAccount
			startDay = tt.args.startDay
			endDay = tt.args.endDay
			if err := handler(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("handler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

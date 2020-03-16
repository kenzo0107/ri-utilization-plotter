package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zorkian/go-datadog-api"
)

func TestHandler(t *testing.T) {
	// datadog のエンドポイントへメトリクスをプロットする際の必ず200ステータスを返す（成功する）テストサーバ
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer ts.Close()

	// PostMetrics が必ず成功する datadog Client
	ddClient := &datadog.Client{
		HttpClient: http.DefaultClient,
	}
	ddClient.SetBaseUrl(ts.URL)

	// datadog のエンドポイントへメトリクスをプロットする際の必ず 403 ステータスを返す（失敗する）テストサーバ
	tsFailed := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
	}))
	defer tsFailed.Close()

	// PostMetrics が必ず失敗する datadog Client
	ddClientFailed := &datadog.Client{
		HttpClient: http.DefaultClient,
	}
	ddClientFailed.SetBaseUrl(tsFailed.URL)

	services = []string{
		"Amazon Elastic Compute Cloud - Compute",
		"Amazon Redshift",
	}
	type args struct {
		ctx              context.Context
		startDay, endDay string
		datadogClient    *datadog.Client
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
				ctx:           ctx,
				startDay:      now.AddDate(0, 0, -3).Format("2006-01-02"),
				endDay:        now.Format("2006-01-02"),
				datadogClient: ddClient,
			},
			wantErr: false,
		},
		{
			name: "start date cannot be after 2 days ago",
			args: args{
				ctx:           ctx,
				startDay:      now.Format("2006-01-02"),
				endDay:        now.Format("2006-01-02"),
				datadogClient: ddClient,
			},
			wantErr: true,
		},
		{
			name: "failed to post metric to datadog",
			args: args{
				ctx:           ctx,
				startDay:      now.AddDate(0, 0, -3).Format("2006-01-02"),
				endDay:        now.Format("2006-01-02"),
				datadogClient: ddClientFailed,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startDay = tt.args.startDay
			endDay = tt.args.endDay
			datadogClient = tt.args.datadogClient
			if err := handler(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("handler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

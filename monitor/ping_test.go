package monitor_test

import (
	"context"
	"encore.app/monitor"
	"testing"
)

func TestPing(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name string
		URL  string
		Up   bool
	}{
		{name: "encore", URL: "encore.com", Up: true},
		{name: "encore dev", URL: "encore.dev", Up: true},
		{name: "httpbin without prefix", URL: "httpbin.org/status/200", Up: true},
		{name: "httpbin with prefix", URL: "https://httpbin.org/status/200", Up: true},
		{name: "httpbin status 400", URL: "https://httpbin.org/status/400", Up: false},
		{name: "httpbin status 500", URL: "https://httpbin.org/status/500", Up: false},
		{name: "invalid urls", URL: "invalid://schema", Up: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := monitor.Ping(ctx, tt.URL)
			if err != nil {
				t.Errorf("url %s: unexpected error: %v", tt.URL, err)
			} else if resp.Up != tt.Up {
				t.Errorf("url %s: expected %v, got %v", tt.URL, tt.Up, resp.Up)
			}
		})
	}
}

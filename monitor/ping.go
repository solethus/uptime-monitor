// Package monitor Service monitor checks if a website is up or down
package monitor

import (
	"context"
	log "log/slog"
	"net/http"
	"strings"
)

// PingResponse is the response from the ping endpoint
type PingResponse struct {
	Up bool `json:"up"`
}

// Ping pings a specific site and determines whether its up or down right now
//
//encore:api public path=/ping/:url
func Ping(ctx context.Context, url string) (*PingResponse, error) {
	// If the url does not start with "http:" or "https:", default to "https:".
	if !strings.HasPrefix(url, "http:") && !strings.HasPrefix(url, "https:") {
		url = "https://" + url
	}

	// Make an HTTP request to check if it's up.
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.ErrorContext(ctx, err.Error())
		return &PingResponse{Up: false}, nil
	}
	defer resp.Body.Close()

	// 2xx and 3xx codes are considered up
	up := resp.StatusCode < 400

	return &PingResponse{Up: up}, nil
}

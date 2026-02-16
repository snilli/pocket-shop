package ezclient

import (
	"context"
	"math"
	"net/http"
	"time"
)

func isRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp != nil && resp.StatusCode >= 500 {
		return true
	}
	return false
}

func withRetry(ctx context.Context, maxAttempts, backoffSec int, fn func() (*http.Response, error)) (*http.Response, error) {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	if backoffSec < 1 {
		backoffSec = 1
	}
	var lastResp *http.Response
	var lastErr error
	for attempt := range maxAttempts {
		resp, err := fn()
		lastResp, lastErr = resp, err
		if !isRetryable(resp, err) {
			return resp, err
		}
		if attempt < maxAttempts-1 {
			secs := float64(backoffSec) * math.Pow(2, float64(attempt))
			backoff := time.Duration(secs) * time.Second
			select {
			case <-ctx.Done():
				return resp, ctx.Err()
			case <-time.After(backoff):
			}
		}
	}
	return lastResp, lastErr
}

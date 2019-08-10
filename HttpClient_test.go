package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SpyRetryPolicy struct {
	maxRetries        int
	retryCount        int
	backoffCalls      int
	retryDelayCalls   int
	resetRetriesCalls int
	finalTryCalls     int
}

func (s *SpyRetryPolicy) backoff() {
	s.backoffCalls++
}

func (s *SpyRetryPolicy) getRetryDelay() time.Duration {
	s.retryDelayCalls++
	return 0
}

func (s *SpyRetryPolicy) resetRetries() {
	s.resetRetriesCalls++
	s.retryCount = 0
}

func (s *SpyRetryPolicy) isFinalTry(try int) bool {
	s.finalTryCalls++
	return s.maxRetries <= try
}

func (s *SpyRetryPolicy) getMaxRetries() int {
	return s.maxRetries
}

func TestURLFetcher(t *testing.T) {
	t.Run("fetch URL no maxRetries no body, return successful", func(t *testing.T) {
		stubRetryPolicy := &SpyRetryPolicy{
			maxRetries: 3,
		}
		statusCode := http.StatusOK
		httpClient := NewRetryHttpClientWithPolicy(0, stubRetryPolicy)

		ts := buildHttpServer(t, statusCode)
		defer ts.Close()

		response, err := httpClient.getResponse(ts.URL)

		assertNoError(t, err)
		assertStatusCode(t, response.StatusCode, statusCode)
		assertRetryValue(t, stubRetryPolicy.backoffCalls, 0)
		assertRetryValue(t, stubRetryPolicy.resetRetriesCalls, 1)
	})

	t.Run("fail, no retries", func(t *testing.T) {
		maxRetries := 0
		stubRetryPolicy := &SpyRetryPolicy{
			maxRetries: maxRetries,
		}
		statusCode := http.StatusNotFound
		httpClient := NewRetryHttpClientWithPolicy(0, stubRetryPolicy)

		ts := buildHttpServer(t, statusCode)
		defer ts.Close()

		_, err := httpClient.getResponse(ts.URL)

		assertErrorMessage(t, err, fmt.Sprintf(errorMessage, ts.URL, statusCode))
		assertRetryValue(t, stubRetryPolicy.backoffCalls, maxRetries)
		assertRetryValue(t, stubRetryPolicy.resetRetriesCalls, 1)
	})

	t.Run("fetch URL 3 maxRetries, fail, getRetryDelay 3 times, return 500 error", func(t *testing.T) {
		maxRetries := 3
		stubRetryPolicy := &SpyRetryPolicy{
			maxRetries: maxRetries,
		}
		statusCode := http.StatusInternalServerError
		httpClient := NewRetryHttpClientWithPolicy(0, stubRetryPolicy)

		ts := buildHttpServer(t, statusCode)
		defer ts.Close()

		_, err := httpClient.getResponse(ts.URL)

		assertErrorMessage(t, err, fmt.Sprintf(errorMessage, ts.URL, statusCode))
		assertRetryValue(t, stubRetryPolicy.backoffCalls, maxRetries)
	})

	t.Run("failing resets retry count", func(t *testing.T) {
		stubRetryPolicy := &SpyRetryPolicy{
			maxRetries: 3,
		}
		statusCode := http.StatusInternalServerError
		httpClient := NewRetryHttpClientWithPolicy(0, stubRetryPolicy)

		ts := buildHttpServer(t, statusCode)
		defer ts.Close()

		_, err := httpClient.getResponse(ts.URL)

		assertErrorMessage(t, err, fmt.Sprintf(errorMessage, ts.URL, statusCode))
		assertRetryValue(t, stubRetryPolicy.retryCount, 0)
	})
}

func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %v but wanted %v", got, want)
	}
}

func assertRetryValue(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, but wanted %v", got, want)
	}
}

func buildHttpServer(t *testing.T, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))
}

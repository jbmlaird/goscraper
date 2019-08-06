package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLFetcher(t *testing.T) {
	t.Run("fetch URL, return successful", func(t *testing.T) {
		statusCode := http.StatusOK
		httpClient := NewHttpClient(0, 0)

		ts := buildHttpServer(t, statusCode)

		response, err := httpClient.fetchUrl(ts.URL)

		assertNoError(t, err)
		assertStatusCode(t, response.StatusCode, statusCode)
		assertRetryValue(t, httpClient.retries, 0)
	})

	t.Run("fetch URL, fail, retry, return timeout", func(t *testing.T) {
		statusCode := http.StatusGatewayTimeout
		httpClient := NewHttpClient(3, 0)

		ts := buildHttpServer(t, statusCode)

		response, err := httpClient.fetchUrl(ts.URL)

		assertNoError(t, err)
		assertStatusCode(t, response.StatusCode, statusCode)
		assertRetryValue(t, httpClient.retries, 3)
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

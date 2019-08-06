package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLFetcher(t *testing.T) {
	t.Run("fetch URL no retries no body, return successful", func(t *testing.T) {
		statusCode := http.StatusOK
		httpClient := NewHttpClient(0, 0, 0, 0)

		ts := buildHttpServer(t, statusCode)
		defer ts.Close()

		response, err := httpClient.getUrl(ts.URL)

		assertNoError(t, err)
		assertStatusCode(t, response.StatusCode, statusCode)
		assertRetryValue(t, httpClient.retryPolicy.retries, 0)
	})

	t.Run("fetch URL 3 retries, fail, getRetryDelay 3 times, return 500 error", func(t *testing.T) {
		statusCode := http.StatusInternalServerError
		httpClient := NewHttpClient(3, 0, 0, 0)

		ts := buildHttpServer(t, statusCode)
		defer ts.Close()

		_, err := httpClient.getUrl(ts.URL)

		assertErrorMessage(t, err, fmt.Sprintf(errorMessage, ts.URL))
		assertRetryValue(t, httpClient.retryPolicy.retries, 3)
	})
	// No unit test for timeout as that is not my code
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

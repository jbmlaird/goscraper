package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLFetcher(t *testing.T) {
	t.Run("fetch URL, return successful", func(t *testing.T) {
		httpClient := HttpClient{
			retries:    0,
			retryDelay: 0,
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		response, err := httpClient.fetchUrl(ts.URL)

		assertNoError(t, err)

		assertStatus(t, response.StatusCode, http.StatusOK)
	})
}

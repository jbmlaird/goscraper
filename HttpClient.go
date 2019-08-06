package main

import (
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

// Golang's standard library http.Client does not contain a retry policy
type RetryHttpClient struct {
	retryPolicy RetryPolicy
	http.Client
}

func NewHttpClient(retries, retryCount, retryDelay, timeoutSeconds int) *RetryHttpClient {
	return &RetryHttpClient{
		RetryPolicy{
			retries:           retries,
			retryCount:        retryCount,
			retryDelaySeconds: retryDelay,
		},
		http.Client{
			Timeout: time.Second * time.Duration(timeoutSeconds),
		},
	}
}

const errorMessage = "Error fetching URL %v"

func (r *RetryHttpClient) getResponse(url string) (*http.Response, error) {
	var (
		response *http.Response
		err      error
	)
	for r.retryPolicy.retryCount = 0; r.retryPolicy.retryCount <= r.retryPolicy.retries; r.retryPolicy.retryCount++ {
		response, err = r.Get(url)
		if err != nil {
			return nil, err
		}
		if response != nil {
			if response.StatusCode == http.StatusOK {
				break
			} else {
				response.Body.Close()
				if r.retryPolicy.isFinalTry() {
					return nil, errors.Errorf(errorMessage, url)
				} else {
					log.Printf("Error fetching URL: %v", url)
					time.Sleep(r.retryPolicy.getRetryDelay())
					r.retryPolicy.backoff()
				}
			}
		}
	}
	return response, nil
}

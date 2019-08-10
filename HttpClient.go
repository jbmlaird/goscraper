package main

import (
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

// Golang's standard library http.client does not contain a retry policy
type RetryHttpClient struct {
	retryPolicy RetryPolicy
	http.Client
}

func NewRetryHttpClient(retries, retryDelay, timeoutSeconds int) *RetryHttpClient {
	return &RetryHttpClient{
		&ProdRetryPolicy{
			maxRetries:        retries,
			retryDelaySeconds: retryDelay,
		},
		http.Client{
			Timeout: time.Second * time.Duration(timeoutSeconds),
		},
	}
}

// Had to create this to provide a stub retry policy. I don't like this.
func NewRetryHttpClientWithPolicy(timeoutSeconds int, retryPolicy RetryPolicy) *RetryHttpClient {
	return &RetryHttpClient{
		retryPolicy,
		http.Client{
			Timeout: time.Second * time.Duration(timeoutSeconds),
		},
	}
}

const errorMessage = "Unable to fetch URL %v with status code %v"

func (r *RetryHttpClient) getResponse(url string) (*http.Response, error) {
	var (
		response *http.Response
		err      error
	)
	defer r.retryPolicy.resetRetries()
	for i := 0; i < r.retryPolicy.getMaxRetries()+1; i++ {
		response, err = r.Get(url)
		if err != nil {
			if r.retryPolicy.isFinalTry(i) {
				return nil, errors.Wrapf(err, "Failed final attempt for URL %v", url)
			}
			err = errors.Wrapf(err, "Failed attempt #%d for URL %v. Retrying", i + 1, url)
		}
		if response != nil {
			if response.StatusCode == http.StatusOK {
				break
			}
			if r.retryPolicy.isFinalTry(i) {
				response.Body.Close()
				return nil, errors.Errorf(errorMessage, url, response.StatusCode)
			}
		}
		log.Printf("Error fetching URL attempt #%d: %v", i + 1, url)
		time.Sleep(r.retryPolicy.getRetryDelay())
		r.retryPolicy.backoff()
	}
	return response, nil
}

package main

import (
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

type HttpClient struct {
	retries    int
	retryCount int
	retryDelay int
}

func NewHttpClient(retries, retryDelay int) *HttpClient {
	return &HttpClient{
		retries:    retries,
		retryCount: 0,
		retryDelay: retryDelay,
	}
}

func (h *HttpClient) fetchUrl(url string) (*http.Response, error) {
	// timeouts
	var (
		response *http.Response
		err      error
	)
	for h.retryCount = 0; h.retryCount < h.retries+1; h.retryCount++ {
		response, err = http.Get(url)
		response.Body.Close()
		if err != nil {
			if h.retryCount == h.retries-1 {
				err = errors.Wrapf(err, "Error fetching URL %v, err: %v", url, err)
			} else {
				log.Printf("Error fetching URL %v, err: %v", url, err)
				time.Sleep(time.Second * time.Duration(h.retryDelay))
				h.retryDelay = h.retryDelay * 2
			}
		}
		if response != nil && response.StatusCode == http.StatusOK {
			break
		}
	}
	return response, err
}

package main

import (
	"errors"
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
		response   *http.Response
		err        error
		logMessage string
	)
	for h.retryCount = 0; h.retryCount < h.retries; h.retries++ {
		response, err = http.Get(url)
		if err != nil {
			if h.retryCount == h.retries-1 {
				logMessage = "Giving up 😱"
			} else {
				logMessage = "Retrying"
			}
			log.Printf("Error fetching URL %v, err: %v. %v", url, err, logMessage)
			time.Sleep(time.Second * time.Duration(h.retryDelay))
			h.retryDelay = h.retryDelay * 2
		}
	}
	if response != nil {
		defer response.Body.Close()
		// do some swag shit with the links
		return response, nil
	}
	return nil, errors.New("unable to fetch request")
}

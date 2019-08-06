package main

import "time"

type RetryPolicy struct {
	retries           int
	retryCount        int
	retryDelaySeconds int
}

func (r *RetryPolicy) backoff() {
	r.retries = r.retries * 2
}

func (r *RetryPolicy) getRetryDelay() time.Duration {
	return time.Second * time.Duration(r.retries)
}

func (r *RetryPolicy) isFinalTry() bool {
	return r.retries >= r.retryCount
}

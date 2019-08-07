package main

import "time"

type RetryPolicy interface {
	backoff()
	resetRetries()
	isFinalTry(try int) bool
	getRetryDelay() time.Duration
	getMaxRetries() int
}

type ProdRetryPolicy struct {
	maxRetries               int
	initialRetryDelaySeconds int
	retryDelaySeconds        int
}

func (p *ProdRetryPolicy) backoff() {
	p.retryDelaySeconds = p.retryDelaySeconds * 2
}

func (p *ProdRetryPolicy) resetRetries() {
	p.retryDelaySeconds = p.initialRetryDelaySeconds
}

func (p *ProdRetryPolicy) isFinalTry(try int) bool {
	return p.maxRetries <= try
}

func (p *ProdRetryPolicy) getRetryDelay() time.Duration {
	return time.Second * time.Duration(p.retryDelaySeconds)
}

func (p *ProdRetryPolicy) getMaxRetries() int {
	return p.maxRetries
}

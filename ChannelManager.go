package main

type ChannelManager struct {
	chanError   chan error
	chanCrawled chan string
}

func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		chanError:   make(chan error),
		chanCrawled: make(chan string),
	}
}

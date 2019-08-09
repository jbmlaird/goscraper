package main

import (
	"fmt"
	"time"
)

type ChannelManager struct {
	chanError   chan error
	chanCrawled chan string
	chanTimeout chan bool
}

func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		chanError:   make(chan error),
		chanCrawled: make(chan string),
		chanTimeout: make(chan bool, 1),
	}
}

func (c *ChannelManager) StartListening() {
	go func() {
		for {
			select {
			case err := <-c.chanError:
				err.Error()
				break
			case crawledUrl := <-c.chanCrawled:
				fmt.Println(crawledUrl)
				break
			case <-time.After(time.Second * 10):
				c.chanTimeout <- true
				break
			}
		}
	}()
}

func (c *ChannelManager) CloseChannels() {
	close(c.chanError)
	close(c.chanCrawled)
	close(c.chanTimeout)
}

func (c *ChannelManager) listenToErrors() {
	for {
		<-c.chanError
	}
}

func (c *ChannelManager) listenToCrawledLinks() {
	for {
		<-c.chanCrawled
	}
}

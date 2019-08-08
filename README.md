# Go Crawler

To build:

1. `go get github.com/PuerkitoBio/goquery`
2. `go get github.com/pkg/errors`
3. go build main

---

Improvements:

* Cancellation of requests
* What if a goroutine fails against a certain URL? Remove it from the sitemap and retry?
* Error handling for goroutines is not done. I have a separate branch for attempting to use channels but the WaitGroup
hung forever

Notes:

* I have two data structures. A map for tracking URLs visited, and a slice for having a sorted sitemap

# Go Scraper

To build:

1. `go get github.com/PuerkitoBio/goquery`
2. `go get github.com/pkg/errors`
3. go build main

---

Improvements:

* `http.Client`'s timeout behaviour is not tested since this belongs to the standard library.
* Support absolute paths rather than just relative paths
* `NewRetryHttpClientWithPolicy` :'(
* Cancellation
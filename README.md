# Go Scraper

To build:

1. `go get github.com/PuerkitoBio/goquery`
2. `go get github.com/pkg/errors`
3. go build main

---

Improvements:

* If the user input an https protocol and the page has absolute paths with http, this will currently fail
* `NewRetryHttpClientWithPolicy` :'(
* Cancellation of requests

Notes:

* I over-engineered the initial URL checking according to the requirements. Should have just accepted a website with the
HTTPS protocol like in the initial email
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
* `/` is a different domain error needs to be fixed. Instead be, `this is the homepage`
* Strip things after `#`
* Stop crawling URLs twice (or sitemapping) if they end in a slash and when they don't
* Fix Twitter links appearing

Notes:

* I over-engineered the initial URL checking according to the requirements. Should have just accepted a website with the
HTTPS protocol like in the initial email
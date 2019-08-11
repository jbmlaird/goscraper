# Go Crawler

To execute:

1. `go build`
2. `./goscraper`

We in business.

Sitemap outputs to `sitemap.txt`, goroutine errors output to `goroutineErrors.txt`.

Logging of errors happen as the program runs but due to the large console output, it's easier to read them in a file.

---

Notes: 

HTTP will always be reformatted to HTTPS. No other protocols are supported.

Improvements:

* Ignore crawling URLs if it links to a file e.g. `.pdf` or `.csv`
* Cancelled request handling
* Allow the passing in of a base URL on execution via command line
* Hierarchical output of sitemap as opposed to sorted list

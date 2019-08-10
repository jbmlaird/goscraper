# Go Crawler

To execute:

1. `go build`
2. `./goscraper`

We in business.

Sitemap outputs to `sitemap.txt`, goroutine errors output to `goroutineErrors.txt`.

Logging of errors happen as the program runs but due to the large console output, it's easier to read them in a file.

---

Notes: 

SMTP and FTP links are not supported. HTTP will always be reformatted to HTTPS.

Improvements:

* Ignore crawling URLs if it links to a file e.g. `.pdf` or `.csv`
* Ignore crawling telephone number URLs
* Cancelled request handling
* Allow the passing in of a base URL on execution
* Hierarchical output of sitemap as opposed to sorted list
* Is a single mutex in SitemapBuilder suitable? Or should there be two for the different adding to maps? Would that
increase performance?

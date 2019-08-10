# Go Crawler

To execute:

1. `go build`
2. `./goscraper`

We in business.

Sitemap outputs to `sitemap.txt`, goroutine errors and logs output to `goroutineErrors.txt`.

---

Notes: 

SMTP and FTP links are not supported. HTTP will always be reformatted to HTTPS.
Sitemap outputs to a text file called `sitemap.txt`. This is just a sorted list rather than any hierarchical relationship.

Improvements:

* Ignore crawling URLs if it links to a file e.g. `.pdf` or `.csv`
* Ignore crawling telephone number URLs
* Cancelled request handling and tests
* Allow the passing in of a base URL on execution
* Hierarchical output of sitemap
* Is a single mutex in SitemapBuilder suitable? Or should there be two for the different adding to maps? Would that
increase performance?

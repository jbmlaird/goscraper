# Go Crawler

To execute:

1. `go build`
2. `./goscraper`

We in business.

---

Notes: 

SMTP and FTP are not supported.
Sitemap outputs to a text file called `sitemap.txt`. This is just a sorted list rather than any hierarchical relationship.

Improvements:

* Ignore crawling URLs if it links to a file e.g. .pdf or .csv
* Cancellation of request handling
* Allow the passing in of a base URL on execution
* Hierarchical output of sitemap
* Is a single mutex in SitemapBuilder suitable? Or should there be two for the different adding to maps?
* Ignore crawling telephone number URLs 
* Passing down the links that were searched before the current URL would help to track down pages that have incorrect
URLs. I found `monzo://overdraft` in my `goroutineErrors.txt`, it would be great to know which pages had that incorrect
URL
# Go Crawler

To build:

1. `go get github.com/PuerkitoBio/goquery`
2. `go get github.com/pkg/errors`
3. go build main

---

Improvements:

* Cancellation of requests

Notes:

* Hangs forever on this channels version. Due to receiving channels blocking forever until the message is processed
    * going to have to think about how to handle this
    * something needs to process the unbuffered channel in order to unblock. Could have something which just
    continuously collects the error messages 
* Didn't know Go a month ago :'(
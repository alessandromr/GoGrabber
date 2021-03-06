package httprequester

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"time"
)

//HTTPRequester manage http request
type HTTPRequester struct {
	URL       string
	UserAgent string
	Timeout   time.Duration
}

//MakeCheck Make HTTP Request to a given url and return a goquery document
func (req HTTPRequester) MakeCheck() (goquery.Document, error) {
	client := &http.Client{
		Timeout: req.Timeout * time.Second,
	}
	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", req.URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", req.UserAgent)

	// Make request
	response, err := client.Do(request)
	if err != nil {
		// log.Println("Error making HTTP request.", err)
		return *new(goquery.Document), err
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		// log.Println("Error loading HTTP response body. ", err)
		return *new(goquery.Document), err
	}
	return *document, err
}

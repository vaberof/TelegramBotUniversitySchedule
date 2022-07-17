package http

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// LoadHtmlPage loads html page.
func LoadHtmlPage(url string) (*goquery.Document, error) {
	rBody, err := makeRequest(url)

	if err != nil {
		emptyDocument := &goquery.Document{}
		return emptyDocument, err
	}

	document, err := goquery.NewDocumentFromReader(rBody)

	if err != nil {
		log.Printf("data cannot be parsed as html %v\n", err)
		return document, err
	}

	return document, nil
}

// makeRequest
// returns response body of type Reader end error.
func makeRequest(url string) (io.Reader, error) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	res, err := client.Get(url)

	if err != nil {
		log.Printf("request on %s failed: %v\n", url, err)
		emptyBody := bytes.NewReader([]byte{})
		return emptyBody, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Printf("read body failed: %v\n", err)
		emptyBody := bytes.NewReader([]byte{})
		return emptyBody, err
	}

	defer res.Body.Close()

	rBody := bytes.NewReader(body)

	return rBody, nil
}

package http

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
)

// LoadHtmlPage loads html page.
func LoadHtmlPage(url string) *goquery.Document {
	document, err := goquery.NewDocumentFromReader(makeRequest(url))
	if err != nil {
		fmt.Errorf("data cannot be parsed as html %v\n", err)
	}

	return document
}

// makeRequest
// returns response body of type Reader
// to use it in LoadHtmlPage.
func makeRequest(url string) io.Reader {

	res, err := http.Get(url)

	if err != nil {
		fmt.Errorf("request on %s failed: %w\n", url, err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("read body failed %w\n", err)
	}

	defer res.Body.Close()

	rBody := bytes.NewReader(body)

	return rBody
}

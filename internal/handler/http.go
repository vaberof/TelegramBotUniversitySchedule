package handler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// makeRequest
// returns response body of type Reader.
func makeRequest(url string) io.Reader {

	res, err := http.Get(url)
	if err != nil {
		fmt.Errorf("request on %s failed: %w\n",url, err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("read body failed %w\n", err)
		log.Println(err)
	}

	rBody := bytes.NewReader(body)

	return rBody
}

// LoadHtmlPage loads html page.
func LoadHtmlPage(url string) *goquery.Document{

	document, err := goquery.NewDocumentFromReader(makeRequest(url))
	if err != nil {
		fmt.Errorf("data cannot be parsed as html %v", err)
	}

	return document
}
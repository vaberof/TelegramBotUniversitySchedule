package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/PuerkitoBio/goquery"
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
		log.WithFields(log.Fields{
			"body":  rBody,
			"error": err,
			"func":  "LoadHtmlPage",
		}).Error("Data cannot be parsed as html")

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
		log.WithFields(log.Fields{
			"url":   url,
			"error": err,
			"func":  "makeRequest",
		}).Error("Request is failed")

		emptyBody := bytes.NewReader([]byte{})
		return emptyBody, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.WithFields(log.Fields{
			"url":   url,
			"error": err,
			"func":  "makeRequest",
		}).Error("Read body is failed")

		emptyBody := bytes.NewReader([]byte{})
		return emptyBody, err
	}

	defer res.Body.Close()

	rBody := bytes.NewReader(body)

	return rBody, nil
}

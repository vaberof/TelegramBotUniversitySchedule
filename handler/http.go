package handler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

// Делаем запрос
func makeRequest(url string) *http.Response {

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	//defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("status code error: %d %s", res.StatusCode, res.Status)
	}
	body := res

	return body
}

// Загружаем HTML страничку
func LoadHtmlPage(url string) *goquery.Document{

	document, err := goquery.NewDocumentFromReader(makeRequest(url).Body)

	if err != nil {
		log.Println(err)
	}
	return document
}
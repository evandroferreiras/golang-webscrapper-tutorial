package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func scrapListURL(urlToProcess []string, rchan chan Result) {
	defer close(rchan)
	var results = []chan Result{}

	for i, url := range urlToProcess {
		results = append(results, make(chan Result))
		go scrapParallel(url, results[i])
	}

	for i := range results {
		for r1 := range results[i] {
			rchan <- r1
		}
	}

}

func scrapParallel(url string, rchan chan Result) {
	defer close(rchan)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: It can't scrap '", url, "'")
	}
	// Close body when function ends
	defer resp.Body.Close()
	body := resp.Body
	htmlParsed, err := html.Parse(body)
	if err != nil {
		fmt.Println("ERROR: It can't parse html '", url, "'")
	}
	header := getFirstElementByClass(htmlParsed, "header", "")
	var r Result
	a := getFirstElementByClass(header, "a", "ds-link--styleSubtle")
	r.userName = getFirstTextNode(a).Data

	div := getFirstElementByClass(htmlParsed, "div", "section-content")
	h1 := getFirstElementByClass(div, "h1", "graf--title")
	r.title = getFirstTextNode(h1).Data

	footer := getFirstElementByClass(htmlParsed, "footer", "u-paddingTop10")
	buttonLikes := getFirstElementByClass(footer, "button", "js-multirecommendCountButton")
	r.likes = getFirstTextNode(buttonLikes).Data

	rchan <- r
}

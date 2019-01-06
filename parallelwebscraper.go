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

	var r Result
	a := getFirstTextNode(getFirstElementByClass(htmlParsed, "a", "ds-link--styleSubtle"))
	if a != nil {
		r.userName = a.Data
	} else {
		fmt.Println("Scrap error: Can't find username. url:'", url, "'")
	}

	div := getFirstElementByClass(htmlParsed, "div", "section-content")
	h1 := getFirstTextNode(getFirstElementByClass(div, "h1", "graf--title"))
	if h1 != nil {
		r.title = h1.Data
	} else {
		fmt.Println("Scrap error: Can't find title. url:'", url, "'")
	}

	footer := getFirstElementByClass(htmlParsed, "footer", "u-paddingTop10")
	buttonLikes := getFirstTextNode(getFirstElementByClass(footer, "button", "js-multirecommendCountButton"))
	if buttonLikes != nil {
		r.likes = buttonLikes.Data
	} else {
		fmt.Println("Scrap error: Can't find button of likes. url:'", url, "'")
	}

	rchan <- r
}

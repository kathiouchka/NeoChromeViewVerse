package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gocolly/colly"
)

func testScrape() {
	fmt.Println("HELLO")
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})
	c.Visit("https://www.blueboard.io")
}

func scrapeIPInfo() {
	pURL, _ := url.Parse(`http://85.25.198.20:5566`)
	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(pURL)}}
	resp, err := httpClient.Get("https://www.ipinfo.io")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
}

func main() {
	testScrape()
	// scrapeIPInfo()
}

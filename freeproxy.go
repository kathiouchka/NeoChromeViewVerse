package main

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

func testScrapeFreeProxy() {
	c := colly.NewCollector(colly.AllowURLRevisit())

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	rp, err := proxy.RoundRobinProxySwitcher("http://80.48.119.28:8080", "http://169.57.1.85:8123")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	c.OnResponse(func(r *colly.Response) {
		log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
		fmt.Println("PAGE STATUS = ", r.StatusCode)
	})

	c.OnHTML("#proxy_list", func(e *colly.HTMLElement) {
		e.ForEach("tbody tr", func(_ int, el *colly.HTMLElement) {
			i := 0
			var proxyURL string
			el.ForEachWithBreak("td", func(_ int, elem *colly.HTMLElement) bool {
				if i == 0 && !strings.Contains(elem.Text, "adsbygoogle") {
					base64 := getStringBetween(elem.Text, `("`, `")`)
					byteURL, err := b64.StdEncoding.DecodeString(base64)
					if err != nil {
						log.Fatalln(err)
					}
					proxyURL = string(byteURL)
				}
				if i == 1 {
					retURL := proxyURL + ":" + elem.Text
					fmt.Println(retURL)
				}
				if i == 2 {
					return false
				}
				i++
				return true
			})
		})
	})
	// TODO - Pourquoi ca scrap pas la bonne page??
	for i := 0; i < 5; i++ {
		fmt.Println("I = ", i)
		c.Visit("http://free-proxy.cz/en/proxylist/country/FR/https/ping/all")
	}
}

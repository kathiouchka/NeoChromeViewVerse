package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func scrapeProxynova() {

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	i := 0
	var ipAddress string
	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, trElem *colly.HTMLElement) {
			i = 0
			trElem.ForEachWithBreak("td", func(_ int, tdElem *colly.HTMLElement) bool {
				if !strings.Contains(tdElem.Text, "google") {
					i++
					if i == 1 {
						ipAddress = getStringBetween(tdElem.Text, `(`, `)`)
						ipAddress = strings.ReplaceAll(ipAddress, " ", "")
						ipAddress = strings.ReplaceAll(ipAddress, "'", "")
						ipAddress = strings.ReplaceAll(ipAddress, "+", "")
						ipAddress = strings.ReplaceAll(ipAddress, `"`, "")
						ipAddress += ":"
					}
					if i == 2 {
						ipAddress += (strings.ReplaceAll(strings.ReplaceAll(tdElem.Text, " ", ""), "\n", ""))
						fmt.Println(ipAddress)
						return false
					}
				}
				return true
			})
		})
	})

	c.Visit("https://proxynova.com/proxy-server-list/country-fr")
}

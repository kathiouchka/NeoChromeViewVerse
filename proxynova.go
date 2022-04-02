package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func scrapeProxynova() {
	f, err := os.Create("data.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	i := -1
	y := 0
	j := 0
	var ipAddress string

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, trElem *colly.HTMLElement) {
			y++
		})
	})

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, trElem *colly.HTMLElement) {
			j++
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

					if i == 2 && j != (y-1) {
						ipAddress += (strings.ReplaceAll(strings.ReplaceAll(tdElem.Text, " ", ""), "\n", ""))
						n, err := fmt.Fprintln(f, ipAddress)

						if err != nil {

							log.Fatal(err)
						}

						fmt.Println(n, "bytes written")
						fmt.Println("done")
						return false
					} else if i == 2 && j == (y-1) {
						ipAddress += (strings.ReplaceAll(strings.ReplaceAll(tdElem.Text, " ", ""), "\n", ""))
						n, err := fmt.Fprint(f, ipAddress)

						if err != nil {

							log.Fatal(err)
						}

						fmt.Println(n, "bytes written")
						fmt.Println("done")
						return false
					}

				}
				return true
			})
		})

	})
	fmt.Println(y)
	c.Visit("https://proxynova.com/proxy-server-list/country-fr")

}

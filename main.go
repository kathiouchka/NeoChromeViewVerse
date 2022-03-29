package main

import (
	b64 "encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

func getStringBetween(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return value
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return value
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

// func CharIsPartOfIPAddress(s string) bool {
//     for i, r := range s {
//         if (r >= 48 && r <= 57) || (r == ':') || (r == '.') {
// 			i++
//         }
// 		else {
// 			&s = &s+1
// 		}
//     }
//     return true
// }

func testScrapeFreeProxy() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
		// fmt.Println(string(r.Body))
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
	c.Visit("http://free-proxy.cz/en/proxylist/country/FR/https/ping/all")
}

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
	// testScrapeFreeProxy()
	scrapeProxynova()
	// scrapeIPInfo()
}

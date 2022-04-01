package main

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	playwright "github.com/playwright-community/playwright-go"
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

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}
func testHeadlessScrapeProxyPlayWright() {
	err := playwright.Install()
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch Chromium: %v", err)
	}
	page, err := browser.NewPage(playwright.BrowserNewContextOptions{
		RecordVideo: &playwright.BrowserNewContextOptionsRecordVideo{
			Dir: playwright.String("videos/"),
		},
	})
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	gotoPage := func(url string) {
		fmt.Printf("Visiting %s\n", url)
		if _, err = page.Goto(url); err != nil {
			log.Fatalf("could not goto: %v", err)
		}
		fmt.Printf("Visited %s\n", url)
	}
	// gotoPage("http://whatsmyuseragent.org")
	// gotoPage("https://github.com")
	// gotoPage("https://microsoft.com")
	gotoPage("https://twitch.tv/kathiouzer")
	if err := page.Close(); err != nil {
		log.Fatalf("failed to close page: %v", err)
	}
	path, err := page.Video().Path()
	if err != nil {
		log.Fatalf("failed to get video path: %v", err)
	}
	fmt.Printf("Saved to %s\n", path)
	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}

func testHeadlessScrapeProxyChromeDP() {

	opts := append(chromedp.DefaultExecAllocatorOptions[:]) // 1) specify the proxy server.
	// Note that the username/password is not provided here.
	// Check the link below for the description of the proxy settings:
	// https://www.chromium.org/developers/design-documents/network-settings
	// chromedp.ProxyServer("103.86.50.186:8080"),
	// By default, Chrome will bypass localhost.
	// The test server is bound to localhost, so we should add the
	// following flag to use the proxy for localhost URLs.
	// chromedp.Flag("proxy-bypass-list", "<-loopback>"),

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))
	defer cancel()

	start := time.Now()
	var res string
	var buf []byte
	myurl := "https://www.twitch.tv/kathiouzer"
	err := chromedp.Run(ctx,
		chromedp.Navigate(myurl),
		chromedp.ScrollIntoView(`h1`),
		chromedp.WaitVisible(`h1`),
		fullScreenshot(myurl, 90, &buf),
		chromedp.Text(`h1`, &res, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err = ioutil.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")

	fmt.Printf("h1 contains: '%s'\n", res)
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
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
	// testHeadlessScrapeProxyChromeDP()
	testHeadlessScrapeProxyPlayWright()
	// scrapeProxynova()
	// scrapeIPInfo()
}

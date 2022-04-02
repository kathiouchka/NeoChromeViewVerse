package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func scrapeIPInfo() {
	pURL, _ := url.Parse(`http://193.31.27.123:80`)
	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(pURL)}}
	resp, err := httpClient.Get("http://www.ipinfo.io")
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
	// testHeadlessScrapeProxyPlayWright()
	scrapeProxynova()
	// scrapeIPInfo()
}

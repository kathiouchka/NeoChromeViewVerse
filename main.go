package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {

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

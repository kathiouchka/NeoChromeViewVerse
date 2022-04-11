package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func scrapeIPInfo(IPAddress string) (string, bool) {
	bool := false
	pURL, _ := url.Parse(`http://` + IPAddress)
	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(pURL)}, Timeout: 40 * time.Second}
	resp, err := httpClient.Get("http://www.ipinfo.io")
	if err != nil {
		log.Println(err)
		bool = true
		return "", bool
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		bool = true
		return "", bool
	}
	fmt.Println(string(b))
	if string(b)[0] != '{' {
		bool = true
		return "", bool
	}
	if strings.Contains(string(b), "error") {
		bool = true
		return "", bool
	}
	return IPAddress, false
}

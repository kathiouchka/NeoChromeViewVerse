package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type data struct {
	// Id string `json:"_id"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
	// AnonymityLevel string `json:"anonymityLevel"
	//  "asn":"AS55507",
	//  "city":"Jaipur",
	//  "country":"IN",
	//  "created_at":"2022-04-06T12:38:15.805Z",
	//  "google":false,
	//  "isp":"Uclix",
	//  "lastChecked":1649250249,
	//  "latency":255,
	//  "org":"",
	//  "port":"83",
	//  "protocols":[
	//     "http"
}
type ipaddress struct {
	Data []data
}

func geonodeScrapper() {

	dataLength := 0
	var proxyList string
	url := "https://proxylist.geonode.com/api/proxy-list?limit=50&page=1&sort_by=lastChecked&sort_type=desc"

	f, err := os.Create("geonodeproxies.txt")
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}

	ipaddress1 := ipaddress{}
	jsonErr := json.Unmarshal(body, &ipaddress1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for i, _ := range ipaddress1.Data {

		dataLength = i
	}

	for i, e := range ipaddress1.Data {
		if i <= dataLength {
			proxyList += e.Ip
			proxyList += ":"
			proxyList += e.Port
			proxyList += "\n"
			i++
		} else {
			proxyList += e.Ip
			proxyList += ":"
			proxyList += e.Port
		}
	}

	n, err := fmt.Fprintln(f, proxyList)

	if err != nil {

		log.Fatal(err)
	}

	fmt.Println(n, "bytes written")
	fmt.Print(proxyList)
}

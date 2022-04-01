package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	playwright "github.com/playwright-community/playwright-go"
)

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
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

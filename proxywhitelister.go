package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func proxyWhiteLister() {
	file, err := os.Open("/home/badakzz/NeoChromeViewVerse/data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		str, ok := scrapeIPInfo(scanner.Text())
		if !ok {
			f, err := os.Create("pogggg.txt")

			if err != nil {
				log.Fatal(err)
			}

			defer f.Close()
			_, err = fmt.Fprintln(f, str)

			if err != nil {

				log.Fatal(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

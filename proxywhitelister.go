package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func proxyWhiteLister(inputPath string, outputPath string) {
	// file, err := os.Open("/home/badakzz/NeoChromeViewVerse/data.txt")
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	f, err := os.Create(outputPath)

	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		str, ok := scrapeIPInfo(scanner.Text())

		defer f.Close()

		if !ok {

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

package main

import (
	"fmt"
	"os"
)

func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 1 {
		fmt.Println("no website provided")
		return
	}

	if len(argsWithoutProg) > 1 {
		fmt.Println("too many arguments provided")
		return
	}

	rawBaseURL := argsWithoutProg[0]
	fmt.Printf("starting crawl of: %v...\n", rawBaseURL)

	pages := make(map[string]int)

	crawlPage(rawBaseURL, rawBaseURL, pages)
	for normalizedURL, count := range pages {
		fmt.Printf("%d - %s\n", count, normalizedURL)
	}
}

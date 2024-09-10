package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 1 {
		fmt.Println("Not enough arguments provided")
		fmt.Println("Usage: crawler <baseURL> <maxConcurrency> <maxPages>")
		return
	}

	if len(argsWithoutProg) > 3 {
		fmt.Println("Too many arguments provided")
		return
	}

	rawBaseURL := argsWithoutProg[0]
	maxConcurrency, err := strconv.Atoi(argsWithoutProg[1])
	if err != nil {
		fmt.Printf("Error - maxConcurrency: convert arg to config: %v", err)
		return
	}
	maxPages, err := strconv.Atoi(argsWithoutProg[2])
	if err != nil {
		fmt.Printf("Error - maxPages: convert arg to config: %v", err)
		return
	}
	cfg, err := configure(rawBaseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("Error - configure: %v", err)
		return
	}

	fmt.Printf("Starting crawl of: %v...\n", rawBaseURL)

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	cfg.printReport()
}

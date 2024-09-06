package main

import (
	"fmt"
	"os"
)

func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(argsWithoutProg) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := argsWithoutProg[0]
	fmt.Printf("starting crawl of: %v...\n", rawBaseURL)

	htmlBody, err := getHTML(rawBaseURL)
	if err != nil {
		fmt.Printf("unable to get html: %v", err)
	}
	fmt.Print(htmlBody)
}

package main

import (
	"fmt"
	"sort"
)

type Page struct {
	URL   string
	Count int
}

// sorts the pages map and prints to CLI
func (cfg *config) printReport() {
	fmt.Printf(`
==========================================================
  REPORT for %v
==========================================================
`, cfg.baseURL)

	sortedPages := sortPages(cfg.pages)
	for _, page := range sortedPages {
		url := page.URL
		count := page.Count
		fmt.Printf("Found %d internal links to %v\n", count, url)
	}
}

func sortPages(pages map[string]int) []Page {
	slice := []Page{}
	for url, count := range pages {
		slice = append(slice, Page{
			URL:   url,
			Count: count,
		})
	}
	sort.Slice(slice, func(i, j int) bool {
		if slice[i].Count == slice[j].Count {
			return slice[i].URL < slice[j].URL
		}
		return slice[i].Count > slice[j].Count
	})
	return slice
}

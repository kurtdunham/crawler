package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func normalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse URL: %v", err)
	}

	if u.Scheme == "" || u.Host == "" {
		return "", errors.New("invalid URL")
	}
	normalized := u.Host + u.Path
	normalized = strings.ToLower(normalized)
	normalized = strings.TrimSuffix(normalized, "/")

	return normalized, nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	htmlReader := strings.NewReader(htmlBody)
	doc, err := html.Parse(htmlReader)
	if err != nil {
		return []string{}, fmt.Errorf("unable to parse html body: %s", err)
	}

	urls := []string{}
	var traverseNodes func(*html.Node)
	traverseNodes = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, a := range node.Attr {
				if a.Key == "href" {
					href, err := url.Parse(a.Val)
					if err != nil {
						fmt.Printf("couldn't parse href '%v': %v\n", a.Val, err)
						continue
					}

					resolvedURL := baseURL.ResolveReference(href)
					urls = append(urls, resolvedURL.String())
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverseNodes(child)
		}
	}

	traverseNodes(doc)

	return urls, nil
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("got Network error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return "", fmt.Errorf("got HTTP error: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("got non-HTML response: %s", contentType)
	}

	htmlBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("couldn't read response body: %v", err)
	}

	return string(htmlBodyBytes), nil
}

// In the first call to crawlPage(), rawCurrentURL is a copy of rawBaseURL,
// but as we make further HTTP requests to all the URLs we find on the rawBaseURL,
// the rawCurrentURL value will change while the base stays the same
func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	// first check rawCurrentURL is on same domain as rawBaseURL
	// if not, return current pages

	fmt.Printf("Parsing URLs: Current: %s, Base: %s\n", rawCurrentURL, rawBaseURL)
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawCurrentURL, err)
		return
	}
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawBaseURL, err)
		return
	}
	// Skip other websites
	if currentURL.Hostname() != baseURL.Hostname() {
		fmt.Print("current URL domain and root URL domain do not match\n")
		return
	}
	fmt.Printf("Domains match: %s\n", currentURL.Hostname())

	// Attempt to normalize the URL and track page counts
	fmt.Printf("Normalizing URL: %s\n", rawCurrentURL)
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error normalizing url: %v", err)
		return
	}
	// Increment the page count and log the change
	if _, visited := pages[normalizedURL]; visited {
		// We've crawled this page before, increment and be done
		pages[normalizedURL]++
		return
	}
	// First encounter with the page, initialize with count 1
	pages[normalizedURL] = 1

	fmt.Printf("Pages: %s has been crawled %d time(s).\n", normalizedURL, pages[normalizedURL])

	// print the HTML from the current URL to watch crawler in real-time
	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error reading HTML: %v", err)
		return
	}
	fmt.Print(htmlBody)

	// extract urls from html response
	nextURLs, err := getURLsFromHTML(htmlBody, rawBaseURL)
	if err != nil {
		fmt.Printf("error parsing URLS from HTML: %v", err)
		return
	}

	// crawl through  each url
	for _, url := range nextURLs {
		fmt.Printf("Crawling URL: %s\n", url)
		crawlPage(rawBaseURL, url, pages)
		fmt.Print("finished crawling URL \n")
	}
}

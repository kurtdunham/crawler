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

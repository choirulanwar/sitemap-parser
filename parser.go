// Package sitemapparser provides functionality to parse XML sitemaps and extract URLs.
// It supports both regular sitemaps and sitemap index files with concurrent processing.
package sitemapparser

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// SitemapIndex represents the XML structure of a sitemap index file
// that contains references to other sitemap files
type SitemapIndex struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	Sitemaps []Sitemap `xml:"sitemap"`
}

// Sitemap represents a single sitemap entry within a sitemap index
type Sitemap struct {
	Loc string `xml:"loc"` // URL location of the sitemap
}

// URLSet represents the XML structure of a regular sitemap file
// containing a collection of URLs
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

// URL represents a single URL entry in a sitemap
type URL struct {
	Loc string `xml:"loc"` // URL location of the page
}

// fetchSitemapContent retrieves XML content from a given URL
// Returns the raw bytes of the sitemap content or an error if the fetch fails
func fetchSitemapContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// isSitemapIndexFile determines if the given XML content represents a sitemap index
// Returns true if the content can be parsed as a sitemap index
func isSitemapIndexFile(body []byte) bool {
	var sitemapIndex SitemapIndex
	err := xml.Unmarshal(body, &sitemapIndex)
	return err == nil
}

// parseSitemapIndex retrieves and parses a sitemap index file from the given URL
// Returns a pointer to SitemapIndex structure or an error if parsing fails
func parseSitemapIndex(url string) (*SitemapIndex, error) {
	body, err := fetchSitemapContent(url)
	if err != nil {
		return nil, err
	}

	var sitemapIndex SitemapIndex
	err = xml.Unmarshal(body, &sitemapIndex)
	if err != nil {
		return nil, err
	}

	return &sitemapIndex, nil
}

// parseURLSet retrieves and parses a regular sitemap file from the given URL
// Returns a pointer to URLSet structure or an error if parsing fails
func parseURLSet(url string) (*URLSet, error) {
	body, err := fetchSitemapContent(url)
	if err != nil {
		return nil, err
	}

	var urlset URLSet
	err = xml.Unmarshal(body, &urlset)
	if err != nil {
		return nil, err
	}

	return &urlset, nil
}

// ExtractURLs processes a sitemap URL and returns all URLs found
// It handles both sitemap index files and regular sitemaps
// For sitemap index files, it processes all sub-sitemaps concurrently
func ExtractURLs(url string) ([]string, error) {
	body, err := fetchSitemapContent(url)
	if err != nil {
		fmt.Println("Error fetching sitemap content:", err)
		return nil, err
	}

	var urls []string

	if isSitemapIndexFile(body) {
		// Handle sitemap index file
		sitemapIndex, err := parseSitemapIndex(url)
		if err != nil {
			fmt.Println("Error parsing sitemap index:", err)
			return nil, err
		}

		var wg sync.WaitGroup
		urlsChan := make(chan string, len(sitemapIndex.Sitemaps))

		// Process each sitemap concurrently
		for _, sitemap := range sitemapIndex.Sitemaps {
			wg.Add(1)
			go func(sitemapURL string) {
				defer wg.Done()
				urlset, err := parseURLSet(sitemapURL)
				if err != nil {
					fmt.Println("Error parsing sitemap:", err)
					return
				}
				for _, url := range urlset.URLs {
					urlsChan <- url.Loc
				}
			}(sitemap.Loc)
		}

		// Wait for all goroutines to complete
		go func() {
			wg.Wait()
			close(urlsChan)
		}()

		// Collect URLs from channel
		for url := range urlsChan {
			urls = append(urls, url)
		}
	} else {
		// Handle regular sitemap file
		urlset, err := parseURLSet(url)
		if err != nil {
			fmt.Println("Error parsing sitemap:", err)
			return nil, err
		}

		for _, url := range urlset.URLs {
			urls = append(urls, url.Loc)
		}
	}

	return urls, nil
}

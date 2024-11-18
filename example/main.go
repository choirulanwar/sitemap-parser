package main

import (
	"fmt"
	"log"

	sitemapparser "github.com/choirulanwar/sitemap-parser"
)

func main() {
	// Example 1: Parse a regular sitemap
	urls, err := sitemapparser.ExtractURLs("https://example.com/sitemap.xml")
	if err != nil {
		log.Printf("Error parsing regular sitemap: %v", err)
	} else {
		fmt.Printf("Found %d URLs in regular sitemap\n", len(urls))
		for i, url := range urls {
			if i < 5 { // Print first 5 URLs as example
				fmt.Printf("- %s\n", url)
			}
		}
	}

	// Example 2: Parse a sitemap index
	urls, err = sitemapparser.ExtractURLs("https://example.com/sitemap_index.xml")
	if err != nil {
		log.Printf("Error parsing sitemap index: %v", err)
	} else {
		fmt.Printf("\nFound %d total URLs across all sitemaps\n", len(urls))
		for i, url := range urls {
			if i < 5 { // Print first 5 URLs as example
				fmt.Printf("- %s\n", url)
			}
		}
	}
}

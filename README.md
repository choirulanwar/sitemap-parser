# sitemap-parser

Simple XML sitemap parser in Go. Supports both regular sitemaps and sitemap index files.

## Install

```bash
go get github.com/choirulanwar/sitemap-parser
```

## Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/choirulanwar/sitemap-parser"
)

func main() {
    urls, err := sitemapparser.ExtractURLs("https://example.com/sitemap.xml")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d URLs\n", len(urls))
    for _, url := range urls {
        fmt.Println(url)
    }
}
```

## License

MIT

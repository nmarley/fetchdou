package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	dou "github.com/nmarley/fetchdou"
)

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36"
var client = &http.Client{}

func main() {
	usage := fmt.Sprintf("usage: %s <YYYY-MM-DD>", os.Args[0])
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}
	strDate := os.Args[1]
	t, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	golden := dou.NewDOUFetcher(&userAgent)
	links, err := golden.FetchPDFDownloadLinks(t)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Found %d PDFs for date %v\n", len(links), strDate)

	// Note: I removed concurrency from the downloads as I think the server in
	// Brasília might be rate limiting them.

	for _, pdfURL := range links {
		fn, err := suggestedFilename(pdfURL)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("Fetching ", fn)

		r, err := golden.FetchPDF(pdfURL)
		if err != nil {
			fmt.Print(err)
			return
		}

		f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Print(err)
			return
		}

		n, err := io.Copy(f, r)
		if err != nil {
			fmt.Print(err)
			return
		}
		r.Close()

		fmt.Printf("Wrote %d bytes to %v\n", n, fn)
	}
}

func suggestedFilename(pdfURL string) (string, error) {
	parsedURL, err := url.Parse(pdfURL)
	if err != nil {
		return "", err
	}

	return filepath.Base(parsedURL.Path), nil
}

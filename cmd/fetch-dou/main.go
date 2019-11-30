package main

import (
	"io"
	"log"
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
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <YYYY-MM-DD>", os.Args[0])
	}
	strDate := os.Args[1]
	t, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal("see-ya")

	golden := dou.NewDOUFetcher(&userAgent)

	links, err := golden.FetchPDFDownloadLinks(t)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("links:", links)

	// Note: I removed concurrency from the downloads as I think the server in
	// Brasília might be rate limiting them.

	for _, pdfURL := range links {
		log.Println("fetching: ", pdfURL)
		fn, err := suggestedFilename(pdfURL)
		if err != nil {
			log.Print(err)
			return
		}
		log.Print("fn:", fn)

		r, err := golden.FetchPDF(pdfURL)
		if err != nil {
			log.Print(err)
			return
		}

		f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Print(err)
			return
		}

		n, err := io.Copy(f, r)
		if err != nil {
			log.Print(err)
			return
		}

		log.Printf("Wrote %d bytes to %v\n", n, fn)
		log.Println("Fetched PDF to %v", fn)
	}
	log.Println("All done!")
}

func suggestedFilename(pdfURL string) (string, error) {
	parsedURL, err := url.Parse(pdfURL)
	if err != nil {
		return "", err
	}

	return filepath.Base(parsedURL.Path), nil
}

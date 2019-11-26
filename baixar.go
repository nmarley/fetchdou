package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36"
var client = &http.Client{}

// fetchPDFDownloadLinks will knock on the sacred door of bullshit Java servers
// to retrieve the PDF links with special parameters to allow for PDF downloads
// from the sacred (piece of shit) download server in Brasília. But only from
// 12:00 - 23:59, because... servers need to sleep? I don't know, these people
// are morons.
func fetchPDFDownloadLinks(date time.Time) ([]string, error) {
	links := []string{}

	params := searchParams(date)
	postData := url.Values{}
	for k, v := range params {
		postData.Set(k, v)
	}
	postURL := "http://pesquisa.in.gov.br/imprensa/core/jornalList.action"
	req, err := http.NewRequest("POST", postURL, strings.NewReader(postData.Encode()))
	if err != nil {
		return links, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	for k, v := range requestHeaders() {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return links, err
	}
	defer resp.Body.Close()

	var r io.Reader
	r = resp.Body
	if resp.Header.Get("content-encoding") == "gzip" {
		fmt.Println("Gzip detected -- decompressing")
		r, err = gzip.NewReader(resp.Body)
		if err != nil {
			return links, err
		}
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	data := buf.Bytes()

	// for debug purposes ...
	// f, err := os.OpenFile("PDFpage.html", os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// 	return links, err
	// }
	// f.Write(data)

	// parse body
	links = parseLinks(data)

	return links, nil
}

// fetchPDF accepts a URL and filename and downloads the PDF file
func downloadPDF(theURL, filename string) error {
	req, err := http.NewRequest("GET", theURL, nil)
	if err != nil {
		return err
	}
	for k, v := range requestHeaders() {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r io.Reader
	r = resp.Body
	if resp.Header.Get("content-encoding") == "gzip" {
		fmt.Println("Gzip detected -- decompressing")
		zr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		defer zr.Close()
		r = zr
	}

	// var buf bytes.Buffer
	// io.Copy(&buf, r)
	// data := buf.Bytes()

	// io.Reader

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote %d bytes to %s\n", n, filename)
	return nil
}

func main() {
	// date := time.Now()
	// date, _ := time.Parse("2006-01-02", "2019-11-25")
	date, _ := time.Parse("2006-01-02", "2019-01-14")

	links, err := fetchPDFDownloadLinks(date)
	if err != nil {
		panic(err)
	}
	fmt.Println("links:", links)

	// concurrent download (fetch) of PDFs
	maxRoutines := 8
	var wg sync.WaitGroup
	guard := make(chan struct{}, maxRoutines)
	defer close(guard)

	for _, link := range links {
		wg.Add(1)
		guard <- struct{}{}
		go func(pdfURL string) {
			fn, err := suggestedFilename(pdfURL)
			// fmt.Println("suggested filename:", fn)
			err = downloadPDF(pdfURL, fn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err)
				return
			}
			fmt.Printf("Downloaded %v to %v\n", pdfURL, fn)
			wg.Done()
			<-guard
		}(link)
	}
	wg.Wait()
	fmt.Printf("All done!")
}

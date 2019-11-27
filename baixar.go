package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36"
var client = &http.Client{}

// fetchPDFDownloadLinks will knock on the sacred door of bullshit Java servers
// to retrieve the PDF links with special parameters to allow for PDF downloads
// from the sacred (piece of shit) download server in Brasília. But only from
// 12:00 - 23:59, because... servers need to sleep? I don't know, these people
// are morons.
//
// The date value is a time.Time, but only uses the year, month and day.
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
		log.Println("gzip detected, decompressing")
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

// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.

// FetcherFunc is a func which fetches data
type FetcherFunc func(io.Reader) error

// fetchPDF accepts a URL and filename and downloads the PDF file
// func downloadPDF(theURL, filename string) error
func fetchPDF(theURL string, fetch FetcherFunc) error {
	log.Print("fetching: ", theURL)
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
		log.Print("gzip detected, decompressing")
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

	err = fetch(r)
	if err != nil {
		return err
	}

	// f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()
	// n, err := io.Copy(f, r)
	// if err != nil {
	// 	return err
	// }
	// log.Printf("wrote %d bytes to %s\n", n, filename)

	return nil
}

func main() {
	// date := time.Now()
	// date, _ := time.Parse("2006-01-02", "2019-11-25")
	// date, _ := time.Parse("2006-01-02", "2019-01-14")
	date, _ := time.Parse("2006-01-02", "2019-01-18")

	links, err := fetchPDFDownloadLinks(date)
	if err != nil {
		panic(err)
	}
	log.Println("links:", links)

	// AWS stuff
	s3Region := os.Getenv("S3_REGION")
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Region == "" || s3Bucket == "" {
		log.Fatal("error: S3_REGION or S3_BUCKET not set")
	}
	s3Prefix := fmt.Sprintf("%04d/%02d/%02d", date.Year(), date.Month(), date.Day())
	log.Println("s3 region:", s3Region)
	log.Println("s3 bucket:", s3Bucket)
	log.Println("s3 prefix:", s3Prefix)

	// Establish AWS session for s3
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3Region),
	})
	if err != nil {
		log.Fatal("error establishing AWS session:", err)
	}
	log.Println("established an AWS session")

	// concurrent download (fetch) of PDFs
	limit := 8
	maxGoroutines := int(math.Min(float64(limit), float64(len(links))))
	log.Println("maxGoroutines:", maxGoroutines)
	var wg sync.WaitGroup
	guard := make(chan struct{}, maxGoroutines)
	defer close(guard)

	for _, link := range links {
		wg.Add(1)
		guard <- struct{}{}
		go func(pdfURL string) {
			log.Println("goroutine: fetching ", pdfURL)
			fn, err := suggestedFilename(pdfURL)
			if err != nil {
				log.Println("goroutine: err ", err)
				return
			}
			err = fetchPDF(pdfURL, func(r io.Reader) error {
				var buf bytes.Buffer
				_, err := io.Copy(&buf, r)
				if err != nil {
					return err
				}
				r2 := bytes.NewReader(buf.Bytes())
				s3Key := fmt.Sprintf("%s/%s", s3Prefix, fn)
				err = S3Put(r2, s3Bucket, s3Key, "public-read", sess)
				if err != nil {
					return err
				}
				log.Printf("uploaded %v to s3://%v/%v\n", fn, s3Bucket, s3Key)
				return nil
			})
			if err != nil {
				// later error (default log package does not have debug levels)
				log.Println("error:", err)
				return
			}
			// log.Printf("downloaded %v to %v\n", pdfURL, fn)
			wg.Done()
			<-guard
		}(link)
	}
	wg.Wait()
	log.Println("All done!")
}

package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36"
var client = &http.Client{}

func main() {
	date := time.Now()
	// date, _ := time.Parse("2006-01-02", "2019-01-18")

	dou := NewDOUFetcher(&userAgent)

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

	links, err := dou.FetchPDFDownloadLinks(date)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("links:", links)

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
				log.Println("goroutine err: ", err)
				return
			}
			data, err := dou.FetchPDF(pdfURL)
			if err != nil {
				log.Println("goroutine err: ", err)
				return
			}
			s3Key := fmt.Sprintf("%s/%s", s3Prefix, fn)
			err = S3Put(bytes.NewReader(data), s3Bucket, s3Key, "public-read", sess)
			if err != nil {
				log.Println("goroutine err: ", err)
				return
			}
			log.Printf("uploaded %v to s3://%v/%v\n", fn, s3Bucket, s3Key)
			wg.Done()
			<-guard
		}(link)
	}
	wg.Wait()
	log.Println("All done!")
}

func suggestedFilename(pdfURL string) (string, error) {
	parsedURL, err := url.Parse(pdfURL)
	if err != nil {
		return "", err
	}

	return filepath.Base(parsedURL.Path), nil
}

package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36"
var client = &http.Client{}

func main() {
	date := time.Now()
	// date, _ := time.Parse("2006-01-02", "2019-11-28")

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

	// Note: I removed concurrency from the downloads as I think the server in
	// Brasilia might be rate limiting them.

	for _, pdfURL := range links {
		log.Println("fetching: ", pdfURL)
		fn, err := suggestedFilename(pdfURL)
		if err != nil {
			log.Println("err: ", err)
			return
		}
		log.Println("fn:", fn)
		data, err := dou.FetchPDF(pdfURL)
		if err != nil {
			log.Println("err: ", err)
			return
		}
		log.Println("fetched PDF")
		s3Key := fmt.Sprintf("%s/%s", s3Prefix, fn)
		log.Println("s3Key:", s3Key)
		err = S3Put(bytes.NewReader(data), s3Bucket, s3Key, "public-read", sess)
		if err != nil {
			log.Println("err: ", err)
			return
		}
		log.Printf("uploaded %v to s3://%v/%v\n", fn, s3Bucket, s3Key)
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

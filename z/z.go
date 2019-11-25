package main

import (
	"fmt"
	"net/url"
	"path/filepath"
)

func parseURL(theURL string) error {
	parsedURL, err := url.Parse(theURL)
	if err != nil {
		return err
	}
	fmt.Println("parsedURL:", parsedURL)
	fmt.Println("parsedURL.Query():", parsedURL.Query())
	fmt.Println("parsedURL.Path:", parsedURL.Path)

	fmt.Println("base:", filepath.Base(parsedURL.Path))

	return nil
}

func suggestedFilename(pdfURL string) (string, error) {
	parsedURL, err := url.Parse(pdfURL)
	if err != nil {
		return "", err
	}

	return filepath.Base(parsedURL.Path), nil
}

// 2018
// 2019
//   01
//   11
//     21
// 2020
//   01
//     01
// sha256sums of each PDF
// index.html of the whole thing
// /sgpub/do/secao1/2019/2019_11_21/2019_11_21_ASSINADO_do1.pdf

// Note: All this to be done in SLS.
//
// Code for parsing the .gov.br code to get PDF links could be in another Go
// package (which also has a command-line tool for downloading for a given
// day).
//
// Needs to have CloudFront simply for the caching if nothing else.
//
// Let's do a s3 structure of YYYY/MM/DD/FILENAME.pdf
//
// W/every file laid down, do a scan of the "directory" and create an
// index.html
//   (This will be a Lambda triggered by the s3 put)
//
// One option is to use DynamoDB for a metadata store. Can keep sha256sums of
// each PDF and also assist in the index.html for each "dir".

func main() {
	link := "http://download.in.gov.br/sgpub/do/secao1/2019/2019_11_21/2019_11_21_ASSINADO_do1.pdf?arg1=kvb16gCssmwGX0riHXHe9A&arg2=1574739296"
	fn, err := suggestedFilename(link)
	if err != nil {
		panic(err)
	}
	fmt.Println("suggested filename:", fn)

	parseURL(link)
}

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
	"time"
)

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36"
var client = &http.Client{}

// decompressPage
func decompressPage(data []byte) []byte {
	var buf bytes.Buffer
	zr, _ := gzip.NewReader(bytes.NewReader(data))
	defer zr.Close()
	io.Copy(&buf, zr)
	return buf.Bytes()
}

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

	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return links, err
	//}

	// TODO: Check response headers before assuming it's gzipped
	// Gzip?  This page is compressed (GOOD, faster wire xfer)
	zr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return links, err
	}
	defer zr.Close()

	var buf bytes.Buffer
	io.Copy(&buf, zr)
	data := buf.Bytes()

	// for debug purposes ...
	f, err := os.OpenFile("PDFpage.html", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return links, err
	}
	f.Write(data)

	// parse body
	links = parseLinks(data)

	return links, nil
}

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

	// write to file (or s3, or whatever...)
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote %d bytes to %s\n", n, filename)
	return nil
}

func main() {
	// page, err := getFirstPage()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("page: ", string(page))

	//link := "http://download.in.gov.br/sgpub/do/secao1/2019/2019_11_21/2019_11_21_ASSINADO_do1.pdf?arg1=kvb16gCssmwGX0riHXHe9A&arg2=1574739296"
	//err := getPDF(link)
	//if err != nil {
	//	panic(err)
	//}

	// date := time.Now()
	date, _ := time.Parse("2006-01-02", "2019-11-25")
	//fmt.Printf("date: %+v\n", date)
	//panic("hi")

	params := searchParams(date)
	fmt.Printf("params: %+v\n", params)

	links, err := fetchPDFDownloadLinks(date)
	if err != nil {
		panic(err)
	}
	fmt.Println("links:", links)

	for _, link := range links {
		fn, err := suggestedFilename(link)
		fmt.Println("suggested filename:", fn)
		err = downloadPDF(link, fn)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Downloaded %v to %v\n", link, fn)
	}
}

// http://download.in.gov.br/sgpub/do/secao1/2019/2019_11_21/2019_11_21_ASSINADO_do1.pdf?arg1=kvb16gCssmwGX0riHXHe9A&arg2=1574739296
// http://download.in.gov.br/sgpub/do/secao1/extra/2019/2019_11_21/2019_11_21_ASSINADO_do1_extra_A.pdf?arg1=v5uFuVtjRoHyXTCBrS-ILA&arg2=1574739296

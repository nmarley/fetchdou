package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36"
var client = &http.Client{}

func decompressPage(data []byte) []byte {
	var buf bytes.Buffer
	zr, _ := gzip.NewReader(bytes.NewReader(data))
	defer zr.Close()
	io.Copy(&buf, zr)
	return buf.Bytes()
}

func requestHeaders() map[string]string {
	headers := make(map[string]string)

	headers["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"
	headers["accept-encoding"] = "gzip, deflate"
	headers["accept-language"] = "pt-BR,pt;q=0.9,en-GB;q=0.8,en;q=0.7,en-US;q=0.6"
	headers["cookie"] = "GUEST_LANGUAGE_ID=pt_BR"
	headers["dnt"] = "1"
	// headers["host"] = "download.in.gov.br"
	headers["proxy-connection"] = "keep-alive"
	// headers["referer"] = "http://pesquisa.in.gov.br/imprensa/core/jornalList.action"
	headers["upgrade-insecure-requests"] = "1"
	headers["user-agent"] = userAgent

	return headers
}

func searchParams(date time.Time) map[string]string {
	strDDMM := fmt.Sprintf("%02d/%02d", date.Day(), date.Month())

	params := make(map[string]string)
	params["search-bar"] = ""
	params["tipo-pesquisa"] = "0"
	params["sistema-busca"] = "2"
	params["termo-pesquisado"] = "0"
	params["jornal"] = "do1"
	params["t"] = "com.liferay.journal.model.JournalArticle"
	params["g"] = "68942"
	params["edicao.jornal"] = "1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,600,601,602,603,612,613,614,615,701"
	params["checkbox_edicao.jornal"] = "1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,2,2000,529,525,3,3000,3020,1040,526,530,600,601,602,603,604,605,606,607,608,609,610,611,612,613,614,615,701,702"
	params["__checkbox_edicao.jornal"] = "1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,600,601,602,603,612,613,614,615,701"
	params["__checkbox_edicao.jornal"] = "2,2000,529,525,604,605,606,607,702"
	params["__checkbox_edicao.jornal"] = "3,3000,3020,1040,526,530,608,609,610,611"
	params["edicao.txtPesquisa"] = ""
	params["edicao.jornal_hidden"] = "1,1000,1010,1020,515,521,522,531,535,536,523,532,540,1040,2,2000,529,525,3,3000,3020,1040,526,530,600,601,602,603,604,605,606,607,608,609,610,611,612,613,614,615,701,702"
	params["edicao.dtInicio"] = strDDMM
	params["edicao.dtFim"] = strDDMM
	params["edicao.ano"] = fmt.Sprintf("%04d", date.Year())

	return params
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return links, err
	}

	// Gzip?  This page is compressed (GOOD, faster wire xfer)
	data := decompressPage(body)
	f, err := os.OpenFile("PDFpage.html", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return links, err
	}
	f.Write(data)

	// parse body
	links = parseLinks(data)

	return links, nil
}

func getPDF(theURL string) error {
	// parsedURL, err := url.Parse(theURL)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("parsedURL:", parsedURL)
	// fmt.Println("parsedURL.Query():", parsedURL.Query())

	req, err := http.NewRequest("GET", theURL, nil)
	if err != nil {
		return err
	}
	// req.Header.Set("User-Agent", userAgent)
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

	filename := "out.pdf"

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
}

// http://download.in.gov.br/sgpub/do/secao1/2019/2019_11_21/2019_11_21_ASSINADO_do1.pdf?arg1=kvb16gCssmwGX0riHXHe9A&arg2=1574739296
// http://download.in.gov.br/sgpub/do/secao1/extra/2019/2019_11_21/2019_11_21_ASSINADO_do1_extra_A.pdf?arg1=v5uFuVtjRoHyXTCBrS-ILA&arg2=1574739296

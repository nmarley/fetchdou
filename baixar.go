package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	// "net/url"
)

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36"
var client = &http.Client{}

// http://www.in.gov.br/leiturajornal?data=10-09-2019#daypicker
func getFirstPage() ([]byte, error) {
	var data []byte
	// urlPattern
	url := "http://www.in.gov.br/leiturajornal?data=22-11-2019#daypicker"
	// url := "http://www.in.gov.br/leiturajornal?data=10-09-2019#daypicker"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	return body, nil
}

func getRequestHeaders() map[string]string {
	headers := make(map[string]string)

	headers["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"
	headers["accept-encoding"] = "gzip, deflate"
	headers["accept-language"] = "pt-BR,pt;q=0.9,en-GB;q=0.8,en;q=0.7,en-US;q=0.6"
	headers["cookie"] = "GUEST_LANGUAGE_ID=pt_BR"
	headers["dnt"] = "1"
	headers["host"] = "download.in.gov.br"
	headers["proxy-connection"] = "keep-alive"
	headers["referer"] = "http://pesquisa.in.gov.br/imprensa/core/jornalList.action"
	headers["upgrade-insecure-requests"] = "1"

	return headers
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
	req.Header.Set("User-Agent", userAgent)
	for k, v := range getRequestHeaders() {
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

	link := "http://download.in.gov.br/sgpub/do/secao1/2019/2019_11_21/2019_11_21_ASSINADO_do1.pdf?arg1=kvb16gCssmwGX0riHXHe9A&arg2=1574739296"
	err := getPDF(link)
	if err != nil {
		panic(err)
	}
}

// http://download.in.gov.br/sgpub/do/secao1/2019/2019_11_21/2019_11_21_ASSINADO_do1.pdf?arg1=kvb16gCssmwGX0riHXHe9A&arg2=1574739296
// http://download.in.gov.br/sgpub/do/secao1/extra/2019/2019_11_21/2019_11_21_ASSINADO_do1_extra_A.pdf?arg1=v5uFuVtjRoHyXTCBrS-ILA&arg2=1574739296

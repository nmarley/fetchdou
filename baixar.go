package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

func main() {
	page, err := getFirstPage()
	if err != nil {
		panic(err)
	}
	fmt.Println("page: ", string(page))
}


// URL from JavaScript (This goes to actual download page w/results. Needs to be POSTed apparently, and might need some tag scraped from the initial site)
//http://pesquisa.in.gov.br/imprensa/core/jornalList.action


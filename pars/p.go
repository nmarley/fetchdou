package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"regexp"
)

func parseLinks(data []byte) []string {
	re := regexp.MustCompile(`(http://download.in.gov.br/[^\'"]*)`)
	preLinks := re.FindAllString(string(data), -1)
	links := make([]string, len(preLinks))
	for i, url := range preLinks {
		links[i] = html.UnescapeString(url)
	}
	return links
}

func main() {
	// fn := "page2.html"
	fn := "PDFpage.html"
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	links := parseLinks(data)
	fmt.Println("links:", links)
}

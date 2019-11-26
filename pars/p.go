package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
    "html"
)

func parseLinks(data []byte) []string {
	re := regexp.MustCompile(`(http://download.in.gov.br/[^\'"]*)`)
	preLinks := re.FindAllString(string(data), -1)
    links := make([]string, len(preLinks))
	for i, url := range links {
        fmt.Printf("i = %d, url = %v\n", i, url)
        links[i] = html.UnescapeString(url)
	}
	return links
}

func main() {
	fn := "page2.html"
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	links := parseLinks(data)
	fmt.Println("links:", links)
}

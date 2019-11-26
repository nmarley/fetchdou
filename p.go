package main

import (
	"html"
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

package main

import (
	"fmt"
	"net/http"
)

// http://www.in.gov.br/leiturajornal?data=10-09-2019#daypicker
func getFirstPage() (string, error) {
	urlPattern := "http://www.in.gov.br/leiturajornal?data=10-09-2019#daypicker"
	sth := ""

	return sth, nil
}

func main() {
	getFirstPage()
}

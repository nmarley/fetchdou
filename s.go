package main

import (
	"net/url"
	"path/filepath"
)

func suggestedFilename(pdfURL string) (string, error) {
	parsedURL, err := url.Parse(pdfURL)
	if err != nil {
		return "", err
	}

	return filepath.Base(parsedURL.Path), nil
}

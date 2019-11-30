package fetchdou

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var client = &http.Client{}

// rePDFDownloadLink is the regular expression for parsing out download links
// from the search page.
var rePDFDownloadLink = regexp.MustCompile(`(http://download.in.gov.br/[^\'"]*)`)

// DOUFetcher is an object used for fetching the DOU in PDF format.
type DOUFetcher struct {
	userAgent string
}

// NewDOUFetcher constructs a new DOUFetcher object with a userAgent string.
func NewDOUFetcher(userAgent *string) *DOUFetcher {
	// hard-code a random common user-agent string in case user doesn't pass one
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	if userAgent != nil {
		ua = *userAgent
	}

	return &DOUFetcher{
		userAgent: ua,
	}
}

// requestHeaders returns a map of headers needed to request the page.
func (f *DOUFetcher) requestHeaders() map[string]string {
	headers := make(map[string]string)

	headers["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"
	headers["accept-encoding"] = "gzip, deflate"
	headers["accept-language"] = "pt-BR,pt;q=0.9,en-GB;q=0.8,en;q=0.7,en-US;q=0.6"
	headers["cookie"] = "GUEST_LANGUAGE_ID=pt_BR"
	headers["dnt"] = "1"
	headers["proxy-connection"] = "keep-alive"
	headers["upgrade-insecure-requests"] = "1"
	headers["user-agent"] = f.userAgent

	// unchecked, leave off for now
	// headers["host"] = "download.in.gov.br"
	// headers["referer"] = "http://pesquisa.in.gov.br/imprensa/core/jornalList.action"

	return headers
}

// searchParams returns a map of POST params needed to request the search page.
func (f *DOUFetcher) searchParams(date time.Time) map[string]string {
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

// FetchPDFDownloadLinks returns a list of links to download the DOU in PDF
// format. The links have certain parameters which must remain intact. If the
// request is not made within the window of 12:00 - 23:59 M-F, the server will
// not allow the downloads and instead return an empty list.
//
// The date argument is a time.Time, but only uses the year, month and day.
func (f *DOUFetcher) FetchPDFDownloadLinks(date time.Time) ([]string, error) {
	links := []string{}

	params := f.searchParams(date)
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
	for k, v := range f.requestHeaders() {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return links, err
	}
	defer resp.Body.Close()

	var r io.Reader
	r = resp.Body
	if resp.Header.Get("content-encoding") == "gzip" {
		log.Println("gzip detected, decompressing")
		zr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return links, err
		}
		defer zr.Close()
		r = zr
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	data := buf.Bytes()

	// extract links from body
	links = parseLinks(data)

	return links, nil
}

// parseLinks extracts the PDF download links from the HTML search results.
func parseLinks(data []byte) []string {
	preLinks := rePDFDownloadLink.FindAllString(string(data), -1)
	links := make([]string, len(preLinks))
	for i, url := range preLinks {
		links[i] = html.UnescapeString(url)
	}
	return links
}

// FetchPDF accepts a PDF URL as emitted by FetchPDFDownloadLinks and downloads
// and returns the raw bytes.
func (f *DOUFetcher) FetchPDF(url string) (io.Reader, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range f.requestHeaders() {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r io.Reader
	r = resp.Body
	if resp.Header.Get("content-encoding") == "gzip" {
		r, _ = gzip.NewReader(resp.Body)
	}

	return r, nil
}

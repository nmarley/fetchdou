package main

import (
	"bytes"
	"fmt"

	// "os"

	// "rsc.io/pdf"
	"github.com/ledongthuc/pdf"
)

func main() {
	fn := "2019_11_21_ASSINADO_do1.pdf"
	fmt.Println("file:", fn)

	// f, err := os.Open(fn)
	// if err != nil {
	//     panic(err)
	// }
	// defer f.Close()

	// r, err := pdf.NewReader(f,
	// if err != nil {
	//     panic(err)
	// }

	f, r, err := pdf.Open(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		panic(err)
	}
	buf.ReadFrom(b)

	fmt.Println(buf.String())

	r, err := pdf.Open(fn)
	if err != nil {
		panic(err)
	}

	numPages := r.NumPage()
	fmt.Printf("there are %d pages\n", numPages)

	pageNo := 37
	page := r.Page(pageNo)

	// fmt.Printf("page %d = %v\n", pageNo, p)
	// fmt.Printf("page %d has %d Texts\n", pageNo, len(p.Content().Text))
	//	for i := 1; i <= numPages; i++ {
	//		page := r.Page(i)
	// cont := page.Content()

	keys := page.V.Keys()
	for _, k := range keys {
		fmt.Printf("k = %v\n", k)
	}

	c := page.V.Key("Contents")
	fmt.Printf("c = %v, type: %v\n", c, c.Kind())

	//cont := page.V.Key("Contents")
	//fmt.Printf("page %d has %d Texts\n", i, len(cont.Text))
	//fmt.Printf("page %d has %d Rects\n", i, len(cont.Rect))
	// fmt.Printf("debug: %v\n", page.V.Kind())
	// fmt.Printf("pdf.Dict: %d\n", pdf.Dict)
	//	}
	// r.GetPlainText()

}

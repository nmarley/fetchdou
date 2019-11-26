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
}

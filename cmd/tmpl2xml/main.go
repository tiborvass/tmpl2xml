package main

import (
	"io"
	"log"
	"os"

	"github.com/tiborvass/tmpl2xml"
)

func main() {
	var r io.Reader = os.Stdin
	var err error
	if len(os.Args) > 1 {
		r, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}
	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	out, err := tmpl2xml.String(string(b))
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.WriteString(out)
	os.Stdout.Write([]byte{'\n'})
}

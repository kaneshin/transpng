package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/kaneshin/transpng"
)

var (
	// Usage: transpng [-decode] [-input in_file] [-output out_file]
	decode = flag.Bool("decode", false, "decodes input")
	input  = flag.String("input", "-", `input file (default: "-" for stdin)`)
	output = flag.String("output", "-", `output file (default: "-" for stdout)`)
)

func init() {
	log.SetFlags(0)
	flag.Parse()
}

func main() {
	var reader io.Reader
	if *input != "-" {
		// reads from file for reading.
		f, err := os.Open(*input)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		reader = f
	}

	if reader == nil {
		reader = os.Stdin
	}

	var buf bytes.Buffer
	if *decode {
		// decode
		dec := transpng.NewDecoder(reader)
		if err := dec.Decode(&buf); err != nil {
			log.Fatal(err)
		}
	} else {
		// encode
		enc := transpng.NewEncoder(&buf)
		if err := enc.Encode(reader); err != nil {
			log.Fatal(err)
		}
	}

	if *output != "-" {
		// writes to file.
		if err := ioutil.WriteFile(*output, buf.Bytes(), 0666); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Fprint(os.Stdout, buf.String())
	}
}

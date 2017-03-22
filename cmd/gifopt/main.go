package main

import (
	"flag"
	"image/gif"
	"log"
	"os"

	"github.com/donatj/gifopt"
)

var (
	filename = flag.String("o", "output.gif", "Where to save the optimized gif")
)

func init() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Stderr.WriteString("requires one image as input")
		os.Exit(1)
	}
}

func main() {
	file, err := os.Open(flag.Arg(0))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	g, err := gif.DecodeAll(file)
	if err != nil {
		log.Fatal(err)
	}

	g = gifopt.InterframeCompress(g, 1500000)

	outfile, err := os.Create(*filename)
	defer outfile.Close()
	if err != nil {
		log.Fatal(err)
	}
	gif.EncodeAll(outfile, g)
}

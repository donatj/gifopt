package main

import (
	"flag"
	"fmt"
	"image/gif"
	"log"
	"os"
	"path/filepath"

	"github.com/donatj/gifopt"
)

const defaultFile = "<orig>.opt.gif"

var (
	filename  = flag.String("o", defaultFile, "Where to save the optimized gif")
	threshold = flag.Float64("t", (1500000/float64(gifopt.MaxDistance))*100, "Max interframe color diff percent threshold")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [options] <gif>:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if *filename == defaultFile {
		ext := filepath.Ext(flag.Arg(0))
		path := filepath.Base(flag.Arg(0))
		path = path[0 : len(path)-len(ext)]

		*filename = path + ".opt.gif"
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

	t := (*threshold * gifopt.MaxDistance) / 100
	g = gifopt.InterframeCompress(g, uint32(t))

	outfile, err := os.Create(*filename)
	defer outfile.Close()
	if err != nil {
		log.Fatal(err)
	}
	gif.EncodeAll(outfile, g)
}

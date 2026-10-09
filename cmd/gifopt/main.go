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
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	g, err := gif.DecodeAll(file)
	if err != nil {
		log.Fatal(err)
	}

	t := (*threshold * gifopt.MaxDistance) / 100
	g = gifopt.InterframeCompress(g, uint32(t))

	outfile, err := os.Create(*filename)
	if err != nil {
		log.Fatal(err)
	}

	if err := gif.EncodeAll(outfile, g); err != nil {
		_ = outfile.Close()
		log.Fatal(err)
	}

	if err := outfile.Close(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"flag"
	"fmt"
	"image/gif"
	"log"
	"os"
	"path/filepath"

	pip "github.com/JamesMilnerUK/pip-go"
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

	polys := []gifopt.PolygonThreshold{
		{
			Threshold: distFromPercent(0),
			Polygon: pip.Polygon{
				Points: []pip.Point{{157, 78}, {194, 62}, {287, 57}, {328, 64}, {335, 110}, {317, 178}, {267, 219}, {222, 219}, {182, 178}, {158, 115}, {157, 78}},
			},
		},
	}

	g = gifopt.InterframeCompress(g, distFromPercent(*threshold), polys)

	outfile, err := os.Create(*filename)
	defer outfile.Close()
	if err != nil {
		log.Fatal(err)
	}
	gif.EncodeAll(outfile, g)
}

func distFromPercent(percent float64) uint32 {
	return uint32((percent * gifopt.MaxDistance) / 100)
}

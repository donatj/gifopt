package main

import (
	"flag"
	"fmt"
	"image/gif"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	pip "github.com/JamesMilnerUK/pip-go"
	"github.com/donatj/gifopt"
)

const defaultFile = "<orig>.opt.gif"

var (
	filename  = flag.String("o", defaultFile, "Where to save the optimized gif")
	threshold = flag.Float64("t", (1500000/float64(gifopt.MaxDistance))*100, "Max interframe color diff percent threshold")
	regions   = flag.String("regions", "", "polygonal threshold regions")
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

	polys, err := parseRegionString(*regions)
	if err != nil {
		log.Fatal(err)
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

func parseRegionString(regionStr string) ([]gifopt.PolygonThreshold, error) {
	polys := []gifopt.PolygonThreshold{}
	if regionStr != "" {
		poly := strings.Split(regionStr, "|")
		for polyI, e := range poly {
			pt := gifopt.PolygonThreshold{
				Polygon: pip.Polygon{
					Points: []pip.Point{},
				},
			}

			polyparts := strings.SplitN(e, ":", 2)
			if len(polyparts) != 2 {
				return nil, fmt.Errorf("failed to parse region polygon %d", polyI)
			}
			t, err := strconv.ParseFloat(polyparts[0], 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse region polygon %d", polyI)
			}
			pt.Threshold = distFromPercent(t)
			points := strings.Split(polyparts[1], ";")
			for pointI, point := range points {
				p := strings.SplitN(point, ",", 2)
				if len(p) != 2 {
					return nil, fmt.Errorf("failed to parse region polygon %d point %d", polyI, pointI)
				}

				x, err1 := strconv.ParseFloat(p[0], 64)
				y, err2 := strconv.ParseFloat(p[1], 64)
				if err1 != nil || err2 != nil {
					return nil, fmt.Errorf("failed to parse region polygon %d point %d", polyI, pointI)
				}

				pt.Polygon.Points = append(pt.Polygon.Points, pip.Point{X: x, Y: y})

			}

			//ensure every polygon is closed
			pt.Polygon.Points = append(pt.Polygon.Points, pt.Polygon.Points[0])

			polys = append(polys, pt)
		}
	}

	return polys, nil
}

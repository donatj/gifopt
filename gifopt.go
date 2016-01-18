package main

import (
	"flag"
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"
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
	if err != nil {
		log.Fatal(err)
	}

	g, err := gif.DecodeAll(file)
	if err != nil {
		log.Fatal(err)
	}

	b := g.Image[0].Bounds()

	visible := image.NewRGBA(b)

	for i, img := range g.Image {
		transInd := -1

		p := color.Palette{}
		for x, pc := range img.Palette {
			_, _, _, pa := pc.RGBA()
			if pa == 0 {
				transInd = x
			}
			p = append(p, pc)
		}
		img.Palette = p

		for y := 0; y <= b.Max.Y; y++ {
			for x := 0; x <= b.Max.X; x++ {
				c := img.At(x, y)
				_, _, _, a := c.RGBA()

				if i > 0 {
					v := visible.At(x, y)
					dd := Dist(c, v)

					if dd < 5000000 {
						if transInd == -1 {
							transInd = img.Palette.Index(c)

							img.Palette[transInd] = color.Transparent
						}

						img.SetColorIndex(x, y, uint8(transInd))
					}
				}

				if a != 0 {
					visible.Set(x, y, c)
				}
			}
		}
	}

	outfile, err := os.Create("output.gif")
	if err != nil {
		log.Fatal(err)
	}
	gif.EncodeAll(outfile, g)

	return
}

func Dist(c color.Color, v color.Color) uint32 {
	// A batch version of this computation is in image/draw/draw.go.
	cr, cg, cb, ca := c.RGBA()
	vr, vg, vb, va := v.RGBA()

	sum := sqDiff(cr, vr) + sqDiff(cg, vg) + sqDiff(cb, vb) + sqDiff(ca, va)

	return sum
}

func sqDiff(x, y uint32) uint32 {
	var d uint32
	if x > y {
		d = x - y
	} else {
		d = y - x
	}
	return (d * d) >> 2
}

package gifopt

import (
	"image"
	"image/color"
	"image/gif"
)

const (
	// MaxDistance is the result of calculating dist(color.White, color.Transparent)
	// It is slightly less than MaxUint32 and shouldn't be swaped for it.
	MaxDistance = 4294836224
)

// InterframeCompress helps optimize gifs by analyzing colors across frames and
// setting pixels to transparent if they are below a threshold difference.
//
// InterframeCompress is a lossy method of removing pixels to allowing gifs LZW
// compression to better do its job.
func InterframeCompress(g *gif.GIF, limit uint32) *gif.GIF {
	if len(g.Image) < 2 {
		return g
	}

	visible := image.NewRGBA(g.Image[0].Bounds())

	for i, img := range g.Image {
		sb := img.Bounds()
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

		// Some strange gifs have frames that don't start at the origin…
		for y := sb.Min.Y; y < sb.Max.Y; y++ {
			for x := sb.Min.X; x < sb.Max.X; x++ {
				c := img.At(x, y)
				_, _, _, a := c.RGBA()

				if i > 0 {
					v := visible.At(x, y)
					dd := dist(c, v)

					if dd < limit {
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

	return g
}

func dist(c color.Color, v color.Color) uint32 {
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

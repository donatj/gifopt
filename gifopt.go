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

		compressFrame(img, visible, i > 0, limit, transInd)
	}

	return g
}

func dist(c color.Color, v color.Color) uint32 {
	// A batch version of this computation is in image/draw/draw.go.
	cr, cg, cb, ca := c.RGBA()
	vr, vg, vb, va := v.RGBA()

	return distRGBA(cr, cg, cb, ca, vr, vg, vb, va)
}

func distRGBA(cr, cg, cb, ca, vr, vg, vb, va uint32) uint32 {
	return sqDiff(cr, vr) + sqDiff(cg, vg) + sqDiff(cb, vb) + sqDiff(ca, va)
}

func compressFrame(img *image.Paletted, visible *image.RGBA, compare bool, limit uint32, transInd int) {
	sb := img.Bounds()
	if !sb.In(visible.Bounds()) {
		compressFrameSlow(img, visible, compare, limit, transInd)
		return
	}

	for y := sb.Min.Y; y < sb.Max.Y; y++ {
		imageOffset := img.PixOffset(sb.Min.X, y)
		visibleOffset := visible.PixOffset(sb.Min.X, y)

		for x := 0; x < sb.Dx(); x++ {
			c := img.Palette[img.Pix[imageOffset+x]]
			cr, cg, cb, ca := c.RGBA()

			pixelOffset := visibleOffset + x*4
			if compare {
				vr := uint32(visible.Pix[pixelOffset])
				vg := uint32(visible.Pix[pixelOffset+1])
				vb := uint32(visible.Pix[pixelOffset+2])
				va := uint32(visible.Pix[pixelOffset+3])
				vr |= vr << 8
				vg |= vg << 8
				vb |= vb << 8
				va |= va << 8

				if distRGBA(cr, cg, cb, ca, vr, vg, vb, va) < limit {
					if transInd == -1 {
						transInd = img.Palette.Index(c)
						img.Palette[transInd] = color.Transparent
					}

					img.Pix[imageOffset+x] = uint8(transInd)
				}
			}

			if ca != 0 {
				setVisiblePixel(visible.Pix[pixelOffset:pixelOffset+4], c, cr, cg, cb, ca)
			}
		}
	}
}

func compressFrameSlow(img *image.Paletted, visible *image.RGBA, compare bool, limit uint32, transInd int) {
	sb := img.Bounds()

	// Some strange gifs have frames that don't start at the origin or extend
	// outside the first frame's bounds. Keep the image package's bounds
	// handling for those uncommon cases.
	for y := sb.Min.Y; y < sb.Max.Y; y++ {
		for x := sb.Min.X; x < sb.Max.X; x++ {
			c := img.At(x, y)

			if compare && dist(c, visible.At(x, y)) < limit {
				if transInd == -1 {
					transInd = img.Palette.Index(c)
					img.Palette[transInd] = color.Transparent
				}

				img.SetColorIndex(x, y, uint8(transInd))
			}

			if _, _, _, a := c.RGBA(); a != 0 {
				visible.Set(x, y, c)
			}
		}
	}
}

func setVisiblePixel(pixel []uint8, c color.Color, r, g, b, a uint32) {
	if c, ok := c.(color.RGBA); ok {
		pixel[0] = c.R
		pixel[1] = c.G
		pixel[2] = c.B
		pixel[3] = c.A
		return
	}

	pixel[0] = uint8(r >> 8)
	pixel[1] = uint8(g >> 8)
	pixel[2] = uint8(b >> 8)
	pixel[3] = uint8(a >> 8)
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

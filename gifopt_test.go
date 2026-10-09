package gifopt

import (
	"image"
	"image/color"
	"image/gif"
	"reflect"
	"testing"
)

func TestMaxDistance(t *testing.T) {
	d := dist(color.White, color.Transparent)
	if d != MaxDistance {
		t.Errorf(`Dist(color.White, color.Transparent) = %#v; want %#v`, d, MaxDistance)
	}
}

func TestInterframeCompressMatchesSlowPath(t *testing.T) {
	for _, limit := range []uint32{0, 1, 100, MaxDistance} {
		got := testGIF()
		want := testGIF()

		InterframeCompress(got, limit)
		interframeCompressSlow(want, limit)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("InterframeCompress result differs from the reference implementation with limit %d", limit)
		}
	}
}

func interframeCompressSlow(g *gif.GIF, limit uint32) {
	if len(g.Image) < 2 {
		return
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

		compressFrameSlow(img, visible, i > 0, limit, transInd)
	}
}

func testGIF() *gif.GIF {
	palette := color.Palette{
		color.NRGBA{R: 255, A: 128},
		color.RGBA{G: 255, A: 255},
		color.Transparent,
	}
	rect := image.Rect(-1, 2, 2, 4)
	first := image.NewPaletted(rect, palette)
	second := image.NewPaletted(rect, palette)
	copy(first.Pix, []uint8{0, 1, 0, 1, 0, 1})
	copy(second.Pix, []uint8{0, 1, 1, 1, 0, 0})

	return &gif.GIF{Image: []*image.Paletted{first, second}}
}

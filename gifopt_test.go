package gifopt

import (
	"image/color"
	"testing"
)

func TestMaxDistance(t *testing.T) {
	d := dist(color.White, color.Transparent)
	if d != MaxDistance {
		t.Errorf(`Dist(color.White, color.Transparent) = %#v; want %#v`, d, MaxDistance)
	}
}

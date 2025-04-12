// -*- coding: utf-8 -*-
// utils.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 12-04-2025 14:24:35.775566190 (1744460675)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// This package contains a few utilities that can be used with golor
package utils

import (
	"iter"
	"math"
)

// Given three bytes with R, G and B values, return the Hue, Saturation and
// Lightness of its combination
func RgbToHsl(r, g, b uint8) (h, s, l float64) {

	rf := float64(r)
	gf := float64(g)
	bf := float64(b)

	rf /= 0xff
	gf /= 0xff
	bf /= 0xff

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	l = (max + min) / 2

	if max == min {
		h, s = 0, 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2.0 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case rf:
			h = (gf - bf) / d
			if gf < bf {
				h += 6
			}
		case gf:
			h = (bf-rf)/d + 2
		case bf:
			h = (rf-gf)/d + 4
		}
		h /= 6
	}
	return
}

// Given the Hue, Saturation and Lightness of a combination of red, green and
// blue, return its three primary colors as bytes
func HslToRgb(h, s, l float64) (r, g, b uint8) {
	var rF, gF, bF float64

	if s == 0 {
		rF, gF, bF = l, l, l
	} else {
		var hue2rgb = func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			if t < 1.0/6 {
				return p + (q-p)*6*t
			}
			if t < 1.0/2 {
				return q
			}
			if t < 2.0/3 {
				return p + (q-p)*(2.0/3-t)*6
			}
			return p
		}

		q := l * (1 + s)
		if l >= 0.5 {
			q = l + s - l*s
		}
		p := 2*l - q
		rF = hue2rgb(p, q, h+1.0/3)
		gF = hue2rgb(p, q, h)
		bF = hue2rgb(p, q, h-1.0/3)
	}

	return uint8(rF * 0xff), uint8(gF * 0xff), uint8(bF * 0xff)
}

// Use the HSL model to create a pleasant gradient of color with the given
// number of steps from the start to the specified end. Note that the start and
// end consist of a combination of red, green and blue and thus, they are given
// as uint32
func HslGradient(startRGB, endRGB uint32, steps int) iter.Seq2[int, uint32] {

	rstart := uint8((startRGB >> 16) & 0x0000ff)
	gstart := uint8((startRGB >> 8) & 0x0000ff)
	bstart := uint8(startRGB & 0x0000ff)

	rend := uint8((endRGB >> 16) & 0x0000ff)
	gend := uint8((endRGB >> 8) & 0x0000ff)
	bend := uint8(endRGB & 0x0000ff)

	h1, s1, l1 := RgbToHsl(rstart, gstart, bstart)
	h2, s2, l2 := RgbToHsl(rend, gend, bend)

	return func(yield func(int, uint32) bool) {

		for i := range steps {
			t := float64(i) / float64(steps-1)

			h := h1 + (h2-h1)*t
			s := s1 + (s2-s1)*t
			l := l1 + (l2-l1)*t

			r, g, b := HslToRgb(h, s, l)
			color := (uint32(r) << 16) | (uint32(g) << 8) | uint32(b)
			if !yield(i, color) {
				return
			}
		}
	}
}

// Local Variables:
// mode:go
// fill-column:80
// End:

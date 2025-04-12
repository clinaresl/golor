package main

import (
	"math"

	"github.com/clinaresl/golor"
)

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

func hslGradient(startRGB, endRGB uint32, steps int) (gradient []uint32) {

	rstart := uint8((startRGB >> 16) & 0x0000ff)
	gstart := uint8((startRGB >> 8) & 0x0000ff)
	bstart := uint8(startRGB & 0x0000ff)

	rend := uint8((endRGB >> 16) & 0x0000ff)
	gend := uint8((endRGB >> 8) & 0x0000ff)
	bend := uint8(endRGB & 0x0000ff)

	h1, s1, l1 := RgbToHsl(rstart, gstart, bstart)
	h2, s2, l2 := RgbToHsl(rend, gend, bend)

	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)

		h := h1 + (h2-h1)*t
		s := s1 + (s2-s1)*t
		l := l1 + (l2-l1)*t

		r, g, b := HslToRgb(h, s, l)
		color := (uint32(r) << 16) | (uint32(g) << 8) | uint32(b)
		gradient = append(gradient, color)
	}

	return gradient
}

// fade-in the given string using the combination of colors red (r), green (g)
// and blue in the foreground without affecting the background
func fadeInFg(str string, r, g, b bool) {

	// for all characters in the given string
	for idx, val := range str {

		// compute a value from 0 to 0xff, so that the first character gets the
		// value 0x00, and the last one gets the value 0xff
		sat := uint32(0xff*idx) / uint32(len(str))

		// now, compute the foreground color using the given combination of red, green and blue
		fg := uint32(0)
		if r {
			fg |= sat << 16
		}
		if g {
			fg |= sat << 8
		}
		if b {
			fg |= sat
		}
		golor.Printf("%C{%c}", fg, val)
	}
	golor.Printf("\n")
}

// fade-in the given string using the combination of colors red (r), green (g)
// and blue in the background without affecting the background
func fadeInBg(str string, r, g, b bool) {

	// for all characters in the given string
	for idx, val := range str {

		// compute a value from 0 to 0xff, so that the first character gets the
		// value 0x00, and the last one gets the value 0xff
		sat := uint64(0xff*idx) / uint64(len(str))

		// now, compute the foreground color using the given combination of red, green and blue
		bg := uint64(0)
		if r {
			bg |= sat << 40
		}
		if g {
			bg |= sat << 32
		}
		if b {
			bg |= sat << 24
		}
		golor.Printf("%C{%c}", bg, val)
	}
	golor.Printf("\n")
}

func fadeIn(str string, from, to uint32) {

	for idx, val := range hslGradient(from, to, len(str)) {
		golor.Printf("%C{%c}", val, str[idx])
	}

	golor.Printf("\n")
}

func main() {

	// Take a random sentence, yeah latin would be nice :)
	ipsum := "Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Vestibulum tortor quam, feugiat vitae, ultricies eget, tempor sit amet, ante. Donec eu libero sit amet quam egestas semper. Aenean ultricies mi vitae est. Mauris placerat eleifend leo. Quisque sit amet est et sapien ullamcorper pharetra. Vestibulum erat wisi, condimentum sed, commodo vitae, ornare sit amet, wisi. Aenean fermentum, elit eget tincidunt condimentum, eros ipsum rutrum orci, sagittis tempus lacus enim ac dui. Donec non enim in turpis pulvinar facilisis. Ut felis. Praesent dapibus, neque id cursus faucibus, tortor neque egestas augue, eu vulputate magna eros eu erat. Aliquam erat volutpat. Nam dui mi, tincidunt quis, accumsan porttitor, facilisis luctus, metus"

	// Foreground
	fadeInFg(ipsum, true, false, false)
	fadeInFg(ipsum, false, true, false)
	fadeInFg(ipsum, false, false, true)
	fadeInFg(ipsum, true, true, false)
	fadeInFg(ipsum, true, false, true)
	fadeInFg(ipsum, false, true, true)
	fadeInFg(ipsum, true, true, true)

	// Background
	fadeInBg(ipsum, true, false, false)
	fadeInBg(ipsum, false, true, false)
	fadeInBg(ipsum, false, false, true)
	fadeInBg(ipsum, true, true, false)
	fadeInBg(ipsum, true, false, true)
	fadeInBg(ipsum, false, true, true)
	fadeInBg(ipsum, true, true, true)

	fadeIn(ipsum, 0xff00ff, 0xffff00)
}

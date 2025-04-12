package main

import (
	"fmt"
	"iter"
	"math"

	"github.com/clinaresl/golor"
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
func hslGradient(startRGB, endRGB uint32, steps int) iter.Seq2[int, uint32] {

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

// Print the given string with a pleasant foreground gradient from the start
// combination of red, green and blue until the specified end
func fadeInForeground(str string, start, end uint32) {

	for idx, val := range hslGradient(start, end, len(str)) {
		golor.Printf("%C{%c}", val, str[idx])
	}

	golor.Printf("\n")
}

// Print the given string with a pleasant background gradient from the start
// combination of red, green and blue until the specified end.
func fadeInBackground(str string, start, end uint32) {

	for idx, val := range hslGradient(start, end, len(str)) {
		golor.Printf("%C{%c}", uint64(val)<<24, str[idx])
	}

	golor.Printf("\n")
}

// Print the given string with a pleasant gradient from the start combination of
// red, green and blue until the specified end. The gradient is computed for the
// foreground and the background is the opposite
func fadeInForegroundBackground(str string, start, end uint32) {

	for idx, val := range hslGradient(start, end, len(str)) {
		golor.Printf("%C{%c}", uint64(val^0x00ffffff)<<24|uint64(val), str[idx])
	}

	golor.Printf("\n")
}

func main() {

	// Take a random sentence, yeah latin would be nice :)
	ipsum := "Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas."

	// Show the string in standard face for the sake of comparison
	fmt.Println(ipsum)
	fmt.Println()

	// Properties with Foreground
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.BOLD}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.DIM}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.ITALIC}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.UNDERLINE}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.SLOW_BLINK}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.RAPID_BLINK}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.CROSSED_OUT}, ipsum)
	fmt.Println()

	// Properties with Background
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.BOLD}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.DIM}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.ITALIC}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.UNDERLINE}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.SLOW_BLINK}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.RAPID_BLINK}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.CROSSED_OUT}, ipsum)
	fmt.Println()

	// Foreground and background colors
	fadeInForeground(ipsum, 0x000000, 0xff0000)
	fadeInForeground(ipsum, 0x000000, 0x00ff00)
	fadeInForeground(ipsum, 0x000000, 0x0000ff)

	fadeInForeground(ipsum, 0x000000, 0xffffff)
	fmt.Println()

	fadeInBackground(ipsum, 0x000000, 0xff0000)
	fadeInBackground(ipsum, 0x000000, 0x00ff00)
	fadeInBackground(ipsum, 0x000000, 0x0000ff)

	fadeInBackground(ipsum, 0x000000, 0xffffff)
	fmt.Println()

	fadeInForegroundBackground(ipsum, 0x000000, 0xff0000)
	fadeInForegroundBackground(ipsum, 0x000000, 0x00ff00)
	fadeInForegroundBackground(ipsum, 0x000000, 0x0000ff)

	fadeInForegroundBackground(ipsum, 0x000000, 0xffffff)
	fmt.Println()
}

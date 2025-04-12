// Small showcase of the capabilities of the golor module.
//
// This demo is not intended to show all possibilities, and it has been mostly
// created to test the output in a specific console.
//
// It also uses the services provided in the utils package of the golor module
package main

import (
	"fmt"

	"github.com/clinaresl/golor"
	"github.com/clinaresl/golor/utils"
)

// Print the given string with a pleasant foreground gradient from the start
// combination of red, green and blue until the specified end
func fadeInForeground(str string, start, end uint32) {

	for idx, val := range utils.HslGradient(start, end, len(str)) {
		golor.Printf("%C{%c}", val, str[idx])
	}

	golor.Printf("\n")
}

// Print the given string with a pleasant background gradient from the start
// combination of red, green and blue until the specified end.
func fadeInBackground(str string, start, end uint32) {

	for idx, val := range utils.HslGradient(start, end, len(str)) {
		golor.Printf("%C{%c}", uint64(val)<<24, str[idx])
	}

	golor.Printf("\n")
}

// Print the given string with a pleasant gradient from the start combination of
// red, green and blue until the specified end. The gradient is computed for the
// foreground and the background is the opposite
func fadeInForegroundBackground(str string, start, end uint32) {

	for idx, val := range utils.HslGradient(start, end, len(str)) {
		golor.Printf("%C{%c}", uint64(val^0x00ffffff)<<24|uint64(val), str[idx])
	}

	golor.Printf("\n")
}

func main() {

	// Take a random sentence, yeah latin would be nice :)
	ipsum := "Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas."

	// Show the string in standard face for the sake of comparison
	fmt.Println(ipsum)
	fmt.Println("---")

	// Properties with Foreground
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.BOLD}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.DIM}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.ITALIC}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.UNDERLINE}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.SLOW_BLINK}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.RAPID_BLINK}, ipsum)
	golor.Printf("%C{%v}\n", golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.CROSSED_OUT}, ipsum)
	fmt.Println("---")

	// Properties with Background
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.BOLD}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.DIM}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.ITALIC}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.UNDERLINE}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.SLOW_BLINK}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.RAPID_BLINK}, ipsum)
	golor.Printf("%C{%v}\n", golor.BgEffect{R: 0x20, G: 0x00, B: 0x80, Properties: golor.CROSSED_OUT}, ipsum)
	fmt.Println("---")

	// Foreground and background colors
	fadeInForeground(ipsum, 0x000000, 0xff0000)
	fadeInForeground(ipsum, 0x000000, 0x00ff00)
	fadeInForeground(ipsum, 0x000000, 0x0000ff)

	fadeInForeground(ipsum, 0x000000, 0xffffff)
	fmt.Println("---")

	fadeInBackground(ipsum, 0x000000, 0xff0000)
	fadeInBackground(ipsum, 0x000000, 0x00ff00)
	fadeInBackground(ipsum, 0x000000, 0x0000ff)

	fadeInBackground(ipsum, 0x000000, 0xffffff)
	fmt.Println("---")

	fadeInForegroundBackground(ipsum, 0x000000, 0xff0000)
	fadeInForegroundBackground(ipsum, 0x000000, 0x00ff00)
	fadeInForegroundBackground(ipsum, 0x000000, 0x0000ff)

	fadeInForegroundBackground(ipsum, 0x000000, 0xffffff)
	fmt.Println("---")
}

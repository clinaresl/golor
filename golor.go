// -*- coding: utf-8 -*-
// golor.go
// -----------------------------------------------------------------------------
//
// Started on <vie 29-09-2023 22:38:25.455020113 (1696019905)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// golor provides means for showing text on the standard output with different
// selections of color for both the background and the foreground, and also
// other properties. It is intended to ease usage.
//
// golor is strongly based on the interface provided by the fmt Print family
// functions. It provides a new verb %C{} which encloses the string between
// curly brackets (which is welcome to contain other verbs as well) with an ANSI
// prefix and suffix that produce the desired effect
package golor

import (
	"fmt"
	"log"
	"regexp"
)

// Constants
// ----------------------------------------------------------------------------

// The following constants define the different components of the ASCII escape
// codes used for showing properties
const (
	prefix             = "\033["
	foreground_prefix  = "38;2"
	background_prefix  = "48;2"
	bold_prefix        = "1"
	dim_prefix         = "2"
	italic_prefix      = "3"
	underline_prefix   = "4"
	slow_blink_prefix  = "5"
	rapid_blink_prefix = "6"
	crossed_out_prefix = "9"
	suffix             = "\033[0m"
)

// Constants used to filter the red, green and blue colors of the foreground and
// background colors when using either uint32 or uint64 types
const (

	// Specification with uin32
	properties32 = 0xff000000
	fg_red32     = 0xff0000
	fg_green32   = 0x00ff00
	fg_blue32    = 0x0000ff

	// Specification with uin64
	properties64 = 0x00ff000000000000
	bg_red64     = 0xff0000000000
	bg_green64   = 0x00ff00000000
	bg_blue64    = 0x0000ff000000
)

// Constants that can be used for defining properties with the types [Effect],
// [FgEffect] and [BgEffect]
const (
	BOLD = 1 << (iota + 0)
	DIM
	ITALIC
	UNDERLINE
	SLOW_BLINK
	RAPID_BLINK
	CROSSED_OUT
)

// Constants that can be used for defining properties with the type [Effect32]
const (
	BOLD32 = 1 << (iota + 24)
	DIM32
	ITALIC32
	UNDERLINE32
	SLOW_BLINK32
	RAPID_BLINK32
	CROSSED_OUT32
)

// Constants that can be used for defining properties with the type [Effect64]
const (
	BOLD64 = 1 << (iota + 48)
	DIM64
	ITALIC64
	UNDERLINE64
	SLOW_BLINK64
	RAPID_BLINK64
	CROSSED_OUT64
)

// Provide a map between properties and their sequence
var propertyPrefix = map[uint8]string{
	BOLD:        bold_prefix,
	DIM:         dim_prefix,
	ITALIC:      italic_prefix,
	UNDERLINE:   underline_prefix,
	SLOW_BLINK:  slow_blink_prefix,
	RAPID_BLINK: rapid_blink_prefix,
	CROSSED_OUT: crossed_out_prefix,
}

// The following regular expression is used for matching any verb, thouse used
// in the fmt Printf family function, and also the verb %C{...}
const all_verbs_regexp = `%(C\{([^\}]+)\}|(?<flags>[-+#0 ])?(?<width>\d+|\*)?(?:\.(?<precision>\d+|\*))?(?<length>[hljztL]|hh|ll)?(?<specifier>[diuoxXfFeEgGaAcspnTv]))`

// The following regular expression is used instead for matching color codes
// only %C{...}
const color_regexp = `^%C\{([^\}]+)\}`

// Types
// ----------------------------------------------------------------------------

// The following type defines an RGB color to be used only with type [Effect]
type Color struct {
	R, G, B uint8
}

// The following type defines a combination of foreground, background colors and
// properties. Note that both the foreground and background colors have to be of
// type [Color]
type Effect struct {
	Fg, Bg     Color
	Properties uint8
}

// The following type defines a combination of foreground color and properties.
// The foregrround color must be given with three different bytes
type FgEffect struct {
	R, G, B    uint8
	Properties uint8
}

// The following type defines a combination of background color and properties.
// The backgrround color must be given with three different bytes
type BgEffect struct {
	R, G, B    uint8
	Properties uint8
}

// It is also possible to define just the foreground color and the properties
// using an uint32
type Effect32 = uint32

// It is also possible to define the background and foreground colors (in that
// order) and the properties using an uint64
type Effect64 = uint64

// Variables
// ----------------------------------------------------------------------------

// (Must)Compiled regexps
var allVerbs = regexp.MustCompile(all_verbs_regexp)
var colorVerb = regexp.MustCompile(color_regexp)

// Functions
// ----------------------------------------------------------------------------

// Process the specified color specification and return a string representing
// the foreground color. It returns an error in case the specification is given
// in an unknown format
func processForegroundColor(arg any) (output string, err error) {

	// This package supports various formats for specifying colors
	switch val := arg.(type) {

	case Effect:

		output = fmt.Sprintf("%v;%v;%v", val.Fg.R, val.Fg.G, val.Fg.B)

	case FgEffect:

		output = fmt.Sprintf("%v;%v;%v", val.R, val.G, val.B)

	case BgEffect:

		// This type does not provide information about the foreground color
		break

	case Effect32:

		output = fmt.Sprintf("%v;%v;%v", (val&fg_red32)>>16, (val&fg_green32)>>8, val&fg_blue32)

	case Effect64:

		output = fmt.Sprintf("%v;%v;%v", (val&fg_red32)>>16, (val&fg_green32)>>8, val&fg_blue32)

	default:
		return "", fmt.Errorf("Unsuported foreground color format: %v\n", arg)
	}

	return
}

// Process the specified color specification and return a string representing
// the background color. It returns an error in case the specification is given
// in an unknown format
func processBackgroundColor(arg any) (output string, err error) {

	// This package supports various formats for specifying colors
	switch val := arg.(type) {

	case Effect:

		output = fmt.Sprintf("%v;%v;%v", val.Bg.R, val.Bg.G, val.Bg.B)

	case FgEffect:

		// This type does not provide information about the background color
		break

	case BgEffect:

		output = fmt.Sprintf("%v;%v;%v", val.R, val.G, val.B)

	case Effect32:

		// This type does not provide information about the background color
		break

	case Effect64:

		output = fmt.Sprintf("%v;%v;%v", (val&bg_red64)>>40, (val&bg_green64)>>32, val&bg_blue64>>24)

	default:
		return "", fmt.Errorf("Unsuported background color format: %v\n", arg)
	}

	return
}

// Process the specified properties and return the string with its ANSI codes
func processProperties(properties uint8) (output string) {

	// Process all properties one by one
	var idx uint8
	for idx = BOLD; idx <= CROSSED_OUT; idx <<= 1 {

		if properties&idx != 0 {
			output += fmt.Sprintf(";%v", propertyPrefix[idx])
		}
	}

	return
}

// Given a string chunk, return it preceded by the color prefix corresponding to
// the given argument and ended with the corresponding suffix
func substituteColorVerb(chunk string, arg any) (output string, err error) {

	// This package supports various formats for specifying colors and
	// properties
	switch val := arg.(type) {

	case Effect:

		// Get the foreground and background specs
		fg, fgerr := processForegroundColor(val)
		bg, bgerr := processBackgroundColor(val)
		if fgerr != nil {
			return "", err
		}
		if bgerr != nil {
			return "", err
		}

		output = fmt.Sprintf(`%v%v;%v;%v;%v%vm%v%v`, prefix, foreground_prefix, fg, background_prefix, bg, processProperties(val.Properties), chunk, suffix)

	case FgEffect:

		// Get the foreground spec
		fg, fgerr := processForegroundColor(val)
		if fgerr != nil {
			return "", err
		}

		output = fmt.Sprintf(`%v%v;%v%vm%v%v`, prefix, foreground_prefix, fg, processProperties(val.Properties), chunk, suffix)

	case BgEffect:

		// Get the background specs
		bg, bgerr := processBackgroundColor(val)
		if bgerr != nil {
			return "", err
		}

		output = fmt.Sprintf(`%v%v;%v%vm%v%v`, prefix, background_prefix, bg, processProperties(val.Properties), chunk, suffix)

	case Effect32:

		// Get the foreground spec
		fg, fgerr := processForegroundColor(val)
		if fgerr != nil {
			return "", err
		}

		output = fmt.Sprintf(`%v%v;%v%vm%v%v`, prefix, foreground_prefix, fg, processProperties(uint8((val&properties32)>>24)), chunk, suffix)

	case Effect64:

		// Get the foreground and background specs
		fg, fgerr := processForegroundColor(val)
		bg, bgerr := processBackgroundColor(val)
		if fgerr != nil {
			return "", err
		}
		if bgerr != nil {
			return "", err
		}

		output = fmt.Sprintf(`%v%v;%v;%v;%v%vm%v%v`, prefix, foreground_prefix, fg, background_prefix, bg, processProperties(uint8(val&properties64)>>48), chunk, suffix)

	default:
		return "", fmt.Errorf("Unsupported format: %v\n", arg)
	}

	return
}

// substitute all occurrences of color verbs by their corresponding prefixes and
// suffixes without affecting the other verbs in the format string. It returns:
//
//  1. The resulting string with all color verbs properly substituted,
//
//  2. The list of arguments to be used in the substitution of the remaining verbs,
//
//  3. The number of arguments consumed in the substitution of non-color verbs.
//
//  4. An error in case any is found.
func processColorVerbs(format string, a ...any) (output string, args []any, nargs int, err error) {

	// Keep a counter over the arguments given in a to know which ones to use
	// for substituting the color verbs, and which to use in the substitutions
	// performed by the Printf family
	var idx int

	// An offset is needed to know the position from which the last chunk of
	// format was copied
	var offset int

	// Get the location of all verbs and process them all
	matches := allVerbs.FindAllStringSubmatchIndex(format, -1)
	for _, match := range matches {

		// Check whether this is a color verb
		if colorVerb.MatchString(format[match[0]:match[1]]) {

			// First, process the contents of the color verb. Notice that all
			// the remaining args are used, but the first one which must be used
			// later for substituting the color verb
			log.Printf(" [nested process: '%v']\n", format[match[0]+3:match[1]-1])
			contents, cargs, cnargs, cerr := processColorVerbs(format[match[0]+3:match[1]-1], a[idx+1:]...)
			log.Printf("\t contents: %v\n", contents)
			log.Printf("\t cargs   : %v\n", cargs)
			log.Printf("\t cnargs  : %v\n", cnargs)
			log.Printf("\t cerr    : %v\n", cerr)
			if cerr != nil {
				return "", nil, 0, err
			}

			// copy from the previous offset until the end of the color verb,
			// substitute it, and update the offset
			if chunk, err := substituteColorVerb(contents, a[idx]); err == nil {
				output += format[offset:match[0]]
				output += chunk
				offset = match[1]
			} else {
				return "", nil, 0, err
			}

			// Next, copy all the necessary args to make the necessary
			// substitutions later inside the color-verb. 'cnargs' is the number
			// of non-color verbs found within this color verb.
			for i := 0; i < cnargs; i++ {

				// Next, we copy all the arguments necessary to substitute the
				// non-color verbs inside the color verb we just processed. Note
				// we add 1 to avoid cosidering the argument of the color verb
				args = append(args, a[i+1])
			}

			// and update the counter of the next arguments to used. Again, note
			// we add 1 to avoid re-using the argument of the color verb
			idx += 1 + cnargs
		} else {

			// Otherwise, copy all elements in the output and update the offset
			output += format[offset:match[1]]
			offset = match[1]

			// and preserve this argument to be used a posteriori
			nargs++
			args = append(args, a[idx])
			log.Printf("\t args    : %v\n\n", args)

			// and update the counter of the next argument to use.
			idx++
		}
	}

	// Finally, add to the output string the rest of it since the last offset,
	// and also any other argument that have not been used (this is relevant
	// when procesing nested chunks of the original string)
	output += format[offset:]
	if len(a) > idx {
		args = append(args, a[idx:]...)
	}

	return
}

// Cprintf is the counterpart of Printf. It just substitutes the color verbs and
// queries fmt.Printf to resolve the rest. It returns the number of bytes
// written and any write error encountered.
func Printf(format string, a ...any) (n int, err error) {

	// First, substitute all the color verbs
	cformat, cargs, cnargs, cerr := processColorVerbs(format, a...)
	if cerr != nil {
		return 0, err
	}

	// Next, format the resulting string with the arguments that were not
	// used in the color verbs
	log.Printf(" cformat: %+v", cformat)
	log.Printf(" cargs  : %+v\n", cargs)
	log.Printf(" cnargs : %+v\n", cnargs)
	return fmt.Printf(cformat, cargs...)
}

// Local Variables:
// mode:go
// fill-column:80
// End:

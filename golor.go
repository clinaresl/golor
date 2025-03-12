// -*- coding: utf-8 -*-
// golor.go
// -----------------------------------------------------------------------------
//
// Started on <vie 29-09-2023 22:38:25.455020113 (1696019905)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// golor provides means for showing text on the standard color with different
// selections of color for both the background and the foreground, and also
// other effects. It is intended to ease usage.
//
// golor is strongly based on the interface provided by the fmt Print family
// functions. It provides a new verb %C{} which encloses the string between
// curly brackets (which is welcome to contain other verbs as well)
package golor

import (
	"fmt"
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

// Next, the different constants that can be used for defining properties
const (
	BOLD = 1 << (iota + 6)
	DIM
	ITALIC
	UNDERLINE
	SLOW_BLINK
	RAPID_BLINK
	CROSSED_OUT
)

// Provide a map between properties and their sequence
var propertyPrefix = map[uint64]string{
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
const ALL_VERBS_REGEXP = `%(C\{([^\}]+)\}|(?<flags>[-+#0 ])?(?<width>\d+|\*)?(?:\.(?<precision>\d+|\*))?(?<length>[hljztL]|hh|ll)?(?<specifier>[diuoxXfFeEgGaAcspnTv]))`

// The following regular expression is used for matching the different verbs
// used in the fmt Printf family function, excluding the new verb %C{...}
const VERB_REGEXP = `%(?<flags>[-+#0 ])?(?<width>\d+|\*)?(?:\.(?<precision>\d+|\*))?(?<length>[hljztL]|hh|ll)?(?<specifier>[diuoxXfFeEgGaAcspn])`

// The following regular expression is used instead for matching color codes
// only %C{...}
const COLOR_REGEXP = `^%C\{([^\}]+)\}`

// Types
// ----------------------------------------------------------------------------

// The following type defines an effect
type Effect struct {
	R, G, B    uint8
	Properties uint64
}

// Variables
// ----------------------------------------------------------------------------

// (Must)Compiled regexps
var allVerbs = regexp.MustCompile(ALL_VERBS_REGEXP)
var colorVerb = regexp.MustCompile(COLOR_REGEXP)

// Functions
// ----------------------------------------------------------------------------

// Process the specified properties and return the string with its ANSI codes
func processProperties(properties uint64) (output string) {

	// Process all properties one by one
	var idx uint64
	for idx = BOLD; idx <= CROSSED_OUT; idx <<= 1 {

		if properties&idx != 0 {
			output += fmt.Sprintf(";%v", propertyPrefix[idx])
		}
	}

	return
}

// Given a string with a color verb and its argument, return the contents of the
// color verb with its prefix and suffix
func substituteColorVerb(chunk string, arg any) (output string, err error) {

	// This package supports providing the argument in various formats
	switch val := arg.(type) {

	case Effect:

		output = fmt.Sprintf(`%v;%v;%v;%v%vm%v%v`, prefix+foreground_prefix, val.R, val.G, val.B, processProperties(val.Properties), chunk[3:len(chunk)-1], suffix)

	default:
		return "", fmt.Errorf("Unsupported format: %v\n", arg)
	}

	return
}

// substitute all occurrences of color verbs by their corresponding prefixes and
// suffixes without affecting the other verbs in the format string. It returns
// the resulting string with all color verbs properly substituted, and an error
// in case any is found.
func ProcessColorVerbs(format string, a ...any) (output string, err error) {

	// Use the all verbs regexp to locate all verbs. This is important because
	// to substitute the color verbs it is necessary to know precisely the
	// arguments that have to be used
	var idx int

	// An offset is needed to know the position from which the last chunk of
	// format was copied
	var offset int

	// Get the location of all verbs and process them all
	matches := allVerbs.FindAllStringSubmatchIndex(format, -1)
	for _, match := range matches {

		// Check whether this is a color verb
		if colorVerb.MatchString(format[match[0]:match[1]]) {

			// copy from the previous offset until the end of the color verb,
			// substitute it, and update the offset
			if chunk, err := substituteColorVerb(format[match[0]:match[1]], a[idx]); err == nil {
				output += format[offset:match[0]]
				output += chunk
				offset = match[1]
			} else {
				return "", err
			}
		} else {

			// Otherwise, copy all elements in the output and update the offset
			output += format[offset:match[1]]
			offset = match[1]
		}

		// And, in any case increment the counter of parameters
		idx++
	}

	return
}

// Cprintf is the counterpart of Printf. It just substitutes the %C{...}
// sections and queries fmt.Printf to resolve the rest. It returns the number of
// bytes written and any write error encountered.
func Printf(format string, a ...any) (n int, err error) {

	return
}

// Local Variables:
// mode:go
// fill-column:80
// End:

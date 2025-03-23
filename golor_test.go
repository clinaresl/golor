// -*- coding: utf-8 -*-
// golor_test.go
// -----------------------------------------------------------------------------
//
// Started on <dom 16-03-2025 15:29:33.689194401 (1742135373)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package golor

// The following example shows how to write a sentence wtih two words, where
// only the foreground colors and its properties have been set. The first word
// is shown in red and bold, whereas the second appears in blue and underlined.
func Example_general1() {
	Printf("%C{Hello} %C{World}", uint32(0xff0000)|BOLD32, uint32(0x00ff00)|UNDERLINE32)
}

// It is also possible to set both the foreground and background colors along
// with its properties using uint64 as shown below. In this case, the
// combination RGB of the background color has to be given before the foreground
// color and the properties are set using those with the general 64
func Example_general2() {
	Printf("%C{Hello} %C{World}", uint64(0xaadd44ff0000)|BOLD64, uint64(0x43207200ff00)|UNDERLINE64)
}

// Local Variables:
// mode:go
// fill-column:80
// End:

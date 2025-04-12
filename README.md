# Introduction
`golor` provides means for showing text on the standard output with different
selections of color for both the background and the foreground, and also other
effects. It is intended to ease usage

# Usage 

`golor` provides the same interface than the Printf family functions, i.e., it
substitutes all verbs appearing in a string with the given arguments. In
addition, it provides a new verb, `%C{...}`:

+ The argument given to `%C{...}` has to be a color specification. See [Color specification](# Color specification)

+ The ellipsis stands for any text which might also contain other verbs to be
  substituted. Currently, *color verbs* can not be nested and thus the text
  enclosed in a color verb can not contain `%C{...}`
  
# Color specification

This section describes both how to specify the foreground and/or background
color, and also how to set various properties

## Foreground color

The foreground color and properties of a color verb can be set with either a
value of type `golor.FgEffect` or with a `uint32`:

``` go
golor.Printf("%C{%v}\n", 
    golor.FgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.BOLD}, 
    "Hello World!")
golor.Printf("%C{%v}\n", uint32(0xffaa00)|golor.BOLD32, "Hello World!")
```

which produce both the same effect, the foreground is shown in a mixture of red
and green, and the text is shown in bold typeface.

## Background color

Likewise, the background color and properties of some text can be set using
either a value of type `golor.BgEffect` or with an `uint64`:

``` go
	golor.Printf("%C{%v}\n", 
        golor.BgEffect{R: 0xff, G: 0xaa, B: 0x00, Properties: golor.ITALIC}, 
        "Hello World!")
	golor.Printf("%C{%v}\n", uint64(0xffaa00000000)|golor.ITALIC64, "Hello World!")
```

which do not produce exactly the same effect. The reason is that the first
statement does not modify the foreground color, while using an `uint64` for
specifying the background color specifically requires setting the foreground
(which in this case is black, `0x000000`). The reason is that a full combination
of background and foreground color, when given in an `uint64` is decoded as
follows: `rBgBbBrFgFbF` where `r`, `g` and `b` represent the values of red,
green and blue, and `B` and `F` stand for the background and foreground color.

## Foreground and Background

It is also possible to set both the foreground and background colors with any
combination of properties either using a value of type `golor.Effect` or an
`uint64` (See [Background Color](# Background color)):

``` go
	golor.Printf("%C{%v}\n",
		golor.Effect{
			Fg:         golor.Color{R: 0xff, G: 0xaa, B: 0x00},
			Bg:         golor.Color{R: 0x10, G: 0x20, B: 0x30},
			Properties: golor.ITALIC | golor.UNDERLINE},
		"Hello World!")
	golor.Printf("%C{%v}\n", uint64(0x102030ffaa00)|golor.ITALIC64|golor.UNDERLINE64, "Hello World!")
```

This example shows:

+ The convenience of using `uint64` when setting the foreground, background and properties

+ Also, that it is possible to set any combination of properties using the bitwise or operator `|`

## Properties

As mentioned above, `golor` allows the user to specify any combination of
properties, listed below:

+ `golor.BOLD`, `golor.BOLD32`, `golor.BOLD64`
+ `golor.DIM`, `golor.DIM32`, `golor.DIM64`
+ `golor.ITALIC`, `golor.ITALIC32`, `golor.ITALIC64`
+ `golor.UNDERLINE`, `golor.UNDERLINE32`, `golor.UNDERLINE64`
+ `golor.SLOW_BLINK`, `golor.SLOW_BLINK32`, `golor.SLOW_BLINK64`
+ `golor.RAPID_BLINK`, `golor.RAPID_BLINK32`, `golor.RAPID_BLINK64`
+ `golor.CROSSED-OUT`, `golor.CROSSED-OUT32`, `golor.CROSSED-OUT64`

When using values of any of the types `Effect`, `FgEffect` or `BgEffect`, the
properties are set using the first form. The second form must be used only when
using `uint32` for setting the foreground. Lastly, when setting both the
foreground and background with an `uint64`, the third form must be used.

# LICENSE

MIT License

Copyright (c) 2025 Carlos Linares López

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

# Author #

Carlos Linares Lopez <carlos.linares@uc3m.es>  
Computer Science Department <https://www.inf.uc3m.es/en>  
Universidad Carlos III de Madrid <https://www.uc3m.es/home>

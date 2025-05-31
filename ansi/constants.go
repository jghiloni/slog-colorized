package ansi

// ANSIColorCode represents one of the 8 standard ANSI color codes
type ANSIColorCode uint8

const (
	Black   ANSIColorCode = 30
	Red     ANSIColorCode = 31
	Green   ANSIColorCode = 32
	Yellow  ANSIColorCode = 33
	Blue    ANSIColorCode = 34
	Purple  ANSIColorCode = 35
	Cyan    ANSIColorCode = 36
	White   ANSIColorCode = 37
	Default ANSIColorCode = 0
)

// ColorModifier changes whether the color is for the foreground
// or background and whether or not it's the intense variant
type ColorModifier uint8

const (
	Foreground ColorModifier = 0
	Background ColorModifier = 10
	Intense    ColorModifier = 60
)

// StyleModifier set things like bold, italic, etc.
// Not all StyleModifiers are supported by all systems
type StyleModifier uint8

const (
	Normal        StyleModifier = 0
	Bold          StyleModifier = 1
	Italic        StyleModifier = 3
	Underline     StyleModifier = 4
	Strikethrough StyleModifier = 9
)

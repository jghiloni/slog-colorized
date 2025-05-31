package ansi

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/jghiloni/go-commonutils/utils"
)

// Style represents a set of print modifications
type Style struct {
	color *ANSIColorCode
	cm    []ColorModifier
	sm    []StyleModifier
}

// Return the current color, or Default if nil
func (c Style) Color() ANSIColorCode {
	if c.color == nil {
		return Default
	}

	return *c.color
}

// SetColor sets the color. The underlying color is a pointer, and is set to nil if the passed color is Default
func (c *Style) SetColor(color ANSIColorCode) {
	c.color = utils.NilRefIfZero(color)
}

// ColorModifiers returns the normalized (sorted, deduplicated) list of color modifiers
func (c Style) ColorModifiers() []ColorModifier {
	return c.cm
}

// SetColorModifiers sorts and deduplicates the list of color modifiers before setting it
func (c *Style) SetColorModifiers(cm ...ColorModifier) {
	c.cm = sortAndDedupe(cm)
}

// StyleModifiers returns the normalized (sorted, deduplicated) list of style modifiers
func (c Style) StyleModifiers() []StyleModifier {
	return c.sm
}

// SetStyleModifiers sorts and deduplicates the list of style modifiers before setting it
func (c *Style) SetStyleModifiers(sm ...StyleModifier) {
	c.sm = sortAndDedupe(sm)
}

// String implements fmt.Stringer and outputs the ANSI control code for this style
func (c Style) String() string {
	codes := intSlice(c.sm)
	color := 0
	if c.color != nil {
		color = int(*c.color)

		color = utils.Reduce(c.cm, func(c int, m ColorModifier) int {
			return c + int(m)
		}, int(color))
	}

	codes = append(codes, color)
	codeString := strings.Join(utils.Map(codes, strconv.Itoa), ";")

	return fmt.Sprintf("\u001b[%sm", codeString)
}

func (c Style) clone() Style {
	newC := Style{}

	if c.color != nil {
		newColor := *c.color
		newC.color = &newColor
	}

	newC.cm = make([]ColorModifier, len(c.cm))
	newC.sm = make([]StyleModifier, len(c.sm))

	copy(newC.cm, c.cm)
	copy(newC.sm, c.sm)

	return newC
}

func sortAndDedupe[T ~uint8](c []T) []T {
	set := map[T]bool{}
	for _, t := range c {
		set[t] = true
	}

	seq := maps.Keys(set)
	uniq := slices.Collect(seq)

	sort.Slice(uniq, func(i, j int) bool {
		return uniq[i] < uniq[j]
	})

	return uniq
}

func toInt[T ~uint8](t T) int {
	return int(t)
}

func intSlice[T ~uint8](t []T) []int {
	return utils.Map(t, toInt)
}

var reset Style = Style{}

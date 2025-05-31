package ansi

import (
	"fmt"
	"io"
	"reflect"

	"github.com/jghiloni/go-commonutils/utils"
)

// Writer implements io.Writer.
type Writer struct {
	current     Style
	lastWritten *Style
	delegate    io.Writer
}

// NewColorize returns a new io.Writer that can style text
func NewWriter(target io.Writer) *Writer {
	if target == nil {
		return nil
	}

	c := &Writer{
		delegate: target,
		current:  Style{},
	}

	return c
}

func (c Writer) CurrentStyle() Style {
	return c.current
}

func (c *Writer) SetCurrentStyle(ps Style) {
	c.current = ps
}

// Reset sends the reset control code and resets all modifications
func (c *Writer) Reset() (int, error) {
	c.current = reset
	return c.writeControl()
}

// Write implements io.Writer
func (c *Writer) Write(p []byte) (int, error) {
	if len, err := c.writeControl(); err != nil {
		return len, err
	}

	if p == nil {
		return 0, nil
	}

	return c.delegate.Write(p)
}

func (c *Writer) writeControl() (int, error) {
	doWrite := true
	if c.lastWritten != nil {
		doWrite = !reflect.DeepEqual(c.current.clone(), c.lastWritten.clone())
	}

	var len int
	var err error

	if doWrite {
		len, err = fmt.Fprint(c.delegate, c.current)
		if err == nil {
			c.lastWritten = utils.NilRefIfZero(c.current.clone())
		}
	}

	return len, err
}

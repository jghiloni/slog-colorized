package handler

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/jghiloni/stylized/ansi"
	"github.com/mattn/go-isatty"
)

// RecordStylizer returns a PrintStyle for a given slog.Record
type RecordStylizer func(slog.Record) ansi.Style

// Options provides for slog.HandlerOptions, plus
// options around colorizing.
type Options struct {
	*slog.HandlerOptions
	Stylizer      RecordStylizer
	ForceColorize bool
}

type ColorizedHandler struct {
	writer         *ansi.Writer
	delegateBuffer *strings.Builder
	delegate       slog.Handler
	stylizer       RecordStylizer
	mu             *sync.Mutex
}

// Enabled implements slog.Handler.
func (c *ColorizedHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return c.delegate.Enabled(ctx, l)
}

// Handle implements slog.Handler.
func (c *ColorizedHandler) Handle(ctx context.Context, r slog.Record) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.writer.SetCurrentStyle(c.stylizer(r))

	var err error
	if err = c.delegate.Handle(ctx, r); err != nil {
		return err
	}

	out := c.delegateBuffer.String()
	c.delegateBuffer.Reset()

	if _, err = io.WriteString(c.writer, strings.TrimSuffix(out, "\n")); err != nil {
		return err
	}

	if _, err = c.writer.Reset(); err != nil {
		return err
	}

	_, err = c.writer.Write([]byte("\n"))
	return err
}

// WithAttrs implements slog.Handler.
func (c *ColorizedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	c.delegate = c.delegate.WithAttrs(attrs)
	return c
}

// WithGroup implements slog.Handler.
func (c *ColorizedHandler) WithGroup(name string) slog.Handler {
	c.delegate = c.delegate.WithGroup(name)
	return c
}

// NewColorizedHandler returns a new slog.Handler. If options.ForceColorize is false and the writer is not a TTY,
// this returns the delegate that is created
func NewColorizedHandler[T slog.Handler](targetWriter io.Writer, delegateConstructor func(io.Writer, *slog.HandlerOptions) T, options Options) slog.Handler {

	// if we're not forcing colorization, and we're not writing to a TTY, return a delegate
	if !options.ForceColorize {
		if f, isAFile := targetWriter.(*os.File); isAFile {
			if !isatty.IsTerminal(f.Fd()) || !isatty.IsCygwinTerminal(f.Fd()) {
				return delegateConstructor(targetWriter, options.HandlerOptions)
			}
		}
	}

	stylizer := options.Stylizer
	if stylizer == nil {
		stylizer = DefaultLevelStylizer
	}

	delegateBuffer := &strings.Builder{}
	delegate := delegateConstructor(delegateBuffer, options.HandlerOptions)

	writer := ansi.NewWriter(targetWriter)
	mu := &sync.Mutex{}
	return &ColorizedHandler{
		writer,
		delegateBuffer,
		delegate,
		stylizer,
		mu,
	}
}

package handler_test

import (
	"log/slog"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jghiloni/stylized/handler"
)

var _ = Describe("Handler", func() {
	It("Generates colorized logs", func() {
		out := &strings.Builder{}

		h := handler.NewColorizedHandler(out, slog.NewTextHandler, handler.Options{
			ForceColorize: true,
			HandlerOptions: &slog.HandlerOptions{
				Level: slog.LevelDebug,
				ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						return slog.Any(slog.TimeKey, time.Unix(0, 0).UTC())
					}

					return a
				},
			},
		})

		expectedOutput := "\x1b[1;3;31mtime=1970-01-01T00:00:00.000Z level=ERROR msg=ERROR\x1b[0m\n\x1b[33mtime=1970-01-01T00:00:00.000Z level=WARN msg=WARNING\x1b[0m\ntime=1970-01-01T00:00:00.000Z level=INFO msg=INFO\n\x1b[1;96mtime=1970-01-01T00:00:00.000Z level=DEBUG msg=DEBUG\x1b[0m\n"

		logger := slog.New(h)
		logger.Error("ERROR")
		logger.Warn("WARNING")
		logger.Info("INFO")
		logger.Debug("DEBUG")

		output := out.String()
		Expect(expectedOutput).To(Equal(output))
	})
})

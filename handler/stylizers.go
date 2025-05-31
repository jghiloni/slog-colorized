package handler

import (
	"log/slog"

	"github.com/jghiloni/stylized/ansi"
)

func LevelStylizer(levelStyles map[slog.Level]ansi.Style) RecordStylizer {
	return func(r slog.Record) ansi.Style {
		return levelStyles[r.Level]
	}
}

var defaultStyleMap = map[slog.Level]ansi.Style{
	slog.LevelError: newStyle(ansi.Red, nil, []ansi.StyleModifier{ansi.Bold, ansi.Italic}),
	slog.LevelWarn:  newStyle(ansi.Yellow, nil, nil),
	slog.LevelDebug: newStyle(ansi.Cyan, []ansi.ColorModifier{ansi.Intense}, []ansi.StyleModifier{ansi.Bold}),
}

var DefaultLevelStylizer = LevelStylizer(defaultStyleMap)

func newStyle(c ansi.ANSIColorCode, cm []ansi.ColorModifier, sm []ansi.StyleModifier) ansi.Style {
	p := ansi.Style{}
	p.SetColor(c)
	p.SetColorModifiers(cm...)
	p.SetStyleModifiers(sm...)

	return p
}

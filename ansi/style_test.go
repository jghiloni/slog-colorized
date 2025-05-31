package ansi_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jghiloni/stylized/ansi"
)

var _ = Describe("Style", func() {
	It("Deduplicates and Sorts Color Modifiers", func() {
		ps := ansi.Style{}
		ps.SetColorModifiers(ansi.Intense, ansi.Background, ansi.Background, ansi.Intense)

		Expect(ps.ColorModifiers()).To(HaveExactElements(ansi.Background, ansi.Intense))
	})

	It("Deduplicates and Sorts Style Modifiers", func() {
		ps := ansi.Style{}
		ps.SetStyleModifiers(ansi.Bold, ansi.Strikethrough, ansi.Bold, ansi.Italic)
		Expect(ps.StyleModifiers()).To(HaveExactElements(ansi.Bold, ansi.Italic, ansi.Strikethrough))
	})

	It("Renders the code correctly", func() {
		ps := ansi.Style{}
		ps.SetColor(ansi.Cyan)
		ps.SetColorModifiers(ansi.Intense, ansi.Background)
		ps.SetStyleModifiers(ansi.Bold)

		Expect(ps.String()).To(Equal("\u001b[1;106m"))
	})

	It("Renders the reset code correctly", func() {
		ps := ansi.Style{}
		Expect(ps.String()).To(Equal("\u001b[0m"))
	})
})

package internal

import (
	"strings"
	"testing"
)

func TestMarkdownGenerator(t *testing.T) {
	generator := newMarkdownGenerator()
	generator.WriteVersion("v9.0.0")
	generator.WriteLine()
	generator.WriteSection("BrighterConfiguration")
	generator.WriteLine()
	generator.WriteTOCs([]TOCEntry{
		{
			Name:   "DocumentOne",
			File:   "DocumentOne.md",
			Indent: 1,
		},
		{
			Name:   "DocumentTwo",
			File:   "DocumentTwo.md",
			Indent: 1,
		},
	})

	generator.WriteLine()
	generator.WriteSection("DarkerConfiguration")
	generator.WriteLine()
	generator.WriteTOCs([]TOCEntry{
		{
			Name:   "DocumentThree",
			File:   "DocumentThree.md",
			Indent: 1,
		},
	})

	markdown := generator.buffer.String()
	expected := "## v9.0.0\n\n### BrighterConfiguration\n\n * [DocumentOne](DocumentOne.md)\n * [DocumentTwo](DocumentTwo.md)\n\n### DarkerConfiguration\n\n * [DocumentThree](DocumentThree.md)\n"
	if !strings.EqualFold(markdown, expected) {
		t.Errorf("Markdown does not match, expected %s got %s", expected, markdown)
	}

}

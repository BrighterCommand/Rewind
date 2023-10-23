package book

import (
	"github.com/brightercommand/Rewind/internal/pages"
	"strings"
	"testing"
)

func TestMarkdownGenerator(t *testing.T) {
	generator := newMarkdownGenerator()
	generator.WriteVersion("9")
	generator.WriteLine()
	generator.WriteSection("BrighterConfiguration")
	generator.WriteLine()
	generator.WriteTOCs([]pages.TOCEntry{
		{
			Name:  "DocumentOne",
			File:  "DocumentOne.md",
			Order: 1,
		},
		{
			Name:  "DocumentTwo",
			File:  "DocumentTwo.md",
			Order: 1,
		},
	}, "9")

	generator.WriteLine()
	generator.WriteSection("DarkerConfiguration")
	generator.WriteLine()
	generator.WriteTOCs([]pages.TOCEntry{
		{
			Name:  "DocumentThree",
			File:  "DocumentThree.md",
			Order: 1,
		},
	}, "9")

	markdown := generator.buffer.String()
	expected := "## 9\n\n### BrighterConfiguration\n\n  * [DocumentOne](/contents/9/DocumentOne.md)\n  * [DocumentTwo](/contents/9/DocumentTwo.md)\n\n### DarkerConfiguration\n\n  * [DocumentThree](/contents/9/DocumentThree.md)\n"
	if !strings.EqualFold(markdown, expected) {
		t.Errorf("Markdown does not match, expected %s got %s", expected, markdown)
	}
}

package book

/*
func TestMarkdownGenerator(t *testing.T) {
	generator := newMarkdownGenerator()
	generator.WriteVersion("v9.0.0")
	generator.WriteLine()
	generator.WriteSection("BrighterConfiguration")
	generator.WriteLine()
	generator.WriteTOCs([]pages.TOCEntry{
		{
			Name:   "DocumentOne",
			File:   "DocumentOne.md",
			Order: 1,
		},
		{
			Name:   "DocumentTwo",
			File:   "DocumentTwo.md",
			Order: 1,
		},
	}, "9.0.0")

	generator.WriteLine()
	generator.WriteSection("DarkerConfiguration")
	generator.WriteLine()
	generator.WriteTOCs([]pages.TOCEntry{
		{
			Name:   "DocumentThree",
			File:   "DocumentThree.md",
			Order: 1,
		},
	}, "10.0.0")

	markdown := generator.buffer.String()
	expected := "## v9.0.0\n\n### BrighterConfiguration\n\n  * [DocumentOne](/contents/9.0.0/DocumentOne.md)\n  * [DocumentTwo](/contents/9.0.0/DocumentTwo.md)\n\n### DarkerConfiguration\n\n  * [DocumentThree](/contents/10.0.0/DocumentThree.md)\n"
	if !strings.EqualFold(markdown, expected) {
		t.Errorf("Markdown does not match, expected %s got %s", expected, markdown)
	}
}

*/

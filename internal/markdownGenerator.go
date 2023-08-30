package internal

import (
	"fmt"
	"os"
	"strings"
)

type markdownGenerator struct {
	destPath string
	buffer   *strings.Builder
}

func newMarkdownGenerator() *markdownGenerator {
	return &markdownGenerator{
		buffer: &strings.Builder{},
	}
}

func (g *markdownGenerator) GenerateSummary(entries *versionedToc, summary *os.File) error {

	for version, toc := range *entries {
		g.WriteVersion(version)
		g.WriteLine()
		for section, entries := range toc {
			g.WriteSection(section)
			g.WriteLine()
			g.WriteTOCs(entries, version)
			g.WriteLine()
		}
	}

	_, err := summary.WriteString(g.buffer.String())
	if err != nil {
		return err
	}

	return nil
}

func (g *markdownGenerator) WriteLine() {
	g.buffer.WriteString("\n")
}

func (g *markdownGenerator) WriteVersion(version string) {
	g.buffer.WriteString(g.getTitle(version, 2))
	g.buffer.WriteString("\n")
}

func (g *markdownGenerator) WriteSection(section string) {
	g.buffer.WriteString(g.getTitle(section, 3))
	g.buffer.WriteString("\n")
}

func (g *markdownGenerator) WriteTOCs(entries []TOCEntry, version string) {
	for _, entry := range entries {
		g.buffer.WriteString(g.getListItemWithIndent(entry.Name, g.getLinkPath(entry, version), entry.Indent))
		g.buffer.WriteString("\n")
	}
}

func (g *markdownGenerator) getLinkPath(entry TOCEntry, version string) string {
	link := "/" + ContentDirName + "/" + version + "/" + entry.File
	//spaces in markdown links may cause issues, so replace them with %20
	nospace := strings.ReplaceAll(link, " ", "%20")
	return nospace
}

func (m *markdownGenerator) getLink(desc, url string) string {
	return fmt.Sprintf("[%s](%s)", desc, url)
}

func (m *markdownGenerator) getListItemWithIndent(desc, url string, indent int) string {
	return strings.Repeat("  ", indent) + "*" + " " + m.getLink(desc, url)
}

func (m *markdownGenerator) getTitle(content string, level int) string {
	return strings.Repeat("#", level) + " " + content
}

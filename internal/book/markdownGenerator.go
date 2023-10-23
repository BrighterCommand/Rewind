package book

import (
	"fmt"
	"github.com/brightercommand/Rewind/internal/pages"
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

func (g *markdownGenerator) GenerateSummary(entries []pages.OrderedVersionTocs, summary *os.File) error {

	for _, toc := range entries {
		g.WriteVersion(toc.Version)
		g.WriteLine()
		for _, section := range toc.Sections {
			g.WriteSection(section.Name)
			g.WriteLine()
			g.WriteTOCs(section.Section.Entries, toc.Version)
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

func (g *markdownGenerator) WriteTOCs(s []pages.TOCEntry, version string) {
	for _, entry := range s {
		g.buffer.WriteString(g.getListItemWithIndent(entry.Name, g.getLinkPath(entry, version), entry.Order))
		g.buffer.WriteString("\n")
	}
}

func (g *markdownGenerator) getLinkPath(entry pages.TOCEntry, version string) string {
	link := "/" + pages.ContentDirName + "/" + version + "/" + entry.File
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

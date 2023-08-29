package internal

import (
	"fmt"
	"strings"
)

type markdownGenerator struct {
	buffer *strings.Builder
}

func newMarkdownGenerator() *markdownGenerator {
	return &markdownGenerator{
		buffer: &strings.Builder{},
	}
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

func (g *markdownGenerator) WriteTOCs(entries []TOCEntry) {
	for _, entry := range entries {
		g.buffer.WriteString(g.getListItemWithIndent(entry.Name, entry.File, entry.Indent))
		g.buffer.WriteString("\n")
	}
}

func (m *markdownGenerator) getLink(desc, url string) string {
	return fmt.Sprintf("[%s](%s)", desc, url)
}

func (m *markdownGenerator) getListItemWithIndent(desc, url string, indent int) string {
	return strings.Repeat(" ", indent) + "*" + " " + m.getLink(desc, url)
}

func (m *markdownGenerator) getTitle(content string, level int) string {
	return strings.Repeat("#", level) + " " + content
}

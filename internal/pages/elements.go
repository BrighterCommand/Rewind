package pages

import "os"

const SummaryFileName = "SUMMARY.md"
const ContentDirName = "contents"

type Shared struct {
	Docs map[string]Doc
	TOC  *Doc
}

type Version struct {
	DestPath string
	Docs     map[string]Doc
	TOC      *Doc
	Version  string
}

type Doc struct {
	SourcePath string
	Version    string
	Storage    os.DirEntry
}

type Root struct {
	DestPath   string
	SourcePath string
	GitBook    *Doc
	WorkDir    string
}

type TOCEntry struct {
	Name   string
	File   string
	Indent int
}

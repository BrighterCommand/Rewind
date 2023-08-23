package internal

import "os"

type Shared struct {
	Docs map[string]Doc
	TOC  Doc
}

type Version struct {
	DestPath string
	Docs     map[string]Doc
	TOC      Doc
	Version  string
}

type Doc struct {
	SourcePath string
	Version    string
	Storage    os.DirEntry
}

type Root struct {
	DestPath string
	Summary  Doc
	GitBook  Doc
}

type TOCEntry struct {
	Name   string
	File   string
	Indent int
}

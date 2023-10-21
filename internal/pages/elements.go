package pages

import "os"

const SummaryFileName = "SUMMARY.md"
const ContentDirName = "contents"
const StaticFolderName = "_static"
const ImageFolderName = "images"

// AssetType Where we have an asset that is not a markdown file what is its type
type AssetType int

// An enumerated set of asset types
const (
	Undefined AssetType = iota
	Image
)

// Asset A binary asset used with the markdown files, such as images.
type Asset struct {
	SourcePath string
	What       AssetType
	Version    string
	Storage    os.DirEntry
}

// Doc A markdown document.
type Doc struct {
	SourcePath string
	Version    string
	Storage    os.DirEntry
}

// Root The root of the book.
type Root struct {
	DestPath   string
	SourcePath string
	GitBook    *Doc
	ReadMe     *Doc
	WorkDir    string
}

// Shared Assets & Docs shared by all versions of the book
type Shared struct {
	Docs   map[string]Doc
	Images map[string]Asset
	TOC    *Doc
}

// TOCSection A table of contents section - the name of the section is held in the Toc map below.
type TOCSection struct {
	Order   int        `yaml:"order"`
	Entries []TOCEntry `yaml:"entries"`
}

// TOCEntry A table of contents entry.
type TOCEntry struct {
	Name  string `yaml:"name"`
	File  string `yaml:"file"`
	Order int    `yaml:"order"`
}

// Version Docs & Assets for a version of the book
type Version struct {
	DestPath string
	Docs     map[string]Doc
	Images   map[string]Asset
	TOC      *Doc
	Version  string
}

type OrderedTocSection struct {
	Name    string
	Section TOCSection
}

type Toc map[string]TOCSection

type VersionedToc map[string]Toc

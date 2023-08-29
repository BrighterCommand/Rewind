package internal

import (
	"os"
	"strings"
)

// fakeDirEntry is a fake implementation of os.DirEntry
// It is used for testing.
type fakeDirEntry struct {
	name  string
	isDir bool
}

func (f fakeDirEntry) Name() string {
	return f.name
}

func (f fakeDirEntry) IsDir() bool {
	return f.isDir
}

func (f fakeDirEntry) Type() os.FileMode {
	return os.ModeDir
}

func (f fakeDirEntry) Info() (os.FileInfo, error) {
	return nil, nil
}

func sourceTestDataBuilder(sourcePath string, mydir string) *Sources {
	sources := &Sources{
		Root: &Root{
			GitBook: &Doc{
				SourcePath: sourcePath,
				Storage: fakeDirEntry{
					name:  ".gitbook.yaml",
					isDir: false,
				},
			},
		},
		Shared: &Shared{
			Docs: make(map[string]Doc),
		},
		Versions: make(map[string]Version),
	}

	//add shared docs
	sources.Shared.TOC = &Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  ".toc.yaml",
			isDir: false,
		}}

	sources.Shared.Docs["DocumentOne.md"] = Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "DocumentOne.md",
			isDir: false,
		}}

	sources.Shared.Docs["DocumentTwo.md"] = Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "DocumentTwo.md",
			isDir: false,
		}}

	sources.Shared.Docs["DocumentThree.md"] = Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "DocumentThree.md",
			isDir: false,
		}}

	//add versions

	//add versionNine docs
	versionNine := Version{
		Docs:    make(map[string]Doc),
		Version: "9.0.0",
	}

	versionNine.TOC = &Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/9.0.0", 1),
		Version:    "9.0.0",
		Storage: fakeDirEntry{
			name:  ".toc.yaml",
			isDir: false,
		}}

	versionNine.Docs["DocumentTwo.md"] = Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/9.0.0", 1),
		Version:    "9.0.0",
		Storage: fakeDirEntry{
			name:  "DocumentTwo.md",
			isDir: false,
		}}

	sources.Versions["9.0.0"] = versionNine

	versionTen := Version{
		Docs:    make(map[string]Doc),
		Version: "10.0.0",
	}

	versionTen.TOC = &Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/10.0.0", 1),
		Version:    "10.0.0",
		Storage: fakeDirEntry{
			name:  ".toc.yaml",
			isDir: false,
		}}

	versionTen.Docs["DocumentOne.md"] = Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/10.0.0", 1),
		Version:    "10.0.0",
		Storage: fakeDirEntry{
			name:  "DocumentOne.md",
			isDir: false,
		}}

	versionTen.Docs["DocumentFour.md"] = Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/10.0.0", 1),
		Version:    "10.0.0",
		Storage: fakeDirEntry{
			name:  "DocumentFour.md",
			isDir: false,
		}}

	sources.Versions["10.0.0"] = versionTen

	return sources
}

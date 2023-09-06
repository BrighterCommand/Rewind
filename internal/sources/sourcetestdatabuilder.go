package sources

import (
	"github.com/brightercommand/Rewind/internal/pages"
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

func SourceTestDataBuilder(sourcePath string, mydir string) *Sources {
	sources := &Sources{
		Root: &pages.Root{
			GitBook: &pages.Doc{
				SourcePath: sourcePath,
				Storage: fakeDirEntry{
					name:  "..gitbook.yaml",
					isDir: false,
				},
			},
			ReadMe: &pages.Doc{
				SourcePath: sourcePath,
				Storage: fakeDirEntry{
					name:  "README.md",
					isDir: false,
				},
			},
		},
		Shared: &pages.Shared{
			Docs:   make(map[string]pages.Doc),
			Images: make(map[string]pages.Asset),
		},
		Versions: make(map[string]pages.Version),
	}

	//add shared docs
	sources.Shared.TOC = &pages.Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  ".toc.yaml",
			isDir: false,
		}}

	sources.Shared.Docs["DocumentOne.md"] = pages.Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "DocumentOne.md",
			isDir: false,
		}}

	sources.Shared.Docs["DocumentTwo.md"] = pages.Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "DocumentTwo.md",
			isDir: false,
		}}

	sources.Shared.Docs["DocumentThree.md"] = pages.Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "DocumentThree.md",
			isDir: false,
		}}

	//add shared images
	sources.Shared.Images["ImageOne.png"] = pages.Asset{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared/_static/images", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "ImageOne.png",
			isDir: false,
		},
	}

	sources.Shared.Images["ImageTwo.png"] = pages.Asset{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared/_static/images", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "ImageTwo.png",
			isDir: false,
		},
	}

	//add versions

	//add versionNine docs
	versionNine := pages.Version{
		Docs:    make(map[string]pages.Doc),
		Version: "9.0.0",
	}

	versionNine.TOC = &pages.Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/9.0.0", 1),
		Version:    "9.0.0",
		Storage: fakeDirEntry{
			name:  ".toc.yaml",
			isDir: false,
		}}

	versionNine.Docs["DocumentTwo.md"] = pages.Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/9.0.0", 1),
		Version:    "9.0.0",
		Storage: fakeDirEntry{
			name:  "DocumentTwo.md",
			isDir: false,
		}}

	//add shared images
	versionNine.Images["ImageThree.png"] = pages.Asset{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared/_static/images", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "ImageThree.png",
			isDir: false,
		},
	}

	versionNine.Images["ImageFour.png"] = pages.Asset{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared/_static/images", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "ImageFour.png",
			isDir: false,
		},
	}

	sources.Versions["9.0.0"] = versionNine

	versionTen := pages.Version{
		Docs:    make(map[string]pages.Doc),
		Version: "10.0.0",
	}

	versionTen.TOC = &pages.Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/10.0.0", 1),
		Version:    "10.0.0",
		Storage: fakeDirEntry{
			name:  ".toc.yaml",
			isDir: false,
		}}

	versionTen.Docs["DocumentOne.md"] = pages.Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/10.0.0", 1),
		Version:    "10.0.0",
		Storage: fakeDirEntry{
			name:  "DocumentOne.md",
			isDir: false,
		}}

	versionTen.Docs["DocumentFour.md"] = pages.Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/10.0.0", 1),
		Version:    "10.0.0",
		Storage: fakeDirEntry{
			name:  "DocumentFour.md",
			isDir: false,
		}}

	//add shared images
	versionTen.Images["ImageOne.png"] = pages.Asset{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared/_static/images", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "ImageThree.png",
			isDir: false,
		},
	}

	versionTen.Images["ImageThree.png"] = pages.Asset{
		SourcePath: strings.Replace(mydir, "internal", "test/source/shared/_static/images", 1),
		Version:    "shared",
		Storage: fakeDirEntry{
			name:  "ImageFour.png",
			isDir: false,
		},
	}

	sources.Versions["10.0.0"] = versionTen

	return sources
}

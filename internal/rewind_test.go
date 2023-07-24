package internal

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
	"testing"
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

func TestFindSources(t *testing.T) {

	mydir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %s", err)
	}

	sourcePath := strings.Replace(mydir, "internal", "test/source", 1)

	sources, err := FindSources(sourcePath)
	if err != nil {
		t.Errorf("Error finding sources: %s", err)
	}

	root := sources.Root
	if root.Summary == nil || root.Summary.Name() != "SUMMARY.md" {
		t.Errorf("Expected SUMMARY.md")
	}

	if root.GitBook == nil || root.GitBook.Name() != ".gitbook.yaml" {
		t.Errorf("Expected .gitbook.yaml")
	}

	docs := sources.Shared.Assets.Docs

	if len(docs) != 3 {
		t.Errorf("Expected 3 documents, got %d", len(docs))
	} else {
		//note order here as sort is alphabetical
		if docs[0].Name() != "DocumentOne.md" {
			t.Errorf("Expected DocumentOne.md, got %s", docs[0])
		}

		if docs[1].Name() != "DocumentThree.md" {
			t.Errorf("Expected DocumentThree.md, got %s", docs[1])
		}

		if docs[2].Name() != "DocumentTwo.md" {
			t.Errorf("Expected DocumentTwo.md, got %s", docs[1])
		}
	}

	versions := sources.Versions
	if len(versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(versions))
	} else {

		if versions[0].Version != "10.0.0" {
			t.Errorf("Expected 9.0.0, got %s", versions[0].Version)
		}
		if versions[1].Version != "9.0.0" {
			t.Errorf("Expected 9.0.0, got %s", versions[1].Version)
		}
	}
}

func TestFileCopy(t *testing.T) {
	mydir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %s", err)
	}

	sourcePath := strings.Replace(mydir, "internal", "test/filecopy", 1)
	destPath := strings.Replace(mydir, "internal", fmt.Sprintf("test/filecopy/%s", uuid.New().String()), 1)

	err = copyFile(sourcePath, destPath, "DocumentOne.md")
	if err != nil {
		t.Errorf("Error copying file: %s", err)
	}

	_, err = os.Stat(destPath)
	if err != nil {
		t.Errorf("Error copying file: %s", err)
	}

	// Remove the directory
	err = os.RemoveAll(destPath)
	if err != nil {
		t.Errorf("Error removing directory: %s", err)
	}

}

// TestBookBuilder tests the book builder.
// It creates a fake directory structure and then runs the book builder.
func TestBookBuilder(t *testing.T) {

	mydir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %s", err)
	}

	sourcePath := strings.Replace(mydir, "internal", "test/filecopy", 1)
	destPath := strings.Replace(mydir, "internal", fmt.Sprintf("test/filecopy/%s", uuid.New().String()), 1)

	sources := sourceTestDataBuilder(sourcePath, mydir)

	book := sources.BuildBook(destPath)

	if book == nil {
		t.Errorf("Error building book: %s", err)
	}

}

func sourceTestDataBuilder(sourcePath string, mydir string) *Sources {
	sources := &Sources{
		Root: &Root{
			Path: sourcePath,
			Summary: fakeDirEntry{
				name:  "SUMMARY.md",
				isDir: false,
			},
			GitBook: fakeDirEntry{
				name:  ".gitbook.yaml",
				isDir: false,
			},
		},
		Shared: &Shared{
			Assets: &Assets{
				Path: strings.Replace(mydir, "internal", "test/source/shared", 1),
				Docs: []os.DirEntry{
					fakeDirEntry{
						name:  "DocumentOne.md",
						isDir: false,
					},
					fakeDirEntry{
						name:  "DocumentTwo.md",
						isDir: false,
					},
					fakeDirEntry{
						name:  "DocumentThree.md",
						isDir: false,
					},
				},
			},
		},
		Versions: []Version{
			Version{
				Assets: &Assets{
					Path: strings.Replace(mydir, "internal", "test/source/9.0.0", 1),
					Docs: []os.DirEntry{
						fakeDirEntry{
							name:  "DocumentTwo.md",
							isDir: false,
						},
					},
				},
				Version: "9.0.0",
			},
			Version{
				Assets: &Assets{
					Path: strings.Replace(mydir, "internal", "test/source/10.0.0", 1),
					Docs: []os.DirEntry{
						fakeDirEntry{
							name:  "DocumentOne.md",
							isDir: false,
						},
					},
				},
				Version: "10.0.0",
			},
		},
	}

	return sources
}

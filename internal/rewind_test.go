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
		if docs["DocumentOne.md"].Storage.Name() != "DocumentOne.md" {
			t.Errorf("Expected DocumentOne.md, got %s", docs["DocumentOne.md"])
		}

		if docs["DocumentTwo.md"].Storage.Name() != "DocumentTwo.md" {
			t.Errorf("Expected DocumentTwo.md, got %s", docs["DocumentTwo.md"])
		}

		if docs["DocumentThree.md"].Storage.Name() != "DocumentThree.md" {
			t.Errorf("Expected DocumentThree.md, got %s", docs["DocumentThree.md"])
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

	sourcePath := strings.Replace(mydir, "internal", "test/source", 1)
	destPath := strings.Replace(mydir, "internal", fmt.Sprintf("test/docs/%s", uuid.New().String()), 1)

	sources := sourceTestDataBuilder(sourcePath, mydir)

	book := sources.BuildBook(destPath)

	if book == nil {
		t.Errorf("Error building book: %s", err)
	}

	if book.Root.SourcePath != sourcePath {
		t.Errorf("Expected %s but found %s", sourcePath, book.Root.SourcePath)
	}
	if book.Root.DestPath != destPath {
		t.Errorf("Expected %s but found %s", destPath, book.Root.DestPath)
	}
	if book.Root.Summary.Name() != sources.Root.Summary.Name() {
		t.Errorf("Expected %s but found %s", sources.Root.Summary.Name(), book.Root.Summary.Name())
	}
	if book.Root.GitBook.Name() != sources.Root.GitBook.Name() {
		t.Errorf("Expected %s but found %s", sources.Root.GitBook.Name(), book.Root.GitBook.Name())
	}

	if len(book.Versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(book.Versions))
	}

	firstVersion := book.Versions[0]
	if firstVersion.Version != "10.0.0" {
		t.Errorf("Expected 10.0.0, got %s", book.Versions[0].Version)
	}

	if firstVersion.Assets.SourcePath != sources.Versions[0].Assets.SourcePath {
		t.Errorf("Expected %s, got %s", sources.Versions[0].Assets.SourcePath, firstVersion.Assets.SourcePath)
	}

	if firstVersion.Assets.DestPath != fmt.Sprintf("%s/10.0.0", destPath) {
		t.Errorf("Expected %s, got %s", fmt.Sprintf("%s/10.0.0", destPath), firstVersion.Assets.DestPath)
	}

	//NOTE: this won't tell us which version of a document in shared and version we got, which we need to consider
	if len(firstVersion.Assets.Docs) != 3 {
		t.Errorf("Expected 3 docs, got %d", len(firstVersion.Assets.Docs))
	}

	secondVersion := book.Versions[1]
	if secondVersion.Version != "9.0.0" {
		t.Errorf("Expected 9.0.0, got %s", book.Versions[1].Version)
	}

	if secondVersion.Assets.SourcePath != sources.Versions[1].Assets.SourcePath {
		t.Errorf("Expected %s, got %s", sources.Versions[1].Assets.SourcePath, secondVersion.Assets.SourcePath)
	}

	if secondVersion.Assets.DestPath != fmt.Sprintf("%s/9.0.0", destPath) {
		t.Errorf("Expected %s, got %s", fmt.Sprintf("%s/9.0.0", destPath), secondVersion.Assets.DestPath)
	}

	//NOTE: this won't tell us which version of a document in shared and version we got, which we need to consider
	if len(secondVersion.Assets.Docs) != 3 {
		t.Errorf("Expected 3 docs, got %d", len(secondVersion.Assets.Docs))
	}
}

func sourceTestDataBuilder(sourcePath string, mydir string) *Sources {
	sources := &Sources{
		Root: &Root{
			SourcePath: sourcePath,
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
				SourcePath: strings.Replace(mydir, "internal", "test/source/shared", 1),
				Docs:       make(map[string]Doc),
			},
		},
		Versions: []Version{
			Version{
				Assets: &Assets{
					SourcePath: strings.Replace(mydir, "internal", "test/source/9.0.0", 1),
					Docs:       make(map[string]Doc),
				},
				Version: "9.0.0",
			},
			Version{
				Assets: &Assets{
					SourcePath: strings.Replace(mydir, "internal", "test/source/10.0.0", 1),
					Docs:       make(map[string]Doc),
				},
				Version: "10.0.0",
			},
		},
	}

	//add shared docs
	sources.Shared.Assets.Docs["DocumentOne.md"] = Doc{
		Version: "Shared", Storage: fakeDirEntry{
			name:  "DocumentOne.md",
			isDir: false,
		}}

	sources.Shared.Assets.Docs["DocumentTwo.md"] = Doc{
		Version: "Shared", Storage: fakeDirEntry{
			name:  "DocumentTwo.md",
			isDir: false,
		}}

	sources.Shared.Assets.Docs["DocumentThree.md"] = Doc{
		Version: "Shared", Storage: fakeDirEntry{
			name:  "DocumentThree.md",
			isDir: false,
		}}

	//add version docs
	sources.Versions[0].Assets.Docs["DocumentTwo.md"] = Doc{
		Version: "9.0.0", Storage: fakeDirEntry{
			name:  "DocumentTwo.md",
			isDir: false,
		}}

	sources.Versions[1].Assets.Docs["DocumentThree.md"] = Doc{
		Version: "10.0.0", Storage: fakeDirEntry{
			name:  "DocumentThree.md",
			isDir: false,
		}}

	return sources
}

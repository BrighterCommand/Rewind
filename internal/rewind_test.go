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
	if root.Summary.Storage.Name() != "SUMMARY.md" {
		t.Errorf("Expected SUMMARY.md")
	}

	if root.GitBook.Storage.Name() != ".gitbook.yaml" {
		t.Errorf("Expected .gitbook.yaml")
	}

	docs := sources.Shared.Docs

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

		if versions["10.0.0"].Version != "10.0.0" {
			t.Errorf("Expected 9.0.0, got %s", versions["10.0.0"].Version)
		}
		if versions["9.0.0"].Version != "9.0.0" {
			t.Errorf("Expected 9.0.0, got %s", versions["9.0.0"].Version)
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

	if book.Root.DestPath != destPath {
		t.Errorf("Expected %s but found %s", destPath, book.Root.DestPath)
	}

	if book.Root.Summary.Storage.Name() != sources.Root.Summary.Storage.Name() {
		t.Errorf("Expected %s but found %s", sources.Root.Summary.Storage.Name(), book.Root.Summary.Storage.Name())
	}

	if book.Root.Summary.SourcePath != sourcePath {
		t.Errorf("Expected %s but found %s", sourcePath, book.Root.Summary.SourcePath)
	}

	if book.Root.GitBook.Storage.Name() != sources.Root.GitBook.Storage.Name() {
		t.Errorf("Expected %s but found %s", sources.Root.GitBook.Storage.Name(), book.Root.GitBook.Storage.Name())
	}

	if book.Root.GitBook.SourcePath != sourcePath {
		t.Errorf("Expected %s but found %s", sourcePath, book.Root.GitBook.SourcePath)
	}

	if len(book.Versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(book.Versions))
	}

	versionNine := book.Versions["9.0.0"]
	if versionNine.Version != "9.0.0" {
		t.Errorf("Expected 9.0.0, got %s", book.Versions["9.0.0"].Version)
	}

	if versionNine.DestPath != fmt.Sprintf("%s/9.0.0", destPath) {
		t.Errorf("Expected %s, got %s", fmt.Sprintf("%s/9.0.0", destPath), versionNine.DestPath)
	}

	if len(versionNine.Docs) != 3 {
		t.Errorf("Expected 3 docs, got %d", len(versionNine.Docs))
	}

	var expectedSourcePath string = fmt.Sprintf("%s/9.0.0", sourcePath)
	sharedSourcePath := fmt.Sprintf("%s/Shared", sourcePath)
	var documentOneFound, documentTwoFound, documentThreeFound bool
	for _, doc := range versionNine.Docs {

		if doc.Storage.Name() == "DocumentOne.md" {
			if !strings.EqualFold(doc.SourcePath, sharedSourcePath) {
				t.Errorf("Expected %s, got %s", sharedSourcePath, doc.SourcePath)
			} else {
				documentOneFound = true
			}
		}

		if doc.Storage.Name() == "DocumentTwo.md" {
			if !strings.EqualFold(doc.SourcePath, expectedSourcePath) {
				t.Errorf("Expected %s, got %s", expectedSourcePath, doc.SourcePath)
			} else {
				documentTwoFound = true
			}
		}
		if doc.Storage.Name() == "DocumentThree.md" {
			if !strings.EqualFold(doc.SourcePath, sharedSourcePath) {
				t.Errorf("Expected %s, got %s", sharedSourcePath, doc.SourcePath)
			} else {
				documentThreeFound = true
			}
		}
	}

	if !(documentOneFound && documentTwoFound && documentThreeFound) {
		t.Errorf("Expected DocumentOne.md, DocumentTwo.md, DocumentThree.md in v9.0.0")
	}

	versionTen := book.Versions["10.0.0"]
	if versionTen.Version != "10.0.0" {
		t.Errorf("Expected 10.0.0, got %s", book.Versions["10.0.0"].Version)
	}

	if versionTen.DestPath != fmt.Sprintf("%s/10.0.0", destPath) {
		t.Errorf("Expected %s, got %s", fmt.Sprintf("%s/10.0.0", destPath), versionTen.DestPath)
	}

	if len(versionTen.Docs) != 3 {
		t.Errorf("Expected 3 docs, got %d", len(versionTen.Docs))
	}

	expectedSourcePath = fmt.Sprintf("%s/10.0.0", sourcePath)
	sharedSourcePath = fmt.Sprintf("%s/Shared", sourcePath)
	for _, doc := range versionTen.Docs {

		if doc.Storage.Name() == "DocumentOne.md" {
			if !strings.EqualFold(doc.SourcePath, expectedSourcePath) {
				t.Errorf("Expected %s, got %s", expectedSourcePath, doc.SourcePath)
			} else {
				documentOneFound = true
			}
		}

		if doc.Storage.Name() == "DocumentTwo.md" {
			if !strings.EqualFold(doc.SourcePath, sharedSourcePath) {
				t.Errorf("Expected %s, got %s", sharedSourcePath, doc.SourcePath)
			} else {
				documentTwoFound = true
			}
		}
		if doc.Storage.Name() == "DocumentThree.md" {
			if !strings.EqualFold(doc.SourcePath, sharedSourcePath) {
				t.Errorf("Expected %s, got %s", sharedSourcePath, doc.SourcePath)
			} else {
				documentThreeFound = true
			}
		}
	}

	if !(documentOneFound && documentTwoFound && documentThreeFound) {
		t.Errorf("Expected DocumentOne.md, DocumentTwo.md, DocumentThree.md in v10.0.0")
	}
}

func TestBookCreation(t *testing.T) {

	mydir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %s", err)
	}

	sourcePath := strings.Replace(mydir, "internal", "test/source", 1)
	destPath := strings.Replace(mydir, "internal", fmt.Sprintf("test/docs/%s", uuid.New().String()), 1)

	sources, err := FindSources(sourcePath)
	if err != nil {
		t.Errorf("Error finding sources: %s", err)
	}

	book := sources.BuildBook(destPath)

	err = book.Create()
	if err != nil {
		t.Errorf("Error creating book: %s", err)
	}

	entries, err := os.ReadDir(destPath)
	if err != nil {
		t.Errorf("Error reading directory: %s", err)
	}

	var summaryFound, gitbookFound, v9Found, v10Found bool
	for _, entry := range entries {
		if entry.Name() == "SUMMARY.md" {
			summaryFound = true
		} else if entry.Name() == ".gitbook.yaml" {
			gitbookFound = true
		} else if entry.IsDir() {
			if entry.Name() == "9.0.0" {
				v9entries, err := os.ReadDir(fmt.Sprintf("%s/9.0.0", destPath))
				if err != nil {
					t.Errorf("Error reading directory: %s", err)
					v9Found = false
				} else {
					v9Found = findFiles(v9entries)
				}
			} else if entry.Name() == "10.0.0" {
				v10entries, err := os.ReadDir(fmt.Sprintf("%s/10.0.0", destPath))
				if err != nil {
					t.Errorf("Error reading directory: %s", err)
					v10Found = false
				} else {
					v10Found = findFiles(v10entries)
				}
			}
		}
	}

	if summaryFound == false {
		t.Errorf("Expected SUMMARY.md")
	}

	if gitbookFound == false {
		t.Errorf("Expected .gitbook.yaml")
	}

	if v9Found == false {
		t.Errorf("Expected 9.0.0")
	}

	if v10Found == false {
		t.Errorf("Expected 10.0.0")
	}

	// Remove the directory
	err = os.RemoveAll(destPath)
	if err != nil {
		t.Errorf("Error removing directory: %s", err)
	}
}

func findFiles(entries []os.DirEntry) bool {
	var documentOneFound, documentTwoFound, documentThreeFound bool
	for _, doc := range entries {
		if doc.Name() == "DocumentOne.md" {
			documentOneFound = true
		}
		if doc.Name() == "DocumentTwo.md" {
			documentTwoFound = true
		}
		if doc.Name() == "DocumentThree.md" {
			documentThreeFound = true
		}
	}
	return documentOneFound && documentTwoFound && documentThreeFound
}

func sourceTestDataBuilder(sourcePath string, mydir string) *Sources {
	sources := &Sources{
		Root: &Root{
			Summary: Doc{
				SourcePath: sourcePath,
				Storage: fakeDirEntry{
					name:  "SUMMARY.md",
					isDir: false,
				},
			},
			GitBook: Doc{
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
	sources.Versions["9.0.0"] = Version{
		Docs:    make(map[string]Doc),
		Version: "9.0.0",
	}

	sources.Versions["10.0.0"] = Version{
		Docs:    make(map[string]Doc),
		Version: "10.0.0",
	}

	//add version docs
	sources.Versions["9.0.0"].Docs["DocumentTwo.md"] = Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/9.0.0", 1),
		Version:    "9.0.0",
		Storage: fakeDirEntry{
			name:  "DocumentTwo.md",
			isDir: false,
		}}

	sources.Versions["10.0.0"].Docs["DocumentOne.md"] = Doc{
		SourcePath: strings.Replace(mydir, "internal", "test/source/10.0.0", 1),
		Version:    "10.0.0",
		Storage: fakeDirEntry{
			name:  "DocumentOne.md",
			isDir: false,
		}}

	return sources
}

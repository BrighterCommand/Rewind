package internal

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/md"
	"github.com/gomarkdown/markdown/parser"
	"github.com/google/uuid"
	"os"
	"strings"
	"testing"
)

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

	var sources = sourceTestDataBuilder(sourcePath, mydir)
	var book = sources.BuildBook(destPath)
	if book == nil {
		t.Errorf("Error building book: %s", err)
	}

	checkBookRoot(t, book, destPath, sources, sourcePath)

	if len(book.Versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(book.Versions))
	}

	checkVersion9(t, book, destPath, sourcePath)

	checkVersion10(t, book, destPath, sourcePath)
}

func checkBookRoot(t *testing.T, book *Book, destPath string, sources *Sources, sourcePath string) {
	if book.Root.DestPath != destPath {
		t.Errorf("Expected %s but found %s", destPath, book.Root.DestPath)
	}

	if book.Root.GitBook.Storage == nil {
		t.Errorf("Expected GitBook Document")
		return
	}

	if book.Root.GitBook.Storage.Name() != sources.Root.GitBook.Storage.Name() {
		t.Errorf("Expected %s but found %s", sources.Root.GitBook.Storage.Name(), book.Root.GitBook.Storage.Name())
	}

	if book.Root.GitBook.SourcePath != sourcePath {
		t.Errorf("Expected %s but found %s", sourcePath, book.Root.GitBook.SourcePath)
	}
}

func checkVersion9(t *testing.T, book *Book, destPath string, sourcePath string) {
	var documentOneFound, documentTwoFound, documentThreeFound bool

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

	var expectedSourcePath = fmt.Sprintf("%s/9.0.0", sourcePath)
	sharedSourcePath := fmt.Sprintf("%s/Shared", sourcePath)
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
	return
}

func checkVersion10(t *testing.T, book *Book, destPath string, sourcePath string) {
	var documentOneFound, documentTwoFound, documentThreeFound, documentFourFound bool

	versionTen := book.Versions["10.0.0"]
	if versionTen.Version != "10.0.0" {
		t.Errorf("Expected 10.0.0, got %s", book.Versions["10.0.0"].Version)
	}

	if versionTen.DestPath != fmt.Sprintf("%s/10.0.0", destPath) {
		t.Errorf("Expected %s, got %s", fmt.Sprintf("%s/10.0.0", destPath), versionTen.DestPath)
	}

	if len(versionTen.Docs) != 4 {
		t.Errorf("Expected 3 docs, got %d", len(versionTen.Docs))
	}

	expectedSourcePath := fmt.Sprintf("%s/10.0.0", sourcePath)
	sharedSourcePath := fmt.Sprintf("%s/Shared", sourcePath)
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

		if doc.Storage.Name() == "DocumentFour.md" {
			if !strings.EqualFold(doc.SourcePath, expectedSourcePath) {
				t.Errorf("Expected %s, got %s", expectedSourcePath, doc.SourcePath)
			} else {
				documentFourFound = true
			}
		}
	}

	if !(documentOneFound && documentTwoFound && documentThreeFound && documentFourFound) {
		t.Errorf("Expected DocumentOne.md, DocumentTwo.md, DocumentThree.md, DocumentFour.md in v10.0.0")
	}
}

// TestBuildTableOfContents tests the BuildTOC function.
// It creates a fake directory structure and then runs the BuildTOC function.
func TestBuildTableOfContents(t *testing.T) {

	mydir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %s", err)
	}

	sourcePath := strings.Replace(mydir, "internal", "test/source", 1)
	destPath := strings.Replace(mydir, "internal", fmt.Sprintf("test/docs/%s", uuid.New().String()), 1)

	var sources = sourceTestDataBuilder(sourcePath, mydir)
	var book = sources.BuildBook(destPath)
	if book == nil {
		t.Errorf("Error building book: %s", err)
	}

	book.BuildTOC()

	checkTOC(t, book, sourcePath, destPath)
}

func checkTOC(t *testing.T, book *Book, sourcePath string, destPath string) {

	if book.Root.Summary.Storage == nil || book.Root.Summary.Storage.Name() != "SUMMARY.md" {
		t.Errorf("Expected a file for SUMMARY.md")
	}

	//In a book, the SUMMARY.md file is always copied to the destpath
	if book.Root.DestPath != destPath {
		t.Errorf("Expected %s, got %s", destPath, book.Root.Summary.SourcePath)
	}

	if book.Root.Summary.SourcePath != sourcePath+"/"+"summary" {
		t.Errorf("Expected %s, got %s", sourcePath+"/"+"summary", book.Root.Summary.SourcePath)
	}

	mdFile, err := os.ReadFile(book.Root.Summary.SourcePath + "/" + book.Root.Summary.Storage.Name())
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}

	ext := parser.CommonExtensions | parser.OrderedListStart
	parser := parser.NewWithExtensions(ext)
	renderer := md.NewRenderer()
	doc := markdown.Parse(mdFile, parser)
	got := markdown.Render(doc, renderer)
	toc := fmt.Sprintf("%s", got)

	if toc != "## v9.0.0\n### Brighter Configuration\n* [Document One](/v9.0.0/DocumentOne.md)\n* [Document Two](/v9.0.0/DocumentTwo.md)\n### Darker Configuration\n* [Document Three](/v9.0.0/DocumentThree.md)\n## v10.0.0\n### Brighter Configuration\n* [Document One](/v10.0.0/DocumentOne.md)\n* [Document Two](/v10.0.0/DocumentTwo.md)\n* [Document Four](/v10.0.0/DocumentFour.md)\n### Darker Configuration\n* [Document Three](/v10.0.0/DocumentThree.md)\n" {
		t.Errorf("Expected empty string, got %s", toc)
	}

}

func TestBookCreation(t *testing.T) {

	myDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %s", err)
	}

	sourcePath := strings.Replace(myDir, "internal", "test/source", 1)
	destPath := strings.Replace(myDir, "internal", fmt.Sprintf("test/docs/%s", uuid.New().String()), 1)

	sources := NewSources()
	err = sources.FindFromPath(sourcePath)
	if err != nil {
		t.Errorf("Error finding sources: %s", err)
	}

	book := sources.BuildBook(destPath)

	err = book.Publish()
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

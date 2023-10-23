package book

import (
	"fmt"
	"github.com/brightercommand/Rewind/internal/pages"
	"github.com/brightercommand/Rewind/internal/sources"
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

	sourcePath := strings.Replace(mydir, "internal/book", "test/filecopy", 1)
	destPath := strings.Replace(mydir, "internal/book", fmt.Sprintf("test/filecopy/%s", uuid.New().String()), 1)

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
// NOTE: Test can fail because ordering is not guaranteed in maps. A re-run typically fixes the issue.
func TestBookBuilder(t *testing.T) {

	mydir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %s", err)
	}

	sourcePath := strings.Replace(mydir, "internal/book", "test/source", 1)
	destPath := strings.Replace(mydir, "internal/book", fmt.Sprintf("test/docs/%s", uuid.New().String()), 1)

	var src = sources.SourceTestDataBuilder(sourcePath, mydir)
	book, err := MakeBook(src, destPath)
	if err != nil {
		t.Errorf("Error building book: %s", err)
	}

	checkBookRoot(t, book, destPath, src, sourcePath)

	if len(book.Versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(book.Versions))
	}

	checkVersion9(t, book, destPath, sourcePath)

	checkVersion10(t, book, destPath, sourcePath)

	checkTOC(t, book, sourcePath, destPath)

	// Remove the temp directory used for the summary file
	err = book.ClearWorkDir()
	if err != nil {
		t.Errorf("Error removing directory: %s", err)
	}
}

func checkBookRoot(t *testing.T, book *Book, destPath string, sources *sources.Sources, sourcePath string) {
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

	versionNine := book.Versions["9"]
	if versionNine.Version != "9" {
		t.Errorf("Expected 9, got %s", book.Versions["9"].Version)
	}

	if versionNine.DestPath != fmt.Sprintf("%s/contents/9", destPath) {
		t.Errorf("Expected %s, got %s", fmt.Sprintf("%s/contents/9", destPath), versionNine.DestPath)
	}

	if len(versionNine.Docs) != 3 {
		t.Errorf("Expected 3 docs, got %d", len(versionNine.Docs))
	}

	var expectedSourcePath = fmt.Sprintf("%s/9", sourcePath)
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
		t.Errorf("Expected DocumentOne.md, DocumentTwo.md, DocumentThree.md in version 9")
	}
	return
}

func checkVersion10(t *testing.T, book *Book, destPath string, sourcePath string) {
	var documentOneFound, documentTwoFound, documentThreeFound, documentFourFound bool

	versionTen := book.Versions["10"]
	if versionTen.Version != "10" {
		t.Errorf("Expected version 10, got version %s", book.Versions["10"].Version)
	}

	if versionTen.DestPath != fmt.Sprintf("%s/contents/10", destPath) {
		t.Errorf("Expected %s, got %s", fmt.Sprintf("%s/contents/10", destPath), versionTen.DestPath)
	}

	if len(versionTen.Docs) != 4 {
		t.Errorf("Expected 3 docs, got %d", len(versionTen.Docs))
	}

	expectedSourcePath := fmt.Sprintf("%s/10", sourcePath)
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
		t.Errorf("Expected DocumentOne.md, DocumentTwo.md, DocumentThree.md, DocumentFour.md in version 10")
	}
}

func checkTOC(t *testing.T, book *Book, sourcePath string, destPath string) {

	//In a book, the SUMMARY.md file is always copied to the destpath
	if book.Root.DestPath != destPath {
		t.Errorf("Expected %s, got %s", destPath, book.Root.DestPath)
	}

	if book.Root.SourcePath != sourcePath {
		t.Errorf("Expected %s, got %s", sourcePath, book.Root.SourcePath)
	}

	if book.Root.WorkDir == "" {
		t.Errorf("Expected a WorkDir for the summary file")
	}

	mdFile, err := os.ReadFile(book.Root.WorkDir + "/" + pages.SummaryFileName)
	if err != nil {
		t.Errorf("Error reading file: %s", err)
	}

	ext := parser.CommonExtensions | parser.OrderedListStart
	ps := parser.NewWithExtensions(ext)
	renderer := md.NewRenderer()
	doc := markdown.Parse(mdFile, ps)
	got := markdown.Render(doc, renderer)
	toc := fmt.Sprintf("%s", got)

	expectedTOC := "## 9\n### Brighter Configuration\n* [Document One](/contents/9/DocumentOne.md)\n    * [Document Two](/contents/9/DocumentTwo.md)\n### Darker Configuration\n* [Document Three](/contents/9/DocumentThree.md)\n## 10\n### Brighter Configuration\n* [Document One](/contents/10/DocumentOne.md)\n    * [Document Two](/contents/10/DocumentTwo.md)\n* [Document Four](/contents/10/DocumentFour.md)\n### Darker Configuration\n* [Document Three](/contents/10/DocumentThree.md)\n"
	if toc != expectedTOC {
		t.Errorf("Expected %s, got %s", expectedTOC, toc)
	}

}

func TestBookCreation(t *testing.T) {

	myDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %s", err)
	}

	sourcePath := strings.Replace(myDir, "internal", "test/source", 1)
	destPath := strings.Replace(myDir, "internal", fmt.Sprintf("test/docs/%s", uuid.New().String()), 1)

	src := sources.NewSources()
	err = src.FindFromPath(sourcePath)
	if err != nil {
		t.Errorf("Error finding sources: %s", err)
	}

	book, err := MakeBook(src, destPath)
	if err != nil {
		t.Errorf("Error building book: %s", err)
	}

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
		} else if entry.Name() == "..gitbook.yaml" {
			gitbookFound = true
		} else if entry.IsDir() {
			if entry.Name() == "contents" {
				contents, err := os.ReadDir(destPath + "/contents")
				if err != nil {
					t.Errorf("Error reading directory: %s", err)
				}
				for _, entry := range contents {
					if entry.Name() == "9" {
						v9entries, err := os.ReadDir(fmt.Sprintf("%s/contents/9", destPath))
						if err != nil {
							t.Errorf("Error reading directory: %s", err)
							v9Found = false
						} else {
							v9Found = findFiles(v9entries)
						}
					} else if entry.Name() == "10" {
						v10entries, err := os.ReadDir(fmt.Sprintf("%s/contents/10", destPath))
						if err != nil {
							t.Errorf("Error reading directory: %s", err)
							v10Found = false
						} else {
							v10Found = findFiles(v10entries)
						}
					}
				}
			}
		}
	}

	if summaryFound == false {
		t.Errorf("Expected SUMMARY.md")
	}

	if gitbookFound == false {
		t.Errorf("Expected ..gitbook.yaml")
	}

	if v9Found == false {
		t.Errorf("Expected 9")
	}

	if v10Found == false {
		t.Errorf("Expected 10")
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

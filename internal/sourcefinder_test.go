package internal

import (
	"os"
	"strings"
	"testing"
)

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

	if len(docs) != 2 {
		t.Errorf("Expected 2 documents, got %d", len(docs))
	} else {

		if docs[0].Name() != "DocumentOne.md" {
			t.Errorf("Expected DocumentOne.md, got %s", docs[0])
		}

		if docs[1].Name() != "DocumentTwo.md" {
			t.Errorf("Expected DocumentTwo.md, got %s", docs[1])
		}
	}

	versions := sources.Versions
	if len(versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(versions))
	} else {

		if versions[0].Version != "9.0.0" {
			t.Errorf("Expected 9.0.0, got %s", versions[0].Version)
		}
		if versions[1].Version != "10.0.0" {
			t.Errorf("Expected 1.0.1, got %s", versions[1].Version)
		}
	}
}

package sources

import (
	"os"
	"strings"
	"testing"
)

func TestFindSources(t *testing.T) {

	myDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %s", err)
	}

	sourcePath := strings.Replace(myDir, "internal/sources", "test/source", 1)

	sources := NewSources()
	err = sources.FindFromPath(sourcePath)
	if err != nil {
		t.Errorf("Error finding sources: %s", err)
	}

	checkRoot(t, sources)

	checkShared(t, sources)

	checkVersions(t, sources)
}

func checkRoot(t *testing.T, sources *Sources) {
	root := sources.Root

	if root.GitBook.Storage.Name() != ".gitbook.yaml" {
		t.Errorf("Expected .gitbook.yaml")
	}
}

func checkShared(t *testing.T, sources *Sources) {

	if sources.Shared == nil {
		t.Errorf("Expected shared documents")
		return
	}

	if sources.Shared.TOC.Storage == nil || sources.Shared.TOC.Storage.Name() != ".toc.yaml" {
		t.Errorf("Expected .toc.yaml, got %s", sources.Shared.TOC)
	}

	docs := sources.Shared.Docs

	if len(docs) != 3 {
		t.Errorf("Expected 3 documents, got %d", len(docs))
	} else {

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
	return
}

func checkVersions(t *testing.T, sources *Sources) {
	versions := sources.Versions
	if len(versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(versions))
	} else {

		if versions["10.0.0"].Version != "10.0.0" {
			t.Errorf("Expected 9.0.0, got %s", versions["10.0.0"].Version)
		} else {
			version := versions["10.0.0"]

			if version.TOC.Storage == nil || version.TOC.Storage.Name() != ".toc.yaml" {
				t.Errorf("Expected .toc.yaml, got %s", version.TOC)
			}

			docs := version.Docs
			if len(docs) != 2 {
				t.Errorf("Expected 3 documents, got %d", len(docs))
			} else {

				if docs["DocumentOne.md"].Storage.Name() != "DocumentOne.md" {
					t.Errorf("Expected DocumentOne.md, got %s", docs["DocumentOne.md"])
				}

				if docs["DocumentFour.md"].Storage.Name() != "DocumentFour.md" {
					t.Errorf("Expected DocumentFour.md, got %s", docs["DocumentFour.md"])
				}
			}
		}
		if versions["9.0.0"].Version != "9.0.0" {
			t.Errorf("Expected 9.0.0, got %s", versions["9.0.0"].Version)
		} else {
			version := versions["9.0.0"]

			if version.TOC.Storage == nil || version.TOC.Storage.Name() != ".toc.yaml" {
				t.Errorf("Expected .toc.yaml, got %s", version.TOC)
			}

			docs := version.Docs
			if len(docs) != 1 {
				t.Errorf("Expected 2 documents, got %d", len(docs))
			} else {

				if docs["DocumentTwo.md"].Storage.Name() != "DocumentTwo.md" {
					t.Errorf("Expected DocumentTwo.md, got %s", docs["DocumentTwo.md"])
				}
			}
		}
	}
}
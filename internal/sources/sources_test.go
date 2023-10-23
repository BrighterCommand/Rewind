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

	if root.ReadMe.Storage.Name() != "README.md" {
		t.Errorf("Expected README.md")
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

	images := sources.Shared.Images

	if len(images) != 2 {
		t.Errorf("Expected 1 image, got %d", len(images))
	} else {
		if images["ImageOne.png"].Storage.Name() != "ImageOne.png" {
			t.Errorf("Expected ImageOne.png, got %s", images["ImageOne.png"].Storage.Name())
		}

		if images["ImageTwo.png"].Storage.Name() != "ImageTwo.png" {
			t.Errorf("Expected ImageTwo.png, got %s", images["ImageTwo.png"].Storage.Name())
		}
	}

	return
}

func checkVersions(t *testing.T, sources *Sources) {
	versions := sources.Versions
	if len(versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(versions))
	} else {
		if versions["10"].Version != "10" {
			t.Errorf("Expected 10, got %s", versions["10"].Version)
		} else {
			version := versions["10"]

			if version.TOC.Storage == nil || version.TOC.Storage.Name() != ".toc.yaml" {
				t.Errorf("Expected .toc.yaml, got %s", version.TOC)
			}

			docs := version.Docs
			if len(docs) != 2 {
				t.Errorf("Expected 2 documents, got %d", len(docs))
			} else {

				if docs["DocumentOne.md"].Storage.Name() != "DocumentOne.md" {
					t.Errorf("Expected DocumentOne.md, got %s", docs["DocumentOne.md"])
				}

				if docs["DocumentFour.md"].Storage.Name() != "DocumentFour.md" {
					t.Errorf("Expected DocumentFour.md, got %s", docs["DocumentFour.md"])
				}
			}

			images := version.Images
			if len(images) != 2 {
				t.Errorf("Expected 2 image, got %d", len(images))
			} else {
				if images["ImageOne.png"].Storage.Name() != "ImageOne.png" {
					t.Errorf("Expected ImageOne.png, got %s", images["ImageOne.png"].Storage.Name())
				}

				if images["ImageThree.png"].Storage.Name() != "ImageThree.png" {
					t.Errorf("Expected ImageThree.png, got %s", images["ImageThree.png"].Storage.Name())
				}
			}
		}

		if versions["9"].Version != "9" {
			t.Errorf("Expected 9, got %s", versions["9"].Version)
		} else {
			version := versions["9"]

			if version.TOC.Storage == nil || version.TOC.Storage.Name() != ".toc.yaml" {
				t.Errorf("Expected .toc.yaml, got %s", version.TOC)
			}

			docs := version.Docs
			if len(docs) != 1 {
				t.Errorf("Expected 1 document, got %d", len(docs))
			} else {

				if docs["DocumentTwo.md"].Storage.Name() != "DocumentTwo.md" {
					t.Errorf("Expected DocumentTwo.md, got %s", docs["DocumentTwo.md"])
				}
			}

			images := version.Images
			if len(images) != 2 {
				t.Errorf("Expected 2 images, got %d", len(images))
			} else {
				if images["ImageThree.png"].Storage.Name() != "ImageThree.png" {
					t.Errorf("Expected ImageThree.png, got %s", images["ImageThree.png"].Storage.Name())
				}

				if images["ImageFour.png"].Storage.Name() != "ImageFour.png" {
					t.Errorf("Expected ImageFour.png, got %s", images["ImageFour.png"].Storage.Name())
				}
			}
		}
	}
}

package internal

import "os"

type Assets struct {
	Path string
	Docs []os.DirEntry
}

type Root struct {
	Path    string
	Summary os.DirEntry
	GitBook os.DirEntry
}

type Shared struct {
	Assets *Assets
}

type Version struct {
	Assets  *Assets
	Version string
}

type Sources struct {
	Root     *Root
	Shared   *Shared
	Versions []Version
}

// FindSources finds the sources for a book.
// It takes a root directory as a string.
// It returns a Sources struct.
// The Sources struct contains the root directory, shared documents, and versioned documents.
func FindSources(root string) (sources *Sources, err error) {

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	sources = &Sources{
		Root:     &Root{Path: root},
		Shared:   &Shared{},
		Versions: []Version{},
	}

	for _, entry := range entries {
		if entry.Name() == "SUMMARY.md" {
			sources.Root.Summary = entry
		} else if entry.Name() == ".gitbook.yaml" {
			sources.Root.GitBook = entry
		} else if entry.IsDir() && entry.Name() == "shared" {
			shared, err := findShared(root, entry)
			if err != nil {
				return nil, err
			}
			sources.Shared = shared
		} else if entry.IsDir() {
			version, err := findVersion(root, entry)
			if err != nil {
				return nil, err
			}
			sources.Versions = append(sources.Versions, *version)
		}
	}

	return sources, err
}

// findShared finds the shared documents for the book.
// It takes a directory entry and its path.
// It returns a Shared struct.
func findShared(path string, entry os.DirEntry) (shared *Shared, err error) {
	shared = &Shared{
		Assets: &Assets{
			Docs: []os.DirEntry{},
		},
	}

	sharedPath := path + "/" + entry.Name()

	err = findAssets(sharedPath, shared.Assets)
	if err != nil {
		return shared, err
	}

	return shared, nil

}

// findVersion finds the versioned documents for the book.
// It takes a directory entry and its path.
// It returns a Version struct.
func findVersion(path string, entry os.DirEntry) (version *Version, err error) {

	version = &Version{
		Assets: &Assets{
			Docs: []os.DirEntry{},
		},
		Version: entry.Name(),
	}

	versionPath := path + "/" + entry.Name()

	err = findAssets(versionPath, version.Assets)
	if err != nil {
		return version, err
	}

	return version, nil

}

// findAssets finds the documents for the book.
// It takes a directory entry and an Assets struct.
// It returns an error.
// It will recurse over any subdirectories and place all shared assets at the same level.
func findAssets(path string, assets *Assets) (err error) {
	assets.Path = path

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			assets.Docs = append(assets.Docs, entry)
		} else {
			err = findAssets(path+"/"+entry.Name(), assets)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

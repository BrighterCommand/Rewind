package internal

import (
	"io"
	"os"
)

type Assets struct {
	DestPath   string
	SourcePath string
	Docs       map[string]Doc
}

type Book struct {
	Root     *Root
	Versions map[string]Version
}

type Doc struct {
	Version string
	Storage os.DirEntry
}

type Root struct {
	DestPath   string
	SourcePath string
	Summary    os.DirEntry
	GitBook    os.DirEntry
}

type Sources struct {
	Root     *Root
	Shared   *Shared
	Versions map[string]Version
}

type Shared struct {
	Assets *Assets
}

type Version struct {
	Assets  *Assets
	Version string
}

func (s *Sources) BuildBook(destPath string) *Book {
	book := &Book{
		Root: &Root{
			DestPath:   destPath,
			SourcePath: s.Root.SourcePath,
			Summary:    s.Root.Summary,
			GitBook:    s.Root.GitBook,
		},
		Versions: make(map[string]Version),
	}

	for key, version := range s.Versions {

		var bookVersion = &Version{
			Version: version.Version,
			Assets: &Assets{
				DestPath:   destPath + "/" + version.Version,
				SourcePath: version.Assets.SourcePath,
				Docs:       make(map[string]Doc),
			},
		}

		//copy shared assets first
		for key, doc := range s.Shared.Assets.Docs {
			bookVersion.Assets.Docs[key] = doc
		}

		//now copy assets for this version and overwrite any shared assets with the same name
		for key, doc := range version.Assets.Docs {
			bookVersion.Assets.Docs[key] = doc
		}

		book.Versions[key] = *bookVersion
	}

	return book
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
		Root:     &Root{SourcePath: root},
		Shared:   &Shared{},
		Versions: make(map[string]Version),
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
			sources.Versions[entry.Name()] = *version
		}
	}

	return sources, err
}

// copyFile copies a file from sourcePath to destPath
// It takes two string arguments: sourcePath and destPath
// It creates the destination file if it does not exist
// It overwrites the destination file if it exists
// It returns an error if the copy fails
func copyFile(sourcePath string, destPath string, fileName string) (err error) {
	r, err := os.Open(sourcePath + "/" + fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	//if the directory does not exist, create it
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		err = os.MkdirAll(destPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	w, err := os.Create(destPath + "/" + fileName)
	if err != nil {
		return err
	}

	defer func() {
		if c := w.Close(); c != nil && err == nil {
			err = c
		}
	}()

	_, err = io.Copy(w, r)
	return err
}

// findAssets finds the documents for the book.
// It takes a directory entry and an Assets struct.
// It returns an error.
// It will recurse over any subdirectories and place all shared assets at the same level.
func findAssets(path string, version string, assets *Assets) (err error) {
	assets.SourcePath = path

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			assets.Docs[entry.Name()] = Doc{Version: version, Storage: entry}
		} else {
			err = findAssets(path+"/"+entry.Name(), version, assets)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// findShared finds the shared documents for the book.
// It takes a directory entry and its path.
// It returns a Shared struct.
func findShared(path string, entry os.DirEntry) (shared *Shared, err error) {
	shared = &Shared{
		Assets: &Assets{
			Docs: make(map[string]Doc),
		},
	}

	sharedPath := path + "/" + entry.Name()

	err = findAssets(sharedPath, "Shared", shared.Assets)
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
			Docs: make(map[string]Doc),
		},
		Version: entry.Name(),
	}

	versionPath := path + "/" + entry.Name()

	err = findAssets(versionPath, version.Version, version.Assets)
	if err != nil {
		return version, err
	}

	return version, nil

}

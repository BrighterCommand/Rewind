package book

import (
	"github.com/brightercommand/Rewind/internal/pages"
	"github.com/brightercommand/Rewind/internal/sources"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type Book struct {
	Root     *pages.Root
	Versions map[string]pages.Version
}

func MakeBook(s *sources.Sources, destPath string) (*Book, error) {
	b := &Book{
		Root: &pages.Root{
			DestPath:   destPath,
			SourcePath: s.Root.SourcePath,
			GitBook:    s.Root.GitBook,
		},
		Versions: make(map[string]pages.Version),
	}

	log.Print("Making versions...")
	b.MakeVersions(s)

	log.Print("Making TOC...")
	err := b.MakeTOC(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Book) Publish() error {

	//create root directory
	rootPath := b.Root.DestPath
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		err := os.MkdirAll(rootPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	//copy root files
	err := copyFile(b.Root.WorkDir, b.Root.DestPath, pages.SummaryFileName)
	if err != nil {
		return err
	}

	err = copyFile(b.Root.GitBook.SourcePath, b.Root.DestPath, b.Root.GitBook.Storage.Name())
	if err != nil {
		return err
	}

	//copy versioned files
	for _, version := range b.Versions {

		destPath := version.DestPath
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			err := os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return err
			}
		}

		for _, doc := range version.Docs {
			err = copyFile(doc.SourcePath, destPath, doc.Storage.Name())
			if err != nil {
				return err
			}
		}
	}

	//clean up the temporary files
	err = b.ClearWorkDir()
	if err != nil {
		return err
	}

	return nil
}

func (b *Book) MakeVersions(s *sources.Sources) {
	for key, version := range s.Versions {

		log.Print("Making version " + version.Version + "...")

		var bookVersion = &pages.Version{
			Version:  version.Version,
			DestPath: b.Root.DestPath + "/" + pages.ContentDirName + "/" + version.Version,
			Docs:     make(map[string]pages.Doc),
		}

		log.Print("Copying shared assets...")
		//copy shared assets first
		for key, doc := range s.Shared.Docs {
			bookVersion.Docs[key] = doc
		}

		log.Print("Copying version assets...")
		//now copy assets for this sharedVersion and overwrite any shared assets with the same name
		for key, doc := range version.Docs {
			bookVersion.Docs[key] = doc
		}

		b.Versions[key] = *bookVersion
	}
}

func (b *Book) MakeTOC(s *sources.Sources) error {

	entries, err := b.buildEntries(s)
	if err != nil {
		return err
	}

	orderedTocs := entries.Sort()

	//create a workdir
	b.Root.WorkDir, err = os.MkdirTemp(s.Root.SourcePath, "summary")

	//write the summary
	summary, err := os.Create(b.Root.WorkDir + "/" + pages.SummaryFileName)
	if err != nil {
		return err
	}
	defer summary.Close()

	mg := newMarkdownGenerator()
	err = mg.GenerateSummary(orderedTocs, summary)

	return nil
}

func (b *Book) ClearWorkDir() error {
	return os.RemoveAll(b.Root.WorkDir)
}

func (b *Book) buildEntries(s *sources.Sources) (*pages.VersionedToc, error) {

	log.Print("Building TOC entries...")
	var summary = pages.VersionedToc{
		Contents: make(map[string]*pages.Toc),
	}

	shared, err := b.loadSharedEntries(s)
	if err != nil {
		return nil, err
	}

	for _, version := range s.Versions {

		versionEntries, err := b.buildVersionEntries(shared, version)
		if err != nil {
			return nil, err
		}
		summary.Contents[version.Version] = versionEntries
	}
	return &summary, nil
}

func (b *Book) loadSharedEntries(s *sources.Sources) (*pages.Toc, error) {

	log.Print("Loading shared entries from " + s.Shared.TOC.SourcePath + "/" + s.Shared.TOC.Storage.Name() + "...")
	file, err := os.ReadFile(s.Shared.TOC.SourcePath + "/" + s.Shared.TOC.Storage.Name())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	shared := pages.Toc{
		Sections: make(map[string]pages.TOCSection),
	}

	err = yaml.Unmarshal(file, &shared)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &shared, nil
}

func (b *Book) buildVersionEntries(shared *pages.Toc, version pages.Version) (*pages.Toc, error) {

	log.Print("Building version entries for " + version.Version + "...")
	//merge the shared and versioned information
	versionedSections := pages.Toc{
		Sections: make(map[string]pages.TOCSection),
	}

	b.addSharedSections(shared, versionedSections)

	err := b.addVersionedSections(version, versionedSections)
	if err != nil {
		return nil, err
	}
	return &versionedSections, nil
}

func (b *Book) addSharedSections(shared *pages.Toc, versionedEntries pages.Toc) {

	log.Print("Adding shared TOC sections...")
	//add the shared information for each configuration
	for key, section := range shared.Sections {
		versionedSection := pages.TOCSection{
			Order: section.Order,
		}
		versionedSectionEntries := make([]pages.TOCEntry, 0)
		for _, sharedEntry := range section.Entries {
			versionedSectionEntries = append(versionedSectionEntries, sharedEntry)
		}
		versionedSection.Entries = versionedSectionEntries
		versionedEntries.Sections[key] = versionedSection
	}
}

func (b *Book) addVersionedSections(version pages.Version, versionedTOC pages.Toc) error {

	log.Print("Loading versioned TOC sections for " + version.TOC.SourcePath + "/" + version.TOC.Storage.Name() + "...")
	//read the versioned information
	file, err := os.ReadFile(version.TOC.SourcePath + "/" + version.TOC.Storage.Name())
	if err != nil {
		log.Fatal(err)
		return err
	}

	versioned := pages.Toc{
		Sections: make(map[string]pages.TOCSection),
	}

	err = yaml.Unmarshal(file, &versioned)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//for each ection in the version
	for v, i := range versioned.Sections {
		//iterate over the shared toc sections we have already added
		sectionExists := false
		for s, j := range versionedTOC.Sections {
			//if we already have that section
			if v == s {

				sectionExists = true
				j.Order = i.Order

				//we already have this section and map values need to be re-assigned not changed
				var entries = make([]pages.TOCEntry, 0)

				//add the shared entries
				for _, sharedEntry := range j.Entries {
					entries = append(entries, sharedEntry)
				}

				//iterate over versioned entries
				for _, versionedEntry := range i.Entries {
					found := false
					for k, se := range entries {
						//if we have a matching shared entry replace it with the versioned one
						if se.Name == versionedEntry.Name {
							entries[k] = versionedEntry
							found = true
							break
						}
					}
					//it was a new entry add it
					if !found {
						entries = append(entries, versionedEntry)
					}
				}

				j.Entries = entries
			}
		}
		//if it is not an existing section, just add it
		if !sectionExists {
			versionedTOC.Sections[v] = i
		}
	}
	return nil
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
	defer func(r *os.File) {
		_ = r.Close()
	}(r)

	//if the directory does not exist, create it
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		err = os.MkdirAll(destPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	log.Print("Copying " + sourcePath + "/" + fileName + " to " + destPath + "/" + fileName + "...")

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

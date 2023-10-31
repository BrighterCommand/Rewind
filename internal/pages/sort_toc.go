package pages

import "strconv"

func (t *VersionedToc) Sort() (ordered []OrderedVersionTocs) {

	//Turn a VersionedToc into an array of OrderedVersionToc
	for versionName, toc := range t.Contents {
		order, _ := strconv.Atoi(versionName)
		orderedVersion := OrderedVersionTocs{
			Version: versionName,
			Order:   order,
		}

		for sectionName, section := range toc.Sections {
			orderedSection := OrderedTocSection{
				Name:    sectionName,
				Section: section,
			}
			orderedVersion.Sections = append(orderedVersion.Sections, orderedSection)
		}

		ordered = append(ordered, orderedVersion)

	}

	//Now sort the OrderedVersionToc by their order fields

	//sort the versions
	for i := 0; i < len(ordered); i++ {
		for j := i + 1; j < len(ordered); j++ {
			if ordered[i].Order > ordered[j].Order {
				ordered[i], ordered[j] = ordered[j], ordered[i]
			}
		}
	}

	//now for each of the versions, sort their sections
	for i := 0; i < len(ordered); i++ {
		for k := 0; k < len(ordered[i].Sections); k++ {
			for l := k + 1; l < len(ordered[i].Sections); l++ {
				if ordered[i].Sections[k].Section.Order > ordered[i].Sections[l].Section.Order {
					ordered[i].Sections[k], ordered[i].Sections[l] = ordered[i].Sections[l], ordered[i].Sections[k]
				}
			}
		}
	}

	// now sort the entries in those sections
	for i := 0; i < len(ordered); i++ {
		for j := 0; j < len(ordered[i].Sections); j++ {
			for k := 0; k < len(ordered[i].Sections[j].Section.Entries); k++ {
				for l := k + 1; l < len(ordered[i].Sections[j].Section.Entries); l++ {
					if ordered[i].Sections[j].Section.Entries[k].Order > ordered[i].Sections[j].Section.Entries[l].Order {
						ordered[i].Sections[j].Section.Entries[k], ordered[i].Sections[j].Section.Entries[l] = ordered[i].Sections[j].Section.Entries[l], ordered[i].Sections[j].Section.Entries[k]
					}
				}
			}
		}
	}

	return ordered
}

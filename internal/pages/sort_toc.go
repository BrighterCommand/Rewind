package pages

func (t *Toc) Sort() (ordered []OrderedTocSection) {

	for sectionName, section := range *t {
		orderedSection := OrderedTocSection{
			Name:    sectionName,
			Section: section,
		}
		ordered = append(ordered, orderedSection)
	}

	//sort the sections
	for i := 0; i < len(ordered); i++ {
		for j := i + 1; j < len(ordered); j++ {
			if ordered[i].Section.Order > ordered[j].Section.Order {
				ordered[i], ordered[j] = ordered[j], ordered[i]
			}
		}
	}

	//sort the entries
	for i := 0; i < len(ordered); i++ {
		for k := 0; k < len(ordered[i].Section.Entries); k++ {
			for l := k + 1; l < len(ordered[i].Section.Entries); l++ {
				if ordered[i].Section.Entries[k].Order > ordered[i].Section.Entries[l].Order {
					ordered[i].Section.Entries[k], ordered[i].Section.Entries[l] = ordered[i].Section.Entries[l], ordered[i].Section.Entries[k]
				}
			}
		}
	}

	return ordered
}

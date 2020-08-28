package mcfile

// st3_Refs gathers all the various types of internal and
// external references.
// - DoBlockList()
// - DoGLinks()
// - DoTableOfContents()
func (p *MCFile) st3_Refs() *MCFile {
	if p.GetError() != nil {
		return p
	}
	println("--> (3) Refs")
	return p.DoBlockList().DoGLinks().DoTableOfContents()
}

// DoBlockList is Step 3a
func (p *MCFile) DoBlockList() *MCFile {
	if p.GetError() != nil {
		return p
	}
	switch p.FileType() {
	case "XML":
		println("TODO> 3a. DoBlockList XML")
	case "HTML":
		println("TODO> 3a. DoBlockList HTML")
	case "MKDN":
		println("TODO> 3a. DoBlockList MKDN")
	}
	return p
}

// DoGLinks is Step 3b
func (p *MCFile) DoGLinks() *MCFile {
	if p.GetError() != nil {
		return p
	}
	switch p.FileType() {
	case "XML":
		p.GatherXmlGLinks()
	case "HTML":
		println("TODO> 3b. DoGLinks HTML")
	case "MKDN":
		println("TODO> 3b. DoGLinks MKDN")
	}
	return p
}

// DoTableOfContents is Step 3c
func (p *MCFile) DoTableOfContents() *MCFile {
	if p.GetError() != nil {
		return p
	}
	switch p.FileType() {
	case "XML":
		println("TODO> 3c. DoTableOfContents XML")
	case "HTML":
		println("TODO> 3c. DoTableOfContents HTML")
	case "MKDN":
		println("TODO> 3c. DoTableOfContents MKDN")
	}
	return p
}

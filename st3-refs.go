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
	case "MKDN":
		println("TODO> 3c. DoTableOfContents MKDN")
	}
	return p
}

/*
	switch p.TypeSpecific.(type) {
	case TypeMkdn:
		if p.MType[0] != "md" {
			s := "Not markdown (Mtype[0]!=\"md\") ?!"
			logerr.Println(s)
			p.SetError(errors.New(s))
		}
	case TypeXml:
		if !(p.IsXML && p.MType[0] == "xml") {
			s := "Not XML (Mtype[0]!=\"xml\") ?!"
			logerr.Println(s)
			p.SetError(errors.New(s))
		}
	default:
		panic("SanityCheck failed")
	}
	return p
}
*/

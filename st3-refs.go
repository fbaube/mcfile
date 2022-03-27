package mcfile

import L "github.com/fbaube/mlog"

// st3_Refs gathers all the various types of internal and
// external references.
// - DoBlockList()
// - DoGLinks()
// - DoTableOfContents()
//
func (p *Contentity) st3_Refs() *Contentity {
	if p.HasError() {
		return p
	}
	p.logStg = "33"
	p.L(LProgress, "Refs")
	return p.DoBlockList().DoGLinks().DoTableOfContents()
}

// DoBlockList is Step 3a
//
func (p *Contentity) DoBlockList() *Contentity {
	if p.HasError() {
		return p
	}
	switch p.FileType() {
	case "XML":
		L.L.Warning("TODO> 3a. DoBlockList XML")
	case "HTML":
		L.L.Warning("TODO> 3a. DoBlockList HTML")
	case "MKDN":
		L.L.Warning("TODO> 3a. DoBlockList MKDN")
	}
	return p
}

// DoGLinks is Step 3b
//
func (p *Contentity) DoGLinks() *Contentity {
	if p.HasError() {
		return p
	}
	switch p.FileType() {
	case "XML":
		p.GatherXmlGLinks()
	case "HTML":
		L.L.Warning("TODO> 3b. DoGLinks HTML")
	case "MKDN":
		L.L.Warning("TODO> 3b. DoGLinks MKDN")
	}
	return p
}

// DoTableOfContents is Step 3c
//
func (p *Contentity) DoTableOfContents() *Contentity {
	if p.HasError() {
		return p
	}
	switch p.FileType() {
	case "XML":
		L.L.Warning("TODO> 3c. DoTableOfContents XML")
	case "HTML":
		L.L.Warning("TODO> 3c. DoTableOfContents HTML")
	case "MKDN":
		L.L.Warning("TODO> 3c. DoTableOfContents MKDN")
	}
	return p
}

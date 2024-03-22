package mcfile

import (
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
)

// st3_Refs gathers all the various types of internal and
// external references. Summary:
//   - 3a DoBlockList() // make list of blocks
//   - 3b DoGLinks()    // gather links
//   - 3c DoTableOfContents() // make ToC
//
// Some elms are not clearly either BLK or INL, so we
// do (or should) permit an indeterminate third state.
// .
func (p *Contentity) st3_Refs() *Contentity {
	if p.HasError() {
		return p
	}
	p.logStg = "33"
	p.L(LProgress, "=== 33:Refs ===")
	return p.DoBlockList().DoGLinks().DoTableOfContents()
}

// DoBlockList makes a list of all the nodes that are
// blocks, so that they cn be traversed for rendering,
// and targeted for references.
// .
func (p *Contentity) DoBlockList() *Contentity {
	if p.HasError() {
		return p
	}
	switch p.MarkupTypeOfMType() {
	case SU.MU_type_XML:
		L.L.Warning("TODO> 3a. DoBlockList XML")
	case SU.MU_type_HTML:
		L.L.Warning("TODO> 3a. DoBlockList HTML")
	case SU.MU_type_MKDN:
		L.L.Warning("TODO> 3a. DoBlockList MKDN")
	}
	return p
}

// DoGLinks gathers links.
// .
func (p *Contentity) DoGLinks() *Contentity {
	if p.HasError() {
		return p
	}
	switch p.MarkupTypeOfMType() {
	case SU.MU_type_XML:
		L.L.Info("Calling GatherXmlGLinks...")
		p.GatherXmlGLinks()
		L.L.Info("Called! GatherXmlGLinks")
	case SU.MU_type_HTML:
		L.L.Warning("TODO> 3b. DoGLinks HTML")
	case SU.MU_type_MKDN:
		L.L.Warning("TODO> 3b. DoGLinks MKDN")
	}
	return p
}

// DoTableOfContents makes a ToC.
// .
func (p *Contentity) DoTableOfContents() *Contentity {
	if p.HasError() {
		return p
	}
	switch p.MarkupTypeOfMType() {
	case SU.MU_type_XML:
		L.L.Warning("TODO> 3c. DoTableOfContents XML")
	case SU.MU_type_HTML:
		L.L.Warning("TODO> 3c. DoTableOfContents HTML")
	case SU.MU_type_MKDN:
		L.L.Warning("TODO> 3c. DoTableOfContents MKDN")
	}
	return p
}

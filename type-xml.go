package mcfile

// propserties of a link: tag, att, text raw, ref abspath

import (
	"bufio"
	S "strings"

	// "github.com/fbaube/db"
	"github.com/fbaube/gparse"
	"github.com/fbaube/gtree"
	XM "github.com/fbaube/xmlmodels"
)

// The food chain:
// RelFilePath
// AbsFilePath
// CheckedPath
// MCFile
// - TypeSpecific = TypeXml OR TypeMkdn OR TypeHtml
// GTree
// ForesTree (with x-ref sources & targets, etc)

// NOTE We always create a MCFile, so it is a logical place
// to store a GTokenization and a GTree.

type TypeXml struct {

	// The article about go types for functions
	// MAKE BLOCK LIST
	// SORT OUT RESOLUTION OF GLinks
	// GATHER ToC ELEMENTS
	// Separate the XML types into a map of callback functions

	// XmlItems is (DOCS) IDs & IDREFs, (DTDs) Elm defs (incl. Att defs) & Ent defs
	// *xmlfile.XmlItems
	// *IDinfo
	// it is not precisely defined how to handle relative paths in external
	// IDs and entity substitutions, so we need to maintain this list.
	// TODO EntSearchDirs []string

	// GEnts is "ENTITY"" directives (both with "%" and without).
	GEnts map[string]*gparse.GEnt
	// DElms is "ELEMENT" directives.
	DElms map[string]*gtree.GTag
	// TODO Maybe also add maps for NOTs (Notations)
}

// TryXmlPreamble creates and sets `MCFile.XmlPreamble` *only* if a preamble
// is found. It cheats by sneaking a peek at the first line of the content.
// Calling it "Try" indicates that failure to find a preamble is not fatal.
func (p *MCFile) TryXmlPreamble() *MCFile {
	var e error
	var s string
	var pR *bufio.Reader
	// var pX *TypeXml
	var pXP *XM.XmlPreambleFields

	if p.GetError() != nil {
		return p
	}
	// pX = p.TheXml()
	pR = bufio.NewReader(S.NewReader(p.Raw))
	s, e = pR.ReadString('\n')
	// Quick failure ?
	if !S.HasPrefix(s, "<?xml ") {
		return p
	}
	pXP, e = XM.NewXmlPreambleFields(p.Raw)
	// NOTE An error is not fatal! Not here anyways.
	if e != nil {
		println("==> TryXmlPreamble:", e.Error())
		return nil
	}
	p.XmlPreambleFields = pXP
	p.IsXml = 1
	return p
}

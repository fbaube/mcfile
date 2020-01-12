package mcfile

// propserties of a link: tag, att, text raw, ref abspath

import (
	"bufio"
	S "strings"

	"github.com/fbaube/gparse"
	"github.com/fbaube/gtree"
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

// TypeXml is as granular as it gets (for now) for XML files.
//
type TypeXml struct {
	XmlContype
	// XmlPreamble (i.e. <?xml ...?> ) is non-nil
	// IFF we know *definitively* that we have an XML file.
	// If it is `nil` then when writing the file back out in
	// canonical format, also write out `encoding.xml.Header`,
	// which is the default preamble defined in the Go stdlib.
	*gparse.XmlPreamble
	// TagDefCt is for DTD-type files (.dtd, .mod, .ent)
	TagDefCt int // Nr of <!ELEMENT ...>
	// XmlDoctype is non-nil IFF a DOCTYPE directive was found
	*gparse.XmlDoctype
	// CmtdOasisPublicID *ParsedPublicID

	// RootTagIndex int  // Or some sort of pointer into the tree.
	// RootTagCt is >1 means mark the content as a Fragment.
	RootTagCt int
	// These two distinctions are pretty fundamental to processing,
	// so we dedicate booleans to them.
	DoctypeIsDeclared, DoctypeIsGuessed bool

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
	var pX *TypeXml
	var pXP *gparse.XmlPreamble

	if p.GetError() != nil {
		return p
	}
	pX = p.TheXml()
	pR = bufio.NewReader(S.NewReader(p.Raw))
	s, e = pR.ReadString('\n')
	// Quick failure ?
	if !S.HasPrefix(s, "<?xml ") {
		return p
	}
	pXP, e = gparse.NewXmlPreamble(p.Raw)
	// NOTE that an error is not fatal! Not here anyways.
	if e != nil {
		println("==> TryXmlPreamble:", e.Error())
		return nil
	}
	pX.XmlPreamble = pXP
	p.IsXML = true
	return p
}

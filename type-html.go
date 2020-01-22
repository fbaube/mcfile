package mcfile

// properties of a link: tag, att, text raw, ref abspath

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

// TypeHtml is for HTML5 files, but probably not XHTML,
// and probably not for older versions of HTML (3, 4, etc.).
//
type TypeHtml struct {
	// Embedded!
	TypeXml

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
}

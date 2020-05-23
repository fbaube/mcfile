package mcfile

import (
	FU "github.com/fbaube/fileutils"
	"github.com/fbaube/gtree"
)

// GLinkSet is used for (1) intra-file ref resolution,
// (2) inter-file ptr resolution, (3) ToC generation.
type GLinkSet struct {
	// OwnerP points back to the owning struct, so that
	// `GLink`s can be processed easily as simple data structures.
	OwnerP     interface{}
	// OutgoingKeys are outgoing key-based links/references
	OutgoingKeys  []*GLink
	// IncomableKeys are unique key-based definitions that are possible
	// resolution targets (of same or other files' `OutgoingKeys`)
	IncomableKeys []*GLink
	// OutgoingURIs are outgoing URI-based links/references
	OutgoingURIs  []*GLink
	// IncomableURIs are unique URI-based definitions that are possible
	// resolution targets (of same or other files' `OutgoingURIs`)
	IncomableURIs []*GLink
}

// GLink summarizes a link (or key) (or reference) found in markup content.
// It is either URI-based (`href conref id`) or key-based (`key keyref`). It
// applies to all LwDITA formats, but not all fields apply to all LwDITA formats.
type GLink struct {
	// else is Incomable (which are much more numerous)
	IsOutgoing bool
	// "http", "key", "idref", "uri"
	AddressMode string
	// id, idref, href, xref, keyref, etc.
	Att string
	// the tag that has this link-related attribute of interest
	Tag string
	// as redd in during parsing
	Raw string
	// RelFP can be a URI or the resolution of a keyref.
	// "" if target is in same file; NOTE This is relative to the
	// sourcing file, NOT to the current working directory during parsing!
	RelFP string
	// AbsFP can be a URI or the resolution of a keyref.
	// "" if target is in same file
	AbsFP FU.AbsFilePath
	// if present
	TopicID string
	// peeled off from Raw
	FragID string
	// used to narrow in on difficult cases
	Resolved bool
	// the tag where the GLink is defined
	Source *gtree.GTag
	// can be nil: the tag where the GLink is resolved to, quite possibly
	// in another file, which we hope we also have available in memory.
	Target *gtree.GTag
}

/*
// MISCELLANEOUS NOTES
//
// Link varieties:: Outward/Outgoing, Inside
// Link resolution: Inbound/Incoming, Inside
//
// IDs are used in all XML-based markup formats of interest,
// but IDREFs seem to be more restricted - fewer use cases.
// IDs form a flat "namespace", so is we need to type them,
// they need to include string prefixes that specify the
// object domain(s).
// For an ID, HTML uses "<a name=FOO>" but HTML5 has in theory
// deprecated "<a name=...>" and now prefers ""@id=FOO".
// In any case, presented a choice, a reference will prefer
// the "@id" as target over the "<a name=...>".
// In fact, in HTML5 the "<a..>" element no longer has @name.
// Also, in XHTML5, IDs cannot contain an unescaped LT sign (<).
//
// See also https://www.w3.org/TR/REC-xml/#id
//
// See also https://www.w3.org/TR/html4/struct/links.html#h-12.2.3
// You can name a destination anchor with the id attribute:
//   Here's a <a id="anchor-two">photo of my idiot cat.</A>.
// @id and @name share the same name space. Therefore they may not
// both define an anchor with the same name in the same document.
// It is OK to use both attributes to specify an element's unique
// identifier for these elements: A, APPLET, FORM, FRAME, IFRAME,
// IMG, and MAP. When both attributes are used on a single element,
// their values must be identical.
//
// In LwDITA, topics MUST have an ID and maps MAY have an ID.
// Since @id functions as the target of the various key tags,
// they function (in effect) as XML IDREF's.
*/

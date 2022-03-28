package mcfile

import (
	FU "github.com/fbaube/fileutils"
	"github.com/fbaube/gtree"
)

type Flags int

const (
	IsRef      Flags = 1 << iota // 1 << 0 i.e. 0000 0001
	IsExtl                       // 1 << 1 i.e. 0000 0010
	IsURI                        // 1 << 2 i.e. 0000 0100
	IsKey                        // 1 << 3 i.e. 0000 1000
	IsResolved                   // 1 << 4 i.e. 0001 0000
)

func (b Flags) Set(flag Flags) Flags   { return b | flag }
func (b Flags) Reset(flag Flags) Flags { return b & ^flag }
func (b Flags) IsSet(flag Flags) bool  { return b&flag != 0 }

// func Toggle(b, flag Flags) Flags { return b ^ flag }

func (f Flags) String() string {
	var deff = "Def"
	var intl = "Intl"
	var ltyp = "N/A"
	var resd = "unReslvd"
	if f.IsSet(IsRef) {
		deff = "Ref"
	}
	if f.IsSet(IsExtl) {
		intl = "Extl"
	}
	if f.IsSet(IsURI) {
		ltyp = "Uri"
	}
	if f.IsSet(IsKey) {
		ltyp = "Key"
	}
	if f.IsSet(IsResolved) {
		resd = "Resolved"
	}
	return deff + "," + intl + "," + ltyp + "," + resd
}

// GLinks is used for (1) intra-file ref resolution,
// (2) inter-file ptr resolution, (3) ToC generation.
type GLinks struct {
	// OwnerP points back to the owning struct, so that
	// `GLink`s can be processed easily as simple data structures.
	OwnerP interface{}
	// OutgoingKeys are outgoing key-based links/references
	OutgoingKeys []*GLink // (Extl|Intl)KeyRefs
	// IncomableKeys are unique key-based definitions that are possible
	// resolution targets (of same or other files' `OutgoingKeys`)
	IncomableKeys []*GLink // (Extl|Intl)KeyDefs
	// OutgoingURIs are outgoing URI-based links/references
	OutgoingURIs []*GLink // (Extl|Intl)UriRefs
	// IncomableURIs are unique URI-based definitions that are possible
	// resolution targets (of same or other files' `OutgoingURIs`)
	IncomableURIs []*GLink // (Extl|Intl)UriDefs
}

// GLink summarizes a link (or key) (or reference) found in markup content.
// It is either URI-based (`href conref id`) or key-based (`key keyref`). It
// applies to all LwDITA formats, but not all fields apply to all LwDITA formats.
type GLink struct {
	// else is Def (which are much more numerous)
	IsRef bool
	// else is Intl (which are more numerous)
	IsExtl bool
	// "http", "key", "idref", "uri"
	AddressMode string
	// id, idref, href, xref, keyref, etc.
	Att string
	// the tag that has this link-related attribute of interest
	Tag string
	// as redd in during parsing
	Link_raw string
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
	LinkedFrom *gtree.GTag
	// can be nil: the tag where the GLink is resolved to, quite possibly
	// in another file, which we hope we also have available in memory.
	Original *gtree.GTag
}

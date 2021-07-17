package mcfile

import (
	"github.com/fbaube/gtoken"
	"github.com/fbaube/gtree"
	SU "github.com/fbaube/stringutils"
)

// RefineDirectives scans to patch Directives with correct keyword.
func (p *MCFile) RefineDirectives() error {
	// pX := p.TheXml()

	var pTag *gtree.GTag
	for _, pTag = range p.GTags {
		if pTag.TTType != "Dir" {
			continue
		}
		// Here are the directives we are dealing with.
		// Promote the Directive subtype keyword at "keynoun" to "Type", and
		// then promote the first word of "keyargs" (the Name) to "keynoun"".
		// DOCTYPE  Name ExtID [SysID]
		// ENTITY   Name EntDef   // General entity
		// ENTITY % Name EntDef   // Parameter entity (DTD only)
		// ELEMENT  Name contentspec
		// ATTLIST  Name AttDef's
		// NOTATION Name ExtID
		// ilog.Printf("Dir.PRE |%s|%s|", RT.string1, RT.string2)

		pTag.TTType = gtoken.TTType(pTag.Keyword)
		if pTag.TTType == "Dir" {
			panic("YIKES, leftover Dir Tagtype")
		}
		pTag.Keyword, pTag.Otherwords = SU.SplitOffFirstWord(pTag.Otherwords)

		println("    --> RefineDirectives:", pTag.TTType, ",", pTag.Keyword)

		if pTag.TTType == "ENTITY" && pTag.Keyword == "%" {
			pTag.EntityIsParameter = true
			pTag.Keyword, pTag.Otherwords = SU.SplitOffFirstWord(pTag.Otherwords)
			// pXI.GotDtdDecls = true
		}
		// fmt.Printf("Normalized Directive: %+v \n", RT)
		// At this point we are OK:
		// Type:DOCTYPE  keynoun:HTML keyargs:  PUBLIC "-//W3C//DTD HTML 4.0 ..."
		// Type:ENTITY   keynoun:foo  keyargs: "FOO"  entityIsParameter:false
		// Type:ENTITY   keynoun:bar  keyargs: "BAR"  entityIsParameter:true
		// Type:ELEMENT  Name contentspec
		// Type:ATTLIST  Name AttDef's
	}
	return nil
}

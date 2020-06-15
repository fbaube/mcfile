package mcfile

import (
	"github.com/fbaube/gparse"
	"github.com/fbaube/gtoken"
	"github.com/fbaube/gtree"
	SU "github.com/fbaube/stringutils"
	XM "github.com/fbaube/xmlmodels"
	"github.com/pkg/errors"
)

// ProcessMetaGetDoctype can safely assume that the file is XML. We scan
// the file, from the beginning, just far enough to provide the data we seek.
// TODO This should return DOCTYPE Micodo too, cos it
// is even more valid than the results of analysis.
func (p *MCFile) ProcessMetaGetDoctype() error {
	// pX := p.TheXml()

	var pTag *gtree.GTag
	// var pXFM *gparse.XmlFileMeta
	// pX.XmlFileMeta = new(gparse.XmlFileMeta)
	// println("DoXmlFileMetaAndGetDoctype()...")
	// pXFM = pX.XmlFileMeta

	// =============================
	//  Try to process XML PREAMBLE
	// =============================

	// if it occurs, it has to be the very first item in the file.
	// In other words, pET.Tags[0] .
	// (Except perhaps if preceded by a Unicode BOM = byte-order mark.)
	// <?xml version="1.0" encoding="UTF-8" standalone="no"?>
	// This contains: version, encoding, standalone, in this order ONLY.
	// Only comments and PIs may appear btwn XML preamble and a DOCTYPE decl.
	// About character sets:
	// https://www.iana.org/assignments/character-sets/character-sets.xhtml
	// https://www.w3.org/TR/xml/#sec-guessing
	pTag = p.GTags[0]
	if pTag.TTType == "PI" && pTag.Keyword == "xml" {
		// pXI.GotXmlPreamble = true
		// ilog.Printf("Pass 3: Got XML preamble |%s|%s|", RT.string1, RT.string2)
		pp, e := XM.NewXmlPreambleFields(pTag.Otherwords)
		// If this failed, there's somethign seriously wrong, so fail.
		if e != nil {
			return errors.Wrapf(e, "fx.procfmdt.newpreamble<%s>", pTag.Otherwords)
		}
		if gparse.ExtraInfo {
			println("    --> XML preamble processed OK")
		}
		p.XmlPreambleFields = pp
	}

	// =============================
	//  Try to process XML DOCTYPE
	// =============================

	// Scan the entire file, until either we find a DOCTYPE declaration, or
	// we find a start element (after which point a DOCTYPE is illegal).
	var i int
	for i, pTag = range p.GTags {

		if i == 0 && nil != p.XmlPreambleFields {
			continue
		}
		// if RT.rtType == "CD" { fmt.Printf("CDataDebug|%+v| \n", RT) }
		if pTag.TTType == "PI" {
			continue
		} // legal
		if pTag.TTType == "Cmt" {
			continue
		}
		if pTag.TTType == "SE" {
			// TODO These next two stmts are redundant and can instead be sanity checks
			// pXI.GotRootTag = true
			// pXI.RootTagIndex = i
			if nil == p.XmlDoctypeFields {
				println("    --> No DOCTYPE declaration found")
			}
			// panic("doXmlFileMeta weirdness")
			break
		}
		// Pare down the Happy path
		if pTag.TTType != "Dir" {
			continue
		}
		if pTag.Keyword != "DOCTYPE" {
			continue
		}
		// ====================
		//  DOCTYPE PROCESSING
		// ====================
		pDT, e := XM.NewXmlDoctypeFieldsInclMType(pTag.Otherwords)
		if e == nil {
			p.XmlDoctype = XM.XmlDoctype(pTag.Otherwords)
		}

		// println("\t Got doctype")

		if e != nil {
			return errors.Wrapf(e, "XFM.NewXmlDoctype<%s>", pTag.Otherwords)
		}
		p.DoctypeIsDeclared = true
		p.XmlDoctypeFields   = pDT

		return nil
	}
	return nil
}

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

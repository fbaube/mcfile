package mcfile

import (
	"fmt"
	"errors"
)

// st0_Init does pre-processing prep.
// Input: A bare-bones `MCFile` made from a `fileutils.CheckedContent` made
// from a `fileutils.BasicPath` made from a `fileutils.AbsFilePath` made from
// a CLI argument. <br/>
// For `Init()` to work, the `MCFile` *must* refer to a readable
// file, and the field `MCFile.MType` *must* already be set. <br/>
// Output: An `MCFile` that is in-memory and analyzed (shallowly)
// as `XML` or `MKDN` (Markdown) or `HTML`.
// - SetTypeSpecific()
// - SanityCheck()
//
func (p *MCFile) st0_Init() *MCFile {
	if p.GetError() != nil {
		return p
	}
	println("--> (0) Init")
	// panic("TEST PANIC")
	return p.st0a_SetTypeSpecific().st0b_SanityCheck()
}

// st0a_SetTypeSpecific is Stage 0a: it sets `MCFile.TypeSpecific`
// based on `MCFile.FileType()`, which uses `MCFile.MType[]`.
//
func (p *MCFile) st0a_SetTypeSpecific() *MCFile {
	// println("Init:", p.FileType())
	switch p.FileType() {
	case "XML":
		if !p.IsXML {
			panic("Init error: is XML but:!XML?!")
		}
		p.FFSdataP = new(TypeXml)
	case "MKDN":
		if p.IsXML {
			panic("Init error: is Mkdn but:XML?!")
		}
		p.FFSdataP = new(TypeMkdn)
	case "HTML":
		if !p.IsXML {
			panic("Init error: is HTML but:!XML?!")
		}
		p.FFSdataP = new(TypeHtml)
	default:
		panic("Init error: file type: " + p.FileType())
	}
	return p
}

// st0b_SanityCheck is Step 0b: place all preflight checks here.
//
func (p *MCFile) st0b_SanityCheck() *MCFile {
	var errmsg string
	// fmt.Printf("TypeSpecific: set? <%b> type? <%t> \n",
	//  p.TypeSpecific != nil, p.TypeSpecific)
	switch p.FFSdataP.(type) {
	case *TypeMkdn:
		if p.MType[0] != "mkdn" {
			errmsg = "TypeMarkdown: MType[0]!=\"mkdn\" ?!"
		}
	case *TypeXml:
		if !(p.IsXML && p.MType[0] == "xml") {
			errmsg = "TypeXml: MType[0]!=\"xml\" ?!"
		}
	case *TypeHtml:
		if !(p.IsXML && p.MType[0] == "xml" && p.MType[1] == "html") {
			errmsg = "TypeHtml: MType[0,1]!=\"xml:html\" ?!"
		}
	default:
		errmsg = fmt.Sprintf(
			"Unknown type for file-format-specific data struct: %T", p.FFSdataP)
	}
	if errmsg != "" {
		errmsg = "st[0b] " + errmsg
		p.Blare(p.OLP + errmsg)
		p.SetError(errors.New(errmsg))
	}
	return p
}

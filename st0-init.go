package mcfile

import (
	SU "github.com/fbaube/stringutils"
	L "github.com/fbaube/mlog"
)

// st0_Init does pre-processing prep and checks.
//
// Input: A `Contentity` that has had its contents analyzed.
//
// For `st0_Init()` to work,
//   - the `Contentity` *must* refer to a readable file, and
//   - the field `Contentity.MType` *must* already be set
//
// Output: A `Contentity` that is in-memory and analyzed
// (shallowly) as `XML` or `MKDN` (Markdown) or `HTML`.
//
// - SetTypeSpecific()
// - SanityCheck()
// .
func (p *Contentity) st0_Init() *Contentity {
	if p.HasError() {
		return p
	}
	if p.MType == "" {
		// panic("st0_Init: nil MType")
		L.L.Error("st0_Init: nil MType")
	}
	p.logStg = "00"
	p.L(LProgress, "=== 00:Init ===")
	// panic("TEST PANIC in st0")
	return p.st0a_SanityCheck()
}

// st0a_SanityCheck checks that `mcfile.MarkupType()` &
// `mcfile.IsXML()` are OK and that `MCFile.MType[]` is set.
func (p *Contentity) st0a_SanityCheck() *Contentity {
	p.logStg = "0a"
	if p.MarkupType != SU.MU_type_DIRLIKE && p.MType == "" {
		p.SetErrMsg("MType is empty")
	}
	if len(p.MarkupTypeOfMType()) < 3 || len(p.MarkupTypeOfMType()) > 7 {
		panic("BAD MarkupTypeOfMType in st0a: " + string(p.MarkupTypeOfMType()))
	}
	switch p.MarkupTypeOfMType() {
	case SU.MU_type_XML:
		if !p.IsXML() {
			p.SetErrMsg("is XML but: !XML?!")
		}
	case SU.MU_type_MKDN:
		if p.IsXML() {
			p.SetErrMsg("is Mkdn but: XML?!")
		}
	case SU.MU_type_HTML:
		if !p.IsXML() {
			p.SetErrMsg("is HTML but: !XML?!")
		}
	case SU.MU_type_BIN:
		if p.IsXML() {
			// panic("Init error: is BIN but: XML?!")
			p.SetErrMsg("is BIN but: XML?!")
		}
	case SU.MU_type_SQL, SU.MU_type_DIRLIKE:
	     // No problem!
	default:
		p.SetErrMsg("bad/missing contentitype: " +
			string(p.MarkupTypeOfMType()))
	}
	return p
}

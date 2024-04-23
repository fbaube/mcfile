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
	if p.RawMT != SU.MU_type_DIRLIKE && p.MType == "" {
		p.SetError("MType is empty")
	}
	if len(p.MarkupType()) < 3 || len(p.MarkupType()) > 7 {
		panic("BAD MarkupType in st0a: " + string(p.MarkupType()))
	}
	switch p.MarkupType() {
	case SU.MU_type_XML:
		if !p.IsXML() {
			p.SetError("is XML but: !XML?!")
		}
	case SU.MU_type_MKDN:
		if p.IsXML() {
			p.SetError("is Mkdn but: XML?!")
		}
	case SU.MU_type_HTML:
		if !p.IsXML() {
			p.SetError("is HTML but: !XML?!")
		}
	case SU.MU_type_BIN:
		if p.IsXML() {
			// panic("Init error: is BIN but: XML?!")
			p.SetError("is BIN but: XML?!")
		}
	case SU.MU_type_SQL, SU.MU_type_DIRLIKE:
	     // No problem!
	default:
		p.SetError("bad/missing markup type: " +
			string(p.MarkupType()))
	}
	return p
}

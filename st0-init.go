package mcfile

import L "github.com/fbaube/mlog"

// st0_Init does pre-processing prep.
// Input: A `Contentity` that has had its contents analyzed.
//
// For `st0_Init()` to work, the `Contentity` *must* refer to a readable
// file, and the field `Contentity.MType` *must* already be set.
//
// Output: A `Contentity` that is in-memory and analyzed (shallowly)
// as `XML` or `MKDN` (Markdown) or `HTML`.
//
// - SetTypeSpecific()
// - SanityCheck()
//
func (p *Contentity) st0_Init() *Contentity {
	if p.GetError() != nil {
		return p
	}
	p.logStg = "0:"
	L.L.Progress("Init")
	// panic("TEST PANIC")
	return p.st0a_SanityCheck()
}

// st0a_SanityCheck is Stage 0a: it sets `MCFile.TypeSpecific`
// based on `MCFile.FileType()`, which uses `MCFile.MType[]`.
//
func (p *Contentity) st0a_SanityCheck() *Contentity {
	p.logStg = "0a"
	// println("Init:", p.FileType())
	switch p.FileType() {
	case "XML":
		if !p.IsXML() {
			panic("Init error: is XML but:!XML?!")
		}
	case "MKDN":
		if p.IsXML() {
			panic("Init error: is Mkdn but:XML?!")
		}
	case "HTML":
		if !p.IsXML() {
			panic("Init error: is HTML but:!XML?!")
		}
	default:
		L.L.Panic("Init: File type: " + p.FileType())
		panic("st0a_SanityCheck")
	}
	return p
}

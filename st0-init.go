package mcfile

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
	if p.HasError() {
		return p
	}
	p.logStg = "00"
	p.L(LProgress, "00:Init")
	// panic("TEST PANIC in st0")
	return p.st0a_SanityCheck()
}

// st0a_SanityCheck checks that `mcfile.FileType()` &
// `mcfile.IsXML()` are OK and that `MCFile.MType[]` is set.
//
func (p *Contentity) st0a_SanityCheck() *Contentity {
	p.logStg = "0a"
	if p.MType == "" {
		// p.L(LWarning, "MType is empty")
		p.SetErrMsg("MType is empty")
	}
	switch p.FileType() {
	case "XML":
		if !p.IsXML() {
			// panic("Init error: is XML but: !XML?!")
			p.SetErrMsg("Init error: is XML but: !XML?!")
		}
	case "MKDN":
		if p.IsXML() {
			// panic("Init error: is Mkdn but: XML?!")
			p.SetErrMsg("Init error: is Mkdn but: XML?!")
		}
	case "HTML":
		if !p.IsXML() {
			// panic("Init error: is HTML but: !XML?!")
			p.SetErrMsg("Init error: is HTML but: !XML?!")
		}
	case "BIN":
		if p.IsXML() {
			// panic("Init error: is BIN but: XML?!")
			p.SetErrMsg("Init error: is BIN but: XML?!")
		}
	default:
		// L.L.Panic("Bad contentitype: " + p.FileType())
		// p.L(LError, "Bad/missing contentitype: "+p.FileType())
		p.SetErrMsg("Bad/missing contentitype: " + p.FileType())
	}
	return p
}

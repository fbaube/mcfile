package mcfile

// st0_Init does pre-processing prep.
// Input: A bare-bones `Contentity`.
//
// For `Init()` to work, the `Contentity` *must* refer to a readable
// file, and the field `Contentity.MType` *must* already be set.
//
// Output: An `Contentity` that is in-memory and analyzed (shallowly)
// as `XML` or `MKDN` (Markdown) or `HTML`.
// - SetTypeSpecific()
// - SanityCheck()
//
// func (p *MCFile) st0_Init() *MCFile {
func (p *Contentity) st0_Init() *Contentity {
	if p.GetError() != nil {
		return p
	}
	println("--> (0) Init")
	// panic("TEST PANIC")
	return p.st0a_SanityCheck()
}

// st0a_SanityCheck is Stage 0a: it sets `MCFile.TypeSpecific`
// based on `MCFile.FileType()`, which uses `MCFile.MType[]`.
//
// func (p *MCFile) st0a_SanityCheck() *MCFile {
func (p *Contentity) st0a_SanityCheck() *Contentity {
	// println("Init:", p.FileType())
	switch p.FileType() {
	case "XML":
		if !p.IsXML() {
			panic("Init error: is XML but:!XML?!")
		}
		// p.FFSdataP = new(TypeXml)
	case "MKDN":
		if p.IsXML() {
			panic("Init error: is Mkdn but:XML?!")
		}
		// p.FFSdataP = new(TypeMkdn)
	case "HTML":
		if !p.IsXML() {
			panic("Init error: is HTML but:!XML?!")
		}
		// p.FFSdataP = new(TypeHtml)
	default:
		println("==> Init ERROR: file type:", p.FileType())
		panic("st0a_SanityCheck")
	}
	return p
}

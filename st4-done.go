package mcfile

// st4_Done does final cleanup and beautification.
//
func (p *MCFile) st4_Done() *MCFile {
	if p.GetError() != nil {
		return p
	}
	println("--> (4) Done")
	switch p.FileType() {
	case "XML":
		println("TODO> 4. Done XML")
	case "HTML":
		println("TODO> 4. Done HTML")
	case "MKDN":
		println("TODO> 4. Done MKDN")
	}
	return p
}

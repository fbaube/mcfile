package mcfile

import L "github.com/fbaube/mlog" 

// st4_Done does final cleanup and beautification.
//
func (p *MCFile) st4_Done() *MCFile {
	if p.GetError() != nil {
		return p
	}
	// p.L(LProgress, "Done")
	L.L.Progress("Done")
	switch p.FileType() {
	case "XML":
		L.L.Warning("TODO> 4. Done XML")
	case "HTML":
		L.L.Warning("TODO> 4. Done HTML")
	case "MKDN":
		L.L.Warning("TODO> 4. Done MKDN")
	}
	return p
}

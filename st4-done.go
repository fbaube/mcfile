package mcfile

import (
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
)

// st4_Done does final cleanup and beautification.
// .
func (p *Contentity) st4_Done() *Contentity {
	if p.HasError() {
		return p
	}
	// p.L(LProgress, "Done")
	L.L.Progress("=== 44:Done ===")
	switch p.MarkupTypeOfMType() {
	case SU.MU_type_XML:
		// L.L.Warning("TODO> 4. Done XML")
	case SU.MU_type_HTML:
		// L.L.Warning("TODO> 4. Done HTML")
	case SU.MU_type_MKDN:
		// L.L.Warning("TODO> 4. Done MKDN")
	}
	return p
}

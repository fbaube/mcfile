package mcfile

import (
	"fmt"
	"github.com/fbaube/gparse"
	"github.com/fbaube/gtoken"
	"github.com/fbaube/gtree"
	ON "github.com/fbaube/orderednodes"
	"github.com/fbaube/m5db"
	SU "github.com/fbaube/stringutils"
	"io"
)

type ContentityStage func(*Contentity) *Contentity

// For the record, ignore the API of
// https://godoc.org/golang.org/x/net/html#Node

// Contentity is awesome. It includes a ContentityRow,
// which includes an FSItem, which includes an Errer. 
// .
type Contentity struct { // has Raw

     	// Nord provides hierarchical structure, only.
	ON.Nord
	// ContentityRow includes all fields what get persisted 
	// to the DB. It contains the field Raw (deeply embedded),
	// and also an FSItem that contains an Errer. 
	m5db.ContentityRow
	
	// LogInfo is (the index of the Contentity in 
	// the larger slice) + (the processing stage ID)
	LogInfo
	// logIdx int
	// logStg string
	
	// ParserResults is parseutils.ParserResults_ffs
	// (ffs = file format -specific = "html" or "mkdn" but not
	// "xml" cos Go's XML parser does not produce a tree structure) 
	ParserResults interface{}

	GTokens      []*gtoken.GToken
	GTags        []*gtree.GTag
	*gtree.GTree // maybe not need GRootTag or RootOfASTptr
	GTknsWriter, GTreeWriter,
	GEchoWriter io.Writer
	GLinks

	// GEnts is "ENTITY"" directives (both with "%" and without).
	GEnts map[string]*gparse.GEnt
	// DElms is "ELEMENT" directives.
	DElms map[string]*gtree.GTag

	TagTally StringTally
	AttTally StringTally

	// FU.OutputFiles // This was useful at one point
}

// LogInfo exists mainly to provide a grep'able string:
// for example "(01:4a)", where 01 is the index of the
// Contentity and 4a is the processing stage. This is
// obv a candidate for replacement by stdlib's slog.
//
// The [io.Writer] field W exists outside of the
// [github.com/fbaube/mlog] logging subsystem 
// and should only be used if `mlog` is not.
// .
type LogInfo struct {
	Lindex int
	Lstage string
	W io.Writer
	}

func (p *LogInfo) String() string {
     return fmt.Sprintf("(%02d:%s)", p.Lindex, p.Lstage) 
     }

func (p *Contentity) IsDir() bool {
	// return p.ContentityRow.FSItem.IsDir()
	// return p.FSItem.IsDir()
	return p.FSItem.IsDir()
}

func (p *Contentity) IsDirlike() bool {
	// return p.ContentityRow.FSItem.IsDir()
	return p.FSItem.IsDirlike()
	// return p.Nord.IsDir()
}

type norderCreationState struct {
	// nexSeqID int // reset to 0 when doing another tree ?
	rootPath string
	// summaryString StringFunc
}

// pNCS is the (oops, global) state of Contentity creation.
// FIXME ID assignment should be offloaded to the DB.
// var pNCS *norderCreationState = new(norderCreationState)

// String is developer output. Hafta dump:
// FU.InputFile, FU.OutputFiles, GTree,
// GRefs, *XmlFileMeta, *XmlItems, *DitaInfo
func (p Contentity) String() string {
	var sGTree string
	if p.GTree != nil {
		sGTree = p.GTree.String()
	}
	// s := fmt.Sprintf("[len:%d]", p.Size())
	s := fmt.Sprintf("||%s||GTree|%s||OutbKeyLinks|%+v|KeyLinkTgts|%+v|OutbUriLinks|%+v|UriLinkTgts|%+v||",
		SU.Tildotted(p.FSItem.FPs.AbsFP) /* p.OutputFiles.String(), */, sGTree,
		p.KeyRefncs, p.KeyRefnts, p.UriRefncs, p.UriRefnts)
	/* code to use ?
			if p.XmlFileMeta != nil {
				s += fmt.Sprintf("XmlFileMeta|%s||", p.XmlFileMeta.String())
			}
		* /
		if p.IDinfo != nil {
			s += fmt.Sprintf("xf.IDinfo|%s||", p.IDinfo.String())
		}
	if p.GEnts != nil {
		// FIXME s += fmt.Sprintf("GEnts|%s||", p.GEnts.String())
		 * 	}
		 * 	if p.DElms != nil {
		// FIXME s += fmt.Sprintf("DElms|%s||", p.DElms.String())
	}
	== */
	// if p.DitaInfo != nil {
	s += fmt.Sprintf("DitaInfo|Flav:%s|Cntp:%s|", p.DitaFlavor, p.DitaContype)
	// }
	return s
}

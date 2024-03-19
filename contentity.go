package mcfile

import (
	"fmt"
	"github.com/fbaube/gparse"
	"github.com/fbaube/gtoken"
	"github.com/fbaube/gtree"
	MU "github.com/fbaube/miscutils"
	ON "github.com/fbaube/orderednodes"
	DRM "github.com/fbaube/datarepo/rowmodels"
	SU "github.com/fbaube/stringutils"
	"io"
)

type ContentityStage func(*Contentity) *Contentity

// For the record, ignore the API of
// https://godoc.org/golang.org/x/net/html#Node

// Contentity is awesome.
// .
type Contentity struct { // has has has has Raw
	ON.Nord
	MU.Errer
	// CFU.GCtx // utils/cliflagutils

	LogInfo
	logIdx int
	logStg string
	
	// ContentityRow is what gets persisted to the DB (and has Raw)
	DRM.ContentityRow
	// FU.OutputFiles // This was useful at one point

	// ParserResults is parseutils.ParserResults_ffs
	// (ffs = file format -specific)
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
}

// LogInfo exists mainly to provide a grep'able string -
// for example "(01:4a)". The [io.Writer] exists outside
// of the [github.com/fbaube/mlog] logging subsystem and
// should only be used if `mlog` is not.
// 
type LogInfo struct {
	index int
	stage string
	W io.Writer
	}

func (p *LogInfo) String() string {
     return fmt.Sprintf("(%02d:%s)", p.index, p.stage) 
     }

func (p *Contentity) IsDir() bool {
	return p.ContentityRow.FSItem.IsDir()
}

type norderCreationState struct {
	nexSeqID int // reset to 0 when doing another tree ?
	rootPath string
	// summaryString StringFunc
}

// pNCS is the (oops, global) state of Contentity creation.
// FIXME ID assignment should be offloaded to the DB.
// .
var pNCS *norderCreationState = new(norderCreationState)

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
		SU.Tildotted(p.FSItem.FPs.AbsFP.S()) /* p.OutputFiles.String(), */, sGTree,
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

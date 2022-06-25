package mcfile

import (
	"errors"
	"fmt"
	"io"
	FP "path/filepath"

	FU "github.com/fbaube/fileutils"
	"github.com/fbaube/gparse"
	"github.com/fbaube/gtoken"
	"github.com/fbaube/gtree"
	MU "github.com/fbaube/miscutils"
	L "github.com/fbaube/mlog"
	ON "github.com/fbaube/orderednodes"
	"github.com/fbaube/repo"
	SU "github.com/fbaube/stringutils"
)

type ContentityStage func(*Contentity) *Contentity

// For the record, ignore the API of
// https://godoc.org/golang.org/x/net/html#Node

// Contentity is awesome.
type Contentity struct {
	ON.Nord
	MU.Errer
	// CFU.GCtx // utils/cliflagutils
	logIdx int
	logStg string
	// ContentityRecord is what gets persisted to the DB
	repo.ContentityRecord
	// FU.OutputFiles // NOTE Does this belong here ? Not sure.

	// ParserResults is parseutils.ParserResults_ffs
	ParserResults interface{}
	GTokens       []*gtoken.GToken
	GTags         []*gtree.GTag
	*gtree.GTree  // maybe not need GRootTag or RootOfASTp

	GTokensWriter, GTreeWriter io.Writer
	GLinks
	// GEnts is "ENTITY"" directives (both with "%" and without).
	GEnts map[string]*gparse.GEnt
	// DElms is "ELEMENT" directives.
	DElms map[string]*gtree.GTag

	TagTally StringTally
	AttTally StringTally
}

func (p *Contentity) IsDir() bool {
	return p.ContentityRecord.PathProps.IsOkayDir()
}

// RootContentityNord is available to make assignments to/from root node explicit.
type RootContentityNord Contentity

type norderCreationState struct {
	nexSeqID int // reset to 0 when doing another tree ?
	rootPath string
	// summaryString StringFunc
}

var pNCS *norderCreationState = new(norderCreationState)

// NewRootContentityNord needs aRootPath to be an absolute filepath.
func NewRootContentityNord(aRootPath string) (*Contentity, error) {
	L.L.Info("NewRootContentityNord: %s", aRootPath)
	p := new(Contentity)
	pNCS.rootPath = aRootPath
	pPP, e := FU.NewPathProps(aRootPath)
	if e != nil || pPP == nil {
		return nil, FU.WrapAsPathPropsError(
			e, "NewRootContentityNord (L63)", pPP)
	}
	// This also does content fetching & analysis !
	pCR, e := repo.NewContentityRecord(pPP)
	p.ContentityRecord = *pCR
	if e != nil || pCR == nil {
		L.L.Error("NewRootContentityNord<%s>: %s", aRootPath, e.Error())
		return p, e
	}
	// This block should not happen. And anyways for a
	// directory, we don't need to worry about any error.
	// Unless maybe there's some weird permissions problem.
	/*
		if pCR.HasError() && !pPP.IsOkayDir() {
			println("newRootCty failed:", pCR.GetError().Error())
			pCR.SetError(fmt.Errorf("newRootCty<%s> failed: %w",
				pCR.AbsFP, pCR.GetError()))
			return nil
		}
	*/
	// Now fill in the Contentity, using code taken from NewMCFile(..)
	p.GLinks = *new(GLinks)
	// println("D=> NewContentity:", p.String()) // p.MType, p.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", p.MType, p.AbsFP())
	p.Nord = *ON.NewRootNord(aRootPath, nil)
	// println("NewRootContentityNord:", FU.Tildotted(p.AbsFP()))
	// fmt.Printf("\t RootNord seqID %d \n", p.SeqID())
	return p, nil
}

func NewContentity(aPath string) (*Contentity, error) {
	if aPath == "" {
		return nil, errors.New("newcontentity: missing path")
	}
	p := new(Contentity)
	p.Nord = *ON.NewNord(aPath)
	// fmt.Printf("\t Nord seqID %d \n", p.SeqID())

	var pPP *FU.PathProps
	var e error
	if FP.IsAbs(aPath) {
		pPP, e = FU.NewPathProps(aPath)
	} else {
		if !FP.IsAbs(pNCS.rootPath) {
			e = errors.New("rootPath not absolute: " + pNCS.rootPath)
		} else {
			pPP, e = FU.NewPathPropsRelativeTo(aPath, pNCS.rootPath)
		}
	}
	if e != nil {
		return nil, fmt.Errorf("newcontentity: %w", e)
	}

	if pPP.IsOkayDir() {
		L.L.Info(SU.Ybg(" Directory " + SU.Tildotted(pPP.AbsFP.S())))
		p.ContentityRecord.PathProps = *pPP
		return p, nil
	}
	L.L.Okay(SU.Gbg(" " + pPP.String() + " "))
	// This also does content fetching & analysis !
	pCR, e := repo.NewContentityRecord(pPP)
	if e != nil || pCR == nil {
		// panic("BAD pCR")
		// L.L.Error("New contentity failed")
		// L.L.Error("NewContentity<%s>: %s", pPP.AbsFP, e.Error())
		return nil, fmt.Errorf("newcontentity<%s>: %w", pPP.AbsFP, e)
	}
	// Now fill in the Contentity, using code taken from NewMCFile(..)
	p.ContentityRecord = *pCR
	p.GLinks = *new(GLinks)
	// println("D=> NewContentity:", p.String()) // p.MType, p.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", p.MType, p.AbsFP())
	return p, nil
}

// String is developer output. Hafta dump:
// FU.InputFile, FU.OutputFiles, GTree,
// GRefs, *XmlFileMeta, *XmlItems, *DitaInfo
func (p Contentity) String() string {
	// var BF BigFields = p.PushBigFields()

	var sGTree string
	if p.GTree != nil {
		sGTree = p.GTree.String()
	}
	// s := fmt.Sprintf("[len:%d]", p.Size())
	s := fmt.Sprintf("||%s||GTree|%s||OutbKeyLinks|%+v|KeyLinkTgts|%+v|OutbUriLinks|%+v|UriLinkTgts|%+v||",
		SU.Tildotted(p.PathProps.AbsFP.S()) /* p.OutputFiles.String(), */, sGTree,
		p.OutgoingKeys, p.IncomableKeys, p.OutgoingURIs, p.IncomableURIs)
	/*
			if p.XmlFileMeta != nil {
				s += fmt.Sprintf("XmlFileMeta|%s||", p.XmlFileMeta.String())
			}
		* /
		if p.IDinfo != nil {
			s += fmt.Sprintf("xf.IDinfo|%s||", p.IDinfo.String())
		}
	*/
	/* ==
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

	// p.PopBigFields(BF)
	return s
}

package mcfile

import (
	"fmt"
	"io"

	"github.com/fbaube/db"
	FU "github.com/fbaube/fileutils"
	"github.com/fbaube/gtoken"
	"github.com/fbaube/gtree"
	L "github.com/fbaube/mlog"
	ON "github.com/fbaube/orderednodes"
	SU "github.com/fbaube/stringutils"
)

// For the record, ignore the API of
// https://godoc.org/golang.org/x/net/html#Node

// Contentity is awesome.
type Contentity struct {
	ON.Nord
	logIdx int
	logStg string
	// ContentRecord is what gets persisted to the DB
	db.ContentRecord
	// ParserResults is parseutils.ParserResults_ffs
	ParserResults interface{}
	GTokens       []*gtoken.GToken
	GTags         []*gtree.GTag
	*gtree.GTree  // maybe not need GRootTag or RootOfASTp

	GTokensOutput, GTreeOutput io.Writer
	GLinks

	TagTally StringTally
	AttTally StringTally
}

func (p *Contentity) IsDir() bool {
	return p.ContentRecord.PathProps.IsOkayDir()
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
func NewRootContentityNord(aRootPath string) *Contentity {
	p := new(Contentity)
	pNCS.rootPath = aRootPath
	pPP := FU.NewPathProps(aRootPath)
	if pPP == nil {
		panic("NewRootContentityNord FAILED on pPP")
	}
	// This also does content fetching & analysis !
	pCR := db.NewContentRecord(pPP)
	if pCR == nil {
		panic("NewRootContentityNord FAILED on pCR")
	}
	// This block should not happen. And anyways for a directory,
	// we don't need to worry about any error. Unless maybe there's
	// some weird permissions problem.
	if pCR.GetError() != nil && !pPP.IsOkayDir() {
		println("newRootCty failed:", pCR.GetError().Error())
		pCR.SetError(fmt.Errorf("newRootCty<%s> failed: %w",
			pCR.AbsFP(), pCR.GetError()))
		return nil
	}
	// Now fill in the Contentity, using code taken from NewMCFile(..)
	p.ContentRecord = *pCR
	p.GLinks = *new(GLinks)
	// println("D=> NewContentity:", p.String()) // p.MType, p.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", p.MType, p.AbsFP())
	p.Nord = *ON.NewRootNord(aRootPath, nil)
	// println("NewRootContentityNord:", FU.Tildotted(p.AbsFP()))
	// fmt.Printf("\t RootNord seqID %d \n", p.SeqID())
	return p
}

func NewContentity(aPath string) *Contentity {
	if aPath == "" {
		println("NewContentity: missing path")
		return nil
	}
	p := new(Contentity)
	p.Nord = *ON.NewNord(aPath)
	// fmt.Printf("\t Nord seqID %d \n", p.SeqID())

	pPP := FU.NewPathPropsRelativeTo(aPath, pNCS.rootPath)
	if pPP.IsOkayDir() {
		L.L.Info(SU.Ybg(" Directory " + FU.Tildotted(pPP.AbsFP())))
		p.ContentRecord.PathProps = *pPP
		return p
	}
	L.L.Okay(SU.Gbg(" " + pPP.String() + " "))
	// This also does content fetching & analysis !
	pCR := db.NewContentRecord(pPP)
	if pCR.GetError() != nil {
		pCR.SetError(fmt.Errorf("newCty<%s> failed: %w",
			pCR.AbsFilePath, pCR.GetError()))
		return nil
	}
	// Now fill in the Contentity, using code taken from NewMCFile(..)
	p.ContentRecord = *pCR
	p.GLinks = *new(GLinks)
	// println("D=> NewContentity:", p.String()) // p.MType, p.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", p.MType, p.AbsFP())
	return p
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
		FU.Tildotted(p.PathProps.AbsFP()) /* p.OutputFiles.String(), */, sGTree,
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
	s += fmt.Sprintf("DitaInfo|ML:%s|Cntp:%s|", p.DitaMarkupLg, p.DitaContype)
	// }

	// p.PopBigFields(BF)
	return s
}

package mcfile

import (
	"errors"
	"fmt"
	"io"
	FP "path/filepath"

	DU "github.com/fbaube/dbutils"
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
	// ContentityRecord is what gets persisted to the DB
	DU.ContentityRecord
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
	pCR := DU.NewContentityRecord(pPP)
	if pCR == nil {
		panic("NewRootContentityNord FAILED on pCR")
	}
	// This block should not happen. And anyways for a directory,
	// we don't need to worry about any error. Unless maybe there's
	// some weird permissions problem.
	/*
		if pCR.HasError() && !pPP.IsOkayDir() {
			println("newRootCty failed:", pCR.GetError().Error())
			pCR.SetError(fmt.Errorf("newRootCty<%s> failed: %w",
				pCR.AbsFP, pCR.GetError()))
			return nil
		}
	*/
	// Now fill in the Contentity, using code taken from NewMCFile(..)
	p.ContentityRecord = *pCR
	p.GLinks = *new(GLinks)
	// println("D=> NewContentity:", p.String()) // p.MType, p.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", p.MType, p.AbsFP())
	p.Nord = *ON.NewRootNord(aRootPath, nil)
	// println("NewRootContentityNord:", FU.Tildotted(p.AbsFP()))
	// fmt.Printf("\t RootNord seqID %d \n", p.SeqID())
	return p, nil
}

func NewContentity(aPath string) *Contentity {
	if aPath == "" {
		println("NewContentity: missing path")
		return nil
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
		p.SetError(fmt.Errorf("NewContentity: %w", e))
		return p
	}

	if pPP.IsOkayDir() {
		L.L.Info(SU.Ybg(" Directory " + SU.Tildotted(pPP.AbsFP.S())))
		p.ContentityRecord.PathProps = *pPP
		return p
	}
	L.L.Okay(SU.Gbg(" " + pPP.String() + " "))
	// This also does content fetching & analysis !
	pCR := DU.NewContentityRecord(pPP)
	if pCR == nil {
		// panic("BAD pCR")
		// L.L.Error("New contentity failed")
		return nil
	}
	if pCR.HasError() {
		pCR.SetError(fmt.Errorf("newCty<%s> failed: %w",
			pCR.AbsFP, pCR.GetError()))
		return p // nil
	}
	// Now fill in the Contentity, using code taken from NewMCFile(..)
	p.ContentityRecord = *pCR
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

/*
WHATWG:

document . images
Returns an HTMLCollection of the img elements in the Document.
document . embeds
document . plugins
Return an HTMLCollection of the embed elements in the Document.
document . links
Returns an HTMLCollection of the a and area elements in the Document that have href attributes.
document . forms
Return an HTMLCollection of the form elements in the Document.
document . scripts
Return an HTMLCollection of the script elements in the Document.

element . innerText [ = value ]
Returns the element's text content "as rendered".
Can be set, to replace the element's children with the given
value, but with line breaks converted to br elements.


element . dataset
https://html.spec.whatwg.org/#domstringmap
Returns a DOMStringMap object for the element's data-* attributes.
Hyphenated names become camel-cased. For example, data-foo-bar=""
becomes element.dataset.fooBar.
*/

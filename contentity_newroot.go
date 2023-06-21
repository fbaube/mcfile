package mcfile

import (
	"errors"
	"fmt"
	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	ON "github.com/fbaube/orderednodes"
	"github.com/fbaube/repo/sqlite"
	RM "github.com/fbaube/rowmodels"
	FP "path/filepath"
)

// RootContentity makes assignments
// to/from root node explicit.
type RootContentity Contentity

// NewRootContentity returns a RootContentity Nord (i.e. node with
// ordered children) that can be the root of a new Contentity tree.
// It requires that argument aRootPath is an absolute filepath AND
// is NOT a directory.
// .
func NewRootContentity(aRootPath string) (*RootContentity, error) {
	L.L.Info("NewRootContentity: %s", aRootPath)
	if aRootPath == "" {
		return nil, errors.New("newrootcontentity: missing path")
	}
	if !FP.IsAbs(aRootPath) {
		return nil, errors.New(
			"NewRootContentity: not an abs.fp: " + aRootPath)
	}
	var pNewCty *RootContentity
	pNewCty = new(RootContentity)
	// Global assignment (oops)
	pNCS.rootPath = aRootPath
	if pNCS.nexSeqID != 0 {
		L.L.Warning("New root cty: seq ID is: %d", pNCS.nexSeqID)
	}

	// ========================
	//  Start with a PathProps
	// ========================
	var pPP *FU.PathProps
	var e error
	pPP, e = FU.NewPathProps(aRootPath)
	if e != nil || pPP == nil {
		return nil, FU.WrapAsPathPropsError(
			e, "NewRootContentity (L47,path=>PP)", pPP)
	}
	// =============================
	//  "Promote" to a PathAnalysis
	// =============================
	pPA, e := FU.NewPathAnalysis(pPP)
	if e != nil || pPA == nil {
		L.L.Error("NewRootContentity(PP=>PA)<%s>: %s", aRootPath, e)
		return nil, fmt.Errorf(
			"NewRootContentity(PP=>PA)<%s>: %w", aRootPath, e)
	}
	// =================================
	//  "Promote" to a ContentityRecord
	// =================================
	var pCR *RM.ContentityRow
	pCR, e = sqlite.NewContentityRow(pPP, pPA)
	if e != nil || pCR == nil {
		L.L.Error("NewRootContentity(PA=>CR)<%s>: %s", aRootPath, e)
		return nil, fmt.Errorf(
			"NewRootContentity(PA=>CR)<%s>: %w", aRootPath, e)
	}
	// L.L.Warning("NewRootCty (PP) %+v", pCR.PathProps)
	// nil! L.L.Warning("NewRootCty (PA) %+v", *pCR.PathAnalysis)
	pNewCty.ContentityRow = *pCR
	// ==================================
	//  Now fill in the ContentityRecord
	// ==================================
	pNewCty.GLinks = *new(GLinks)
	// println("D=> NewContentity:", pNewCty.String()) // pNewCty.MType, pNewCty.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", pNewCty.MType, pNewCty.AbsFP())
	pNewCty.Nord = *ON.NewRootNord(aRootPath, nil)
	// println("NewRootContentity:", FU.Tildotted(pNewCty.AbsFP()))
	// fmt.Printf("\t RootNord seqID %d \n", pNewCty.SeqID())
	// var pRC *RootContentity
	// var C Contentity
	// var R RootContentity
	// pRC = pNewCty
	// C = *pNewCty
	// R = RootContentity(C)
	if pNewCty == nil {
		panic("nil pNewCty")
	}
	return pNewCty, nil
}

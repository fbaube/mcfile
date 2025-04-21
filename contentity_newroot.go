package mcfile

import (
	"errors"
	"io/fs"
	"fmt"
	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	ON "github.com/fbaube/orderednodes"
	"github.com/fbaube/m5db"
	FP "path/filepath"
)

// RootContentity makes assignments
// to/from root node explicit.
type RootContentity Contentity

// NewRootContentity returns a RootContentity Nord (i.e. node with
// ordered children) that can be the root of a new Contentity tree.
// It requires that argument aRootPath is an absolute filepath and
// is a directory.
// .
func NewRootContentity(aRootPath string) (*RootContentity, error) {
	L.L.Info("NewRootContentity: %s", aRootPath)
	if aRootPath == "" {
                return nil, &fs.PathError{Op:"NewRootContentity",
                       Err:errors.New("Missing path"),Path:"(nil)"}
        }
	aRootPath = FU.EnsureTrailingPathSep(aRootPath)
	if !FP.IsAbs(aRootPath) {
		return nil, &fs.PathError{Op:"NewRootContentity",
                       Err:errors.New("Not an absolute filepath"),
		       Path:aRootPath}
	}
	if !FU.IsDirAndExists(aRootPath) {
	   	return nil, &fs.PathError{Op:"NewRootContentity",
                       Err:errors.New("Not a directory"),Path:aRootPath}
	}
	// L.L.Info("CHECK-1")
	var pNewCty *RootContentity
	pNewCty = new(RootContentity)
	// Global assignment (oops)
	CntyEng.rootPath = aRootPath
	/* if CntyEng.nexSeqID != 0 {
		L.L.Warning("New root cty: seq ID is: %d", CntyEng.nexSeqID)
	} */
	// ======================
	//  Start with an FSItem
	// ======================
	// L.L.Info("CHECK-2")
	var pFSI *FU.FSItem
	pFSI = FU.NewFSItem(aRootPath)
	// L.L.Debug("pFSI %p *pFSI %T e %T", pFSI, *pFSI, e)
	if pFSI.HasError() {
	   // L.L.Info("CHECK-2b")
	   return nil, &fs.PathError{Op:"newrootfsitem",
	   	  Err:pFSI.GetError(),Path:aRootPath}
	}
	var e error
	/*
	SKIP this part 
	// =============================
	//  "Promote" to a PathAnalysis
	// =============================
	// NewPathAnalysis returns (nil,nil) for DIRLIKE 
	pPA, e := CA.NewPathAnalysis(pFSI)
	if e != nil { // || pPA == nil {
		L.L.Error("NewRootContentity(PP=>PA)<%s>: %s", aRootPath, e)
		return nil, fmt.Errorf(
			"NewRootContentity(PP=>PA)<%s>: %w", aRootPath, e)
	}
	*/
	// =================================
	//  "Promote" to a ContentityRecord
	// =================================
	var pCR *m5db.ContentityRow
	pCR, e = m5db.NewContentityRow(pFSI, nil)
	if e != nil || pCR == nil {
		L.L.Error("NewRootContentity(PA=>CR)<%s>: %s", aRootPath, e)
		return nil, fmt.Errorf(
			"NewRootContentity(PA=>CR)<%s>: %w", aRootPath, e)
	}
	// L.L.Warning("NewRootCty (PP) %+v", pCR.FSItem)
	// nil! L.L.Warning("NewRootCty (PA) %+v", *pCR.PathAnalysis)
	pCR.FSItem = *pFSI
	pNewCty.ContentityRow = *pCR
	// ==================================
	//  Now fill in the ContentityRecord
	// ==================================
	// pNewCty.GLinks = *new(GLinks)
	// println("D=> NewContentity:", pNewCty.String()) // pNewCty.MType, pNewCty.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", pNewCty.MType, pNewCty.AbsFP())

	pNewCty.Nord = *ON.NewRootNord(aRootPath, nil)
	return pNewCty, nil
}

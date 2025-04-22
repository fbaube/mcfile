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
func NewRootContentity(aRootPath string) *RootContentity {
     	var Empty *RootContentity
	Empty = new(RootContentity)
	L.L.Info("NewRootContentity: %s", aRootPath)
	if aRootPath == "" {
                Empty.SetError(&fs.PathError{Op:"NewRootContentity",
                       Err:errors.New("Missing path"),Path:"(nil)"})
		return Empty
        }
	aRootPath = FU.EnsureTrailingPathSep(aRootPath)
	if !FP.IsAbs(aRootPath) {
		Empty.SetError(&fs.PathError{Op:"NewRootContentity",
                       Err:errors.New("Not an absolute filepath"),
		       Path:aRootPath})
		return Empty
	}
	if !FU.IsDirAndExists(aRootPath) {
	   	Empty.SetError(&fs.PathError{Op:"NewRootContentity",
                       Err:errors.New("Not a directory"),Path:aRootPath})
		return Empty
	}
	// L.L.Info("CHECK-1")
	var pNewCty *RootContentity
	pNewCty = new(RootContentity)
	// Global assignment (oops)
	CntyEng.rootPath = aRootPath
	// =================================
	//  "Promote" to a ContentityRecord
	// =================================
	var pCR *m5db.ContentityRow
	var pFSI *FU.FSItem
	var e error
	pFSI = FU.NewFSItem(aRootPath)
	if pFSI.HasError() {
                println("NewRootContentity: NewFSI error at LINE 56")
		Empty.SetError(pFSI.GetError())
                return Empty
        }
	pCR = m5db.NewContentityRow(pFSI, nil)
	pCR.FSItem = *pFSI
	pNewCty.ContentityRow = *pCR
	if pNewCty.HasError() {
		L.L.Error("NewRootContentity(PA=>CR)<%s>: %s", aRootPath, e)
		pNewCty.SetError(fmt.Errorf(
			"NewRootContentity(PA=>CR)<%s>: %w", aRootPath, e))
		return pNewCty
	}
	// L.L.Warning("NewRootCty (PP) %+v", pCR.FSItem)
	// nil! L.L.Warning("NewRootCty (PA) %+v", *pCR.PathAnalysis)
	// ==================================
	//  Now fill in the ContentityRecord
	// ==================================
	// pNewCty.GLinks = *new(GLinks)
	// println("D=> NewContentity:", pNewCty.String()) // pNewCty.MType, pNewCty.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", pNewCty.MType, pNewCty.AbsFP())

	pNewCty.Nord = *ON.NewRootNord(aRootPath, nil)
	return pNewCty 
}

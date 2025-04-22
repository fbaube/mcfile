package mcfile

import (
	"errors"
	"io/fs"
	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	ON "github.com/fbaube/orderednodes"
	"github.com/fbaube/m5db"
	SU "github.com/fbaube/stringutils"
	CA "github.com/fbaube/contentanalysis"
	FP "path/filepath"
)

// NewContentity returns a Contentity Nord (i.e. a node with content a
// and ordered children) that can NOT be the root of a Contentity tree.
// If there is an error, it is returned in the embedded Errer.
//
// It should accept either an absolute or a relative filepath, altho
// relative is preferred, for various reasons, mainly because of the
// preferences of the path and filepath stdlibs.
// 
// We want everything to be in a nice tree of Nords, and it means that
// we have to create Contenties for directories too (where `Raw_type
// == SU.Raw_type_DIRLIKE`), so we have to handle that case too. 
// .
func NewContentity(aPath string) *Contentity {
     	var Empty *Contentity
	Empty = new(Contentity)
	if aPath == "" {
	   	Empty.SetError(&fs.PathError{Op:"NewContentity",
		       Err:errors.New("missing path"),Path:"(nil)"})
		return Empty
	}
	pFPs, e := FU.NewFilepaths(aPath)
	if e != nil {
	   Empty.SetError(&fs.PathError{ Op:"NewFilepaths", Path:aPath, Err:e })
	   return Empty
	}
	L.L.Debug("NewContentity.FPs: %s", pFPs.String())

	// ===============================
	//  The path is OK, so from here
	//  on we can return the pNewCnty.
	//  Declare some useful vars.
	// ===============================
	var pC       *Contentity
	var pCR *m5db.ContentityRow
	var pFSI  *FU.FSItem
	var pPA   *CA.PathAnalysis
	
	pC = new(Contentity)
	pC.Nord = *ON.NewNord(aPath)
	// fmt.Printf("\t Nord seqID %d \n", p.SeqID())
	L.L.Okay(SU.Ybg("===> New Contentity: %s"), SU.Tildotted(aPath))
	
	// ===============================
	//  Start by creating an FSItem; 
	//  bail if it gets/has an error
	// (altho this should not happen,
	//  cos we know the path is okay.) 
	// ===============================
	// If we were passed an Abs.FP, it's okay.
	if FP.IsAbs(aPath) {
		pFSI = FU.NewFSItem(aPath)
	} else {
		ppp := FP.Join(CntyEng.rootPath, aPath)
		pFSI = FU.NewFSItem(ppp)
	}
	pC.FSItem = *pFSI
	if pFSI.HasError() {
	   	println("NewContentity: NewFSI error at LINE 72")
		return pC
	}
	// ======================================
	//  If it's a directory (or similar,
	//  such as symlink) we handle it here
	//  cos we don't need to do PathAnalysis. 
	// ======================================
	if pFSI.IsDirlike() {
	   	// This should fail only if the item does not exist.
		pCR = m5db.NewContentityRow(pFSI, nil)
		if pCR.HasError() {
			L.L.Error("NewContentity(Dirlike)<%s>: %s", aPath, e)
			println("LINE 85")
			pFSI.SetError(&fs.PathError{Op:"newctyrow.(dirlike)",
			       Err:pCR.GetError(),Path:aPath})
			return pC
		}
		L.L.Info(SU.Ybg(" Dir " + SU.Tildotted(pFSI.FPs.AbsFP)))
                pCR.FSItem = *pFSI
		pC.ContentityRow = *pCR
		return pC 
        }
	e = pFSI.LoadContents()
	// L.L.Warning("LENGTH %d", len(pFSI.TypedRaw.Raw))
	if e != nil {
   	   println("LINE 98")
	   pC.SetError(&fs.PathError{Op:"FSI.GoGetFileContents",
	       	  Err:e,Path:CntyEng.rootPath})
	   return pC
	}
	// =============================
	//  "Promote" to a PathAnalysis
	// =============================
	// NewPathAnalysis return (nil,nil) for DIRLIKE 
	pPA, e = CA.NewPathAnalysis(pFSI)
	if e != nil { 
	   L.L.Error("NewContentity(PP=>PA)<%s>: %s", aPath, e)
	   println("LINE 110")
	   pC.SetError(&fs.PathError{Op:"FSI.NewPathAnalysis.(PP=>PA)",
	       Err:e,Path:aPath})
	   return pC
	}
	if pPA == nil { panic("WTF") }
	// =================================
	//  "Promote" to a ContentityRecord
	// =================================
	pCR = m5db.NewContentityRow(pFSI, pPA)
	if pCR.HasError() {
		L.L.Error("NewContentity(PA=>CR)<%s>: %s", aPath, e)
	   	println("LINE 122")
		pCR.SetError(&fs.PathError{Op:"FSI.NewContentityRow.(PP=>PA)",
                       Err:pCR.GetError(),Path:aPath})
		return pC
	}
	if pCR.RawType() == "" { // or SU.MU_type_UNK {
		panic("UNK MarkupType in NewContentity")
	}
	// NOW if we want to exit, we can
	// do the necessary assignments
	pC.ContentityRow = *pCR
	if pFSI.IsDirlike() {
		L.L.Info(SU.Ybg(" Directory " + SU.Tildotted(pFSI.FPs.AbsFP)))
		pC.ContentityRow.FSItem = *pFSI
		return pC 
	}
	L.L.Info(SU.Gbg(" " + pFSI.String() + " "))

	// ==================================
	//  Now fill in the ContentityRecord
	// ==================================
	pC.GLinks = *new(GLinks)
	// println("D=> NewContentity:", p.String()) // p.MType, p.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", p.MType, p.AbsFP())
	return pC 
}

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
//
// It should accept either an absolute or a relative filepath.
// 
// NOTE that because of interface hassles, BOTH return values might
// be non-nil, in which case, ignore the error. 
//
// We want everything to be in a nice tree of Nords, and it means that
// we have to create Contenties for directories too (where `Raw_type
// == SU.Raw_type_DIRLIKE`), so we have to handle that case too. 
// .
func NewContentity(aPath string) (*Contentity, error) {
	if aPath == "" {
		return nil, &fs.PathError{Op:"NewContentity",
		       Err:errors.New("Missing path"),Path:"(nil)"}
	}
	pFPs, e := FU.NewFilepaths(aPath)
	if e != nil {
	   return nil, &fs.PathError{ Op:"NewFilepaths", Path:aPath, Err:e }
	   }
	L.L.Debug("NewContentity.FPs: %s", pFPs.String())
	var pNewCnty  *Contentity
	pNewCnty = new(Contentity)
	pNewCnty.Nord = *ON.NewNord(aPath)

	// fmt.Printf("\t Nord seqID %d \n", p.SeqID())
	L.L.Okay(SU.Ybg("===> New Contentity: %s"), SU.Tildotted(aPath))
	
	// ======================
	//  Start with an FSItem
	// ======================
	var pFSI *FU.FSItem
	// If we were passed an Abs.FP, it's okay.
	if FP.IsAbs(aPath) {
		pFSI = FU.NewFSItem(aPath)
	/*
	} else if !FP.IsAbs(CntyEng.rootPath) {
	// Else if the 
		e = &fs.PathError{Op:"NewContentity.IsAbs",
		  Err:errors.New("NewContentity: both input path and " +
		  "Contentity Engine's RootPath are relative"),
		  Path:CntyEng.rootPath}
	*/
	} else {
		// pFSI, e = FU.NewFSItemRelativeTo(aPath, CntyEng.rootPath)
		ppp := FP.Join(CntyEng.rootPath, aPath)
		pFSI = FU.NewFSItem(ppp)
	}
	if pFSI.HasError() { 
	   	println("NewContentity: error at LINE 65")
		return nil, pFSI.GetError() // &fs.PathError{Op:"Path-analysis",
		//       Err:e,Path:CntyEng.rootPath}
	}
	// =========================
	//  If it's a directory (or
	//  similar, such as symlink) 
	// ==========================
	if pFSI.IsDirlike() {
		var pCR *m5db.ContentityRow
		pCR, e = m5db.NewContentityRow(pFSI, nil)
		if e != nil || pCR == nil {
			L.L.Error("NewContentity(Dirlike)<%s>: %s", aPath, e)
			println("LINE 78")
			return nil, &fs.PathError{Op:"FSI.NewCtyRow.(dirlike)",
			       Err:e,Path:aPath}
		}
		L.L.Info(SU.Ybg(" Dir " + SU.Tildotted(pFSI.FPs.AbsFP)))
                pCR.FSItem = *pFSI
		pNewCnty.ContentityRow = *pCR
		return pNewCnty, nil
        }
	var pPA *CA.PathAnalysis
	e = pFSI.LoadContents()
	// L.L.Warning("LENGTH %d", len(pFSI.TypedRaw.Raw))
	if e != nil {
   	   println("LINE 91")	
	   return nil, &fs.PathError{Op:"FSI.GoGetFileContents",
	       	  Err:e,Path:CntyEng.rootPath}
	}
	// =============================
	//  "Promote" to a PathAnalysis
	// =============================
	// NewPathAnalysis return (nil,nil) for DIRLIKE 
	pPA, e = CA.NewPathAnalysis(pFSI)
	if e != nil { // || pPA == nil {
	   L.L.Error("NewContentity(PP=>PA)<%s>: %s", aPath, e)
	   println("LINE 102")
	   return nil, &fs.PathError{Op:"FSI.NewPathAnalysis.(PP=>PA)",
	       Err:e,Path:aPath}
	}
	if pPA == nil { panic("WTF") }
	/*
	if pPA != nil && pPA.MarkupTypeOfMType() == SU.MU_type_UNK {
	   	L.L.Panic("UNK MarkupType in NewContentity L121 (%s)",
			pFSI.FPs.AbsFP)
	}
	*/
	// =================================
	//  "Promote" to a ContentityRecord
	// =================================
	var pCR *m5db.ContentityRow
	pCR, e = m5db.NewContentityRow(pFSI, pPA)
	if e != nil || pCR == nil {
		L.L.Error("NewContentity(PA=>CR)<%s>: %s", aPath, e)
	   	println("LINE 120")
		return nil, &fs.PathError{Op:"FSI.NewContentityRow.(PP=>PA)",
                       Err:e,Path:aPath}
	}
	if pCR.RawType() == "" { // or SU.MU_type_UNK {
		panic("UNK MarkupType in NewContentity")
	}
	// NOW if we want to exit, we can
	// do the necessary assignments
	pNewCnty.ContentityRow = *pCR
	if pFSI.IsDirlike() {
		L.L.Info(SU.Ybg(" Directory " + SU.Tildotted(pFSI.FPs.AbsFP)))
		pNewCnty.ContentityRow.FSItem = *pFSI
		return pNewCnty, nil
	}
	L.L.Info(SU.Gbg(" " + pFSI.String() + " "))

	// ==================================
	//  Now fill in the ContentityRecord
	// ==================================
	pNewCnty.GLinks = *new(GLinks)
	// println("D=> NewContentity:", p.String()) // p.MType, p.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", p.MType, p.AbsFP())
	return pNewCnty, nil
}

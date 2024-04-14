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

// NewContentity returns a Contentity Nord (i.e. node with
// ordered children) that can NOT be the root of a Contentity tree.
//
// NOTE: because of interface hassles, BOTH return values might
// be non-nil, in which case, ignore the error. 
//
// We want everything to be in a nice tree of Nords, and that
// means that we have to create Contenties for directories too
// (where MarkupType == SU.MU_type_DIRLIKE). 
//
// When this func is called while walking a DIRECTORY given
// on the command line, aPath is a simple file (or dir) name,
// with no path separators.
//
// When this func is called for a FILE given on the command line,
// aPath can be either absolute or relative, depending on what was
// on the CLI (altho probably a relFP has been upgraded to an absFP).
//
// Alternative hack to achieve a similar end:
// if pPP,e := NewPP(path); e == nil; pPA,e := new PA(pPP);
// e == nil; pCR,e := NewCR(pPA); e == nil { ... }
// .
func NewContentity(aPath string) (*Contentity, error) {
	if aPath == "" {
		return nil, &fs.PathError{Op:"NewContentity",
		       Err:errors.New("Missing path"),Path:"(nil)"}
	}
	var pNewCnty *Contentity
	pNewCnty = new(Contentity)
	pNewCnty.Nord = *ON.NewNord(aPath)

	// fmt.Printf("\t Nord seqID %d \n", p.SeqID())
	// Try black-on-cyan, cos even for white text, the blue is too dark
	if true {
		// // L.L.Info(SU.Cyanbg("\n  ===> crnt RootPath: %s \n"+
		L.L.Okay(SU.Ybg( // "\n===> Crnt Root Path: %s \n"+
			"===> New Contentity: %s"),
			// SU.Tildotted(CntyEng.rootPath),
			SU.Tildotted(aPath))
	} else {
		L.L.Info( /* SU.Wfg( */ SU.Cyanbg( // Blubg(
			"\n  ===> New Contentity: %s <===           "),
			SU.Tildotted(aPath))
	}
	// ======================
	//  Start with an FSItem
	// ======================
	var pFSI *FU.FSItem
	var e error
	// If we were passed an Abs.FP, it's okay.
	if FP.IsAbs(aPath) {
		pFSI, e = FU.NewFSItem(aPath)
	} else if !FP.IsAbs(CntyEng.rootPath) {
	// Else if the 
		e = &fs.PathError{Op:"NewContentity.IsAbs",
		  Err:errors.New("rRootPath is not absolute"),Path:CntyEng.rootPath}
	} else {
		// pFSI, e = FU.NewFSItemRelativeTo(aPath, CntyEng.rootPath)
		ppp := FP.Join(CntyEng.rootPath, aPath)
		pFSI, e = FU.NewFSItem(ppp)
	}
	if pFSI == nil { // e != nil {
	   	println("LINE 75")
		return nil, &fs.PathError{Op:"Path-analysis",
		       Err:e,Path:CntyEng.rootPath}
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
			println("LINE 89")
			return nil, &fs.PathError{Op:"FSI.NewCtyRow.(dirlike)",
			       Err:e,Path:aPath}
		}
		L.L.Info(SU.Ybg(" Dir " + SU.Tildotted(pFSI.FPs.AbsFP.S())))
                pCR.FSItem = *pFSI
		pNewCnty.ContentityRow = *pCR
		return pNewCnty, nil
        }
	var pPA *CA.PathAnalysis
	e = pFSI.GoGetFileContents()
	// L.L.Warning("LENGTH %d", len(pFSI.TypedRaw.Raw))
	if e != nil {
   	   println("LINE 105")	
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
	   println("LINE 116")
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
	   	println("LINE 131")
		return nil, &fs.PathError{Op:"FSI.NewContentityRow.(PP=>PA)",
                       Err:e,Path:aPath}
	}
	if pCR.MarkupTypeOfMType() == SU.MU_type_UNK {
		panic("UNK MarkupType in NewContentity")
	}
	// NOW if we want to exit, we can
	// do the necessary assignments
	pNewCnty.ContentityRow = *pCR
	if pFSI.IsDirlike() {
		L.L.Info(SU.Ybg(" Directory " + SU.Tildotted(pFSI.FPs.AbsFP.S())))
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

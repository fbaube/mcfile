package mcfile

import (
	"errors"
	"fmt"
	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	ON "github.com/fbaube/orderednodes"
	"github.com/fbaube/repo/sqlite"
	RU "github.com/fbaube/repo/util"
	SU "github.com/fbaube/stringutils"
	FP "path/filepath"
)

// NewContentity returns a Contentity Nord (i.e. node with
// ordered children) that can NOT be the root of a Contentity tree.
//
// We want everything to be in a nice tree of Nords, and that
// means that we have to create Contenties for directories too.
//
// When this func is called while walking a DIRECTORY given
// on the command line, aPath is a simple file (or dir) name,
// with no path separators.
//
// When this func is called for a FILE given on the command line,
// aPath can be either absolute or relative, depending on what was
// on the CLI.
//
// Alternative hack to achieve a similar end:
// if pPP,e := NewPP(path); e == nil; pPA,e := new PA(pPP);
// e == nil; pCR,e := NewCR(pPA); e == nil { ... }
// .
func NewContentity(aPath string) (*Contentity, error) {
	if aPath == "" {
		return nil, errors.New("newcontentity: missing path")
	}
	var pNewCnty *Contentity
	pNewCnty = new(Contentity)
	pNewCnty.Nord = *ON.NewNord(aPath)
	// fmt.Printf("\t Nord seqID %d \n", p.SeqID())
	// Try black-on-cyan, cos even for white text, the blue is too dark
	if true {
		L.L.Info(SU.Cyanbg("\n\t ===> crnt RootPath: %s \n"+
			"\t ===> New Contentity: %s"),
			pNCS.rootPath, SU.ElideHomeDir(aPath))
	} else {
		L.L.Info( /* SU.Wfg( */ SU.Cyanbg( // Blubg(
			"\n\t ===> New Contentity: %s <===           "),
			SU.ElideHomeDir(aPath))
	}
	// ========================
	//  Start with a PathProps
	// ========================
	var pPP *FU.PathProps
	var e error
	if FP.IsAbs(aPath) {
		pPP, e = FU.NewPathProps(aPath)
	} else if !FP.IsAbs(pNCS.rootPath) {
		e = errors.New("rootPath not absolute: " + pNCS.rootPath)
	} else {
		pPP, e = FU.NewPathPropsRelativeTo(aPath, pNCS.rootPath)
	}
	if e != nil {
		return nil, fmt.Errorf("newcontentity: %w", e)
	}
	// e = pPP.FetchRaw()
	e = pPP.GoGetFileContents()
	if e != nil {
		return nil, fmt.Errorf("newcontentity: %w", e)
	}
	// =============================
	//  "Promote" to a PathAnalysis
	// =============================
	pPA, e := FU.NewPathAnalysis(pPP)
	if e != nil || pPA == nil {
		L.L.Error("NewContentity(PP=>PA)<%s>: %s", aPath, e)
		return nil, fmt.Errorf(
			"NewContentity(PP=>PA)<%s>: %w", aPath, e)
	}
	// =================================
	//  "Promote" to a ContentityRecord
	// =================================
	var pCR *RU.ContentityRecord
	pCR, e = sqlite.NewContentityRecord(pPP, pPA)
	if e != nil || pCR == nil {
		L.L.Error("NewContentity(PA=>CR)<%s>: %s", aPath, e)
		return nil, fmt.Errorf(
			"NewContentity(PA=>CR)<%s>: %w", aPath, e)
	}
	// NOW if we want to exit, we can
	// do the necessary assignments
	pNewCnty.ContentityRecord = *pCR
	if pPP.IsOkayDir() {
		L.L.Info(SU.Ybg(" Directory " + SU.ElideHomeDir(pPP.AbsFP.S())))
		pNewCnty.ContentityRecord.PathProps = *pPP
		return pNewCnty, nil
	}
	L.L.Okay(SU.Gbg(" " + pPP.String() + " "))

	// ==================================
	//  Now fill in the ContentityRecord
	// ==================================
	pNewCnty.GLinks = *new(GLinks)
	// println("D=> NewContentity:", p.String()) // p.MType, p.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", p.MType, p.AbsFP())
	return pNewCnty, nil
}

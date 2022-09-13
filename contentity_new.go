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
// Alternative hack to achieve a similar end:
// if pPP,e := NewPP(path); e == nil; pPA,e := new PA(pPP);
// e == nil; pCR,e := NewCR(pPA); e == nil { ... }
// .
func NewContentity(aPath string) (*Contentity, error) {
	if aPath == "" {
		return nil, errors.New("newcontentity: missing path")
	}
	var pNewCty *Contentity
	pNewCty = new(Contentity)
	pNewCty.Nord = *ON.NewNord(aPath)
	// fmt.Printf("\t Nord seqID %d \n", p.SeqID())

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
	// =============================
	//  "Promote" to a PathAnalysis
	// =============================
	pPA, e := FU.NewPathAnalysis(pPP)
	if e != nil || pPA == nil {
		L.L.Error("NewContentity(PP=>PA)<%s>: %s", aPath, e)
		return nil, fmt.Errorf(
			"NewContentity(PP=>PA)<%s>: %w", aPath, e)
	}
	pPA.PathProps = pPP
	// =================================
	//  "Promote" to a ContentityRecord
	// =================================
	var pCR *RU.ContentityRecord
	pCR, e = sqlite.NewContentityRecord(pPA)
	if e != nil || pCR == nil {
		L.L.Error("NewContentity(PA=>CR)<%s>: %s", aPath, e)
		return nil, fmt.Errorf(
			"NewContentity(PA=>CR)<%s>: %w", aPath, e)
	}
	// NOW if we want to exit, we can
	// do the necessary assignments
	pNewCty.ContentityRecord = *pCR
	if pPP.IsOkayDir() {
		L.L.Info(SU.Ybg(" Directory " + SU.ElideHomeDir(pPP.AbsFP.S())))
		pNewCty.ContentityRecord.PathProps = pPP
		return pNewCty, nil
	}
	L.L.Okay(SU.Gbg(" " + pPP.String() + " "))

	// ==================================
	//  Now fill in the ContentityRecord
	// ==================================
	pNewCty.GLinks = *new(GLinks)
	// println("D=> NewContentity:", p.String()) // p.MType, p.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", p.MType, p.AbsFP())
	return pNewCty, nil
}

package mcfile

import (
	"io/fs"
	"os"
	"errors"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	L "github.com/fbaube/mlog"
)

// CntyFS is a global, which is a mistake. 
// >> var CntyFS *ContentityFS

// NewContentityFS proceeds as follows:
//  - initialize
//  - create an [os.DirFS]
//  - FIXME: an os.Root
//  - walk the DirFS, creating Contentities
//    and appending them to a slice
//  - process the list to identify and make
//    parent/child links
//
// The path argument should probably be an absolute filepath, 
// because a relative filepath might cause problems.
//
// It uses the global [CntyFS], which precludes
// re-entrancy and concurrency.
//
// Note that when we use [os.DirFS], it appears
// to make no difference whether path 
//  - is relative or absolute
//  - ends with a trailing slash or not
//  - is a directory or a symlink to a directory
// .
func NewContentityFS(aPath string, okayFilexts []string) (*ContentityFS, error){
     	var CntyFS *ContentityFS
	var e error 
	var pFPs *FU.Filepaths
	
	// --------------------
	//  Prepare filepath(s)
	// --------------------
	pFPs, e = FU.NewFilepaths(aPath)
	if e != nil {
	     	// L.L.Error("NewCntyFS: bad path: %s", aPath)
		return nil, &os.PathError { Path:aPath, Err:e,
		       Op:"newcntyfs: bad root path: " + e.Error() }
	}
	pathToUse := FU.EnsureTrailingPathSep(aPath)
	// os.DirFS(..) does not check or report problems
	// with the path argument, so we DIY here 
	if !FU.IsDirAndExists(pathToUse) {
		// L.L.Error("NewCntyFS: Not a directory: %s", path)
		return nil, &os.PathError { Path:aPath, Err:errors.New(
		       "not a valid directory"), Op:"newcntyfs.root" }
	}
	
	CntyFS = new(ContentityFS)
	// 2025.01 Change from path to AbsFP
	CntyFS.rootAbsPath = pFPs.AbsFP // path 
	L.L.Info("Path for new os.DirFS: " + SU.Tildotted(aPath))
	// 2025.01 TODO Change from os.DirFS to os.Root.FS
	// var osRoot *os.Root 
	// osRoot, e = os.OpenRoot(path)
	// CntyFS.FS = osRoot.FS()
	println("MCFILE contentityfs_new L70 FIXME os.Root")
	CntyFS.FS = os.DirFS(pathToUse) 
	// Initialize slice & map
	CntyFS.asSlice = make([]*Contentity, 0)
	CntyFS.asMap = make(map[string]*Contentity)

	// ==================
	//    FIRST PASS
	//  Load slice & map
	// ==================
	// NOTE that rel.path "." seems to be necessary 
	// here or else really weird errors occur.
	// Note that this is the place where [CntyFS]
	// being a global singleton can cause problems. 
	e = fs.WalkDir(CntyFS.FS, ".", wfnBuildContentityTree)
	if e != nil {
		// L.L.Panic("NewCntyFS.WalkDir: " + e.Error())
		return nil, &fs.PathError { Op:"newcntyfs.walkdir",
		       Err:e, Path:aPath } 
	}
	L.L.Okay("NewCntyFS: walked OK %d nords from path %s",
		 len(CntyFS.asSlice), pathToUse)

	// Debuggery 
	for ii, cc := range CntyFS.asSlice {
	    if cc == nil {
	       L.L.Error ("OOPS, CntyFS.asSlice[%02d] is NIL", ii)
	       continue
	    }
	    /* if cc.FSItem == nil || cc.FSItem.FileMeta == nil {
	       L.L.Error("WTF, man!")
	       continue } */
	    if cc.FSItem.IsDirlike() {
	        L.L.Debug("[%02d] isDIRLIKE: AbsFP: %s",
			ii, cc.FSItem.FPs.AbsFP)
	    } else {
		L.L.Debug("[%02d] MarkupType: %s", ii, cc.RawType())
	    }
	}

	// ================================
	//        SECOND PASS
	//    Range over slice to identify
	//  parent Nords and link together
	// ================================
	for i, n := range CntyFS.asSlice {
		if i == 0 {
			continue
		}
		// Is child of root ?
		// println(">>> KOSHER? " + n.Nord.RelFP())
		if !S.Contains(n.RelFP(), FU.PathSep) {
			CntyFS.rootNord.AddKid(n)
		} else {
			itsDir := FP.Dir(n.RelFP())
			itsDir = FU.EnsureTrailingPathSep(itsDir)
			// println(n.Path, "|cnex2|", itsDir)
			var par *Contentity
			var ok bool
			// L.L.Warning("itsDir: " + itsDir)
			// L.L.Warning("theMap: %+v", CntyFS.asMap)
			// PROBLEMS HERE !
			if par, ok = CntyFS.asMap[itsDir]; !ok {
				L.L.Error("findParInMap: failed for: " +
					itsDir + " of " + n.RelFP())
				panic(n.RelFP())
			}
			/*
			if itsDir != par.AbsFP() {
				panic(itsDir + " != " + par.AbsFP())
			}
			*/
			par.AddKid(n)
		}
	}
	/* more debugging
	println("DUMP LIST")
	for _, n := range pFTFS.asSlice {
		println(n.LinePrefixString(), n.LineSummaryString())
	}
	println("DUMP MAP")
	for k, v := range pFTFS.asMap {
		fmt.Printf("%s\t:: %s %s \n", k, v.LinePrefixString(), v.LineSummaryString())
	}
	*/
	// println(SU.Gbg("=== TREE ==="))
	// CntyFS.rootNord.PrintAll(os.Stdout)
	return CntyFS, nil
}

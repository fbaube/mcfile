package mcfile

import (
        "fmt"
	"io/fs"
	"os"
	"errors"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	// FSU "github.com/fbaube/fsutils"
	L "github.com/fbaube/mlog"
)

// CntyFS is a global, which is a mistake. 
// >> var CntyFS *ContentityFS

// NewContentityFS probably should take an absolute filepath, 
// because passing in a relative filepath might cause problems.
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
	
	pFPs, e = FU.NewFilepaths(aPath)
	if e != nil {
	     	// L.L.Error("NewCntyFS: bad path: %s", aPath)
		return nil, errors.New("newcntyfs: bad path: " + aPath)
	}
	var path string
	path = aPath
	
	// We will allow a symlink here, like most std lib functions,
	// so ensure trailing slash (or OS sep) before checking for
	// existence and directoryness.
	path = FU.EnsureTrailingPathSep(path)
	if !FU.IsDirAndExists(path) {
		// L.L.Error("NewCntyFS: Not a directory: %s", path)
		return nil, errors.New("newcntyfs: not a directory: " + path)
	}
	CntyFS = new(ContentityFS)
	// 2025.01 Change from path to AbsFP
	CntyFS.rootAbsPath = pFPs.AbsFP // path 
	L.L.Info("Path for new os.DirFS: " + SU.Tildotted(path))
	// 2025.01 Change from os.DirFS to os.Root.FS
	// var osRoot *os.Root 
	// osRoot, e = os.OpenRoot(path)
	// CntyFS.FS = osRoot.FS()
	println("MCFILE contentityfs_new L61 FIXME os.Root")
	CntyFS.FS = os.DirFS(path) 
	// Initialize slice & map
	CntyFS.asSlice = make([]*Contentity, 0)
	CntyFS.asMap = make(map[string]*Contentity)

	// ==================
	//    FIRST PASS
	//  Load slice & map
	// ==================
	// NOTE that rel.path "." is necessary here 
	// or else really weird errors occur.
	// Note that this is the place where [CntyFS]
	// being a global singleton can cause problems. 
	e = fs.WalkDir(CntyFS.FS, ".", wfnBuildContentityTree)
	if e != nil {
		// L.L.Panic("NewCntyFS.WalkDir: " + e.Error())
		return nil, fmt.Errorf("newcntyfs.walkdir: %w", e)
	}
	L.L.Okay("NewCntyFS: walked OK %d nords from path %s",
		 len(CntyFS.asSlice), path)

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

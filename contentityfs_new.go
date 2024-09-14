package mcfile

import (
	"io/fs"
	"os"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	// FSU "github.com/fbaube/fsutils"
	L "github.com/fbaube/mlog"
)

// CntyFS is a global, which is a mistake. 
var CntyFS *ContentityFS

// (OBS?) NewContentityFS takes an absolute filepath. Passing 
// in a relative filepath is going to cause major problems.
// .
func NewContentityFS(aPath string, okayFilexts []string) *ContentityFS {
     	/*
	// NOTE this will fail on Windoze
	if aa,_ := FP.Abs(aPath); !S.HasPrefix(aa, FU.PathSep) {
		L.L.Error("Not an abs.FP: %s", aPath)
		return nil
	}
	*/
	var path string
	var e error 
	path, e = FP.Abs(aPath) 
	if e != nil {
		L.L.Error("NewCntyFS: bad path: %s (absolute:%s)", aPath, path) 
		return nil
	}
	if !FU.IsDirAndExists(path) {
		L.L.Error("NewCntyFS: Not a directory: %s", path)
		return nil
	}
	path = FU.EnsureTrailingPathSep(path)
	CntyFS = new(ContentityFS)
	CntyFS.rootAbsPath = path // afp.S() 
	L.L.Info("Path for new os.DirFS: " + SU.Tildotted(path))
	CntyFS.FS = os.DirFS(path) // "T/allConTypes")
	// Initialize slice & map
	CntyFS.asSlice = make([]*Contentity, 0)
	CntyFS.asMap = make(map[string]*Contentity)

	// ==================
	//    FIRST PASS
	//  Load slice & map
	// ==================
	// NOTE that rel.path "." is necessary here 
	// or else really weird errors occur.
	e = fs.WalkDir(CntyFS.FS, ".", wfnBuildContentityTree)
	if e != nil {
		L.L.Panic("NewCntyFS.WalkDir: " + e.Error())
	}
	L.L.Okay("NewCntyFS: walked OK %d nords from path %s",
		 len(CntyFS.asSlice), path)

	// DEBUG
	// L.L.Warning("CntyFS.asSlice has len: %d", len(CntyFS.asSlice))
	for ii, cc := range CntyFS.asSlice {
	    // L.L.Warning("[%d]%T...", ii, cc)
	    if cc == nil {
	       L.L.Error ("OOPS, CntyFS.asSlice[%02d] is NIL", ii)
	       continue
	    }
	    // L.L.Warning("Got here!")
	    // L.L.Warning("[%02d] %+v", ii, cc)
	    /* if cc.FSItem == nil || cc.FSItem.FileMeta == nil {
	       L.L.Error("WTF, man!")
	       continue
	    } */
	    if cc.FSItem.IsDirlike() {
	        L.L.Debug("[%02d] isDIRLIKE: AbsFP: %s", ii, cc.FSItem.FPs.AbsFP)
	        } else {
		L.L.Debug("[%02d] MarkupType: %s", ii, cc.MarkupType())
		}
	}
	L.L.Debug(" END")

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
	return CntyFS
}

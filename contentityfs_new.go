package mcfile

import (
	"io/fs"
	"os"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
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
	CntyFS = new(ContentityFS)
	CntyFS.rootAbsPath = path // afp.S() 
	L.L.Info("Path for new os.DirFS: " + path)
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
	L.L.Okay("NewCntyFS: walked os.DirFS OK: " +
		"got %d nords from path %s", len(CntyFS.asSlice), path)

	// DEBUG
	for ii, cc := range CntyFS.asSlice {
	    if cc.IsDir() {
	        L.L.Dbg("[%02d] isDIR - %s", ii, cc.FSItem.FPs.AbsFP)
	        } else {
		L.L.Dbg("[%02d] %s - %s", ii, cc.MarkupType())
		}
	}
	L.L.Dbg(" END")

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
		println(">>> KOSHER? " + n.Nord.RelFP())
		if !S.Contains(n.RelFP(), FU.PathSep) {
			CntyFS.rootNord.AddKid(n)
		} else {
			itsDir := FP.Dir(n.RelFP())
			// println(n.Path, "|cnex2|", itsDir)
			var par *Contentity
			var ok bool
			if par, ok = CntyFS.asMap[itsDir]; !ok {
				L.L.Error("findParInMap: failed for: " +
					itsDir + " of " + n.RelFP())
				panic(n.RelFP())
			}
			if itsDir != par.RelFP() {
				panic(itsDir + " != " + par.RelFP())
			}
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

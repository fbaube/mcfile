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

var pCFS *ContentityFS

// NewContentityFS takes an absolute filepath. Passing in 
// a relative filepath is going to cause major problems.
// .
func NewContentityFS(path string, okayFilexts []string) *ContentityFS {
	// NOTE this will fail on Windoze
	if aa,_ := FP.Abs(path); !S.HasPrefix(aa, FU.PathSep) {
		L.L.Error("Not an abs.FP: %s", path)
		return nil
	}
	var afp FU.AbsFilePath
	afp = FU.AbsFilePath(path)
	if !afp.DirExists() {
		L.L.Error("Not a directory: %s", path)
		return nil
	}
	pCFS = new(ContentityFS)
	pCFS.rootAbsPath = path
	L.L.Info("PATH for new os.DirFS: %s", path)
	pCFS.FS = os.DirFS(path) // "T/allConTypes")
	// Initialize slice & map
	pCFS.asSlice = make([]*Contentity, 0)
	pCFS.asMap = make(map[string]*Contentity)

	// IF path IS A FILE, THIS WILL ALL FAIL !!!!

	// FIRST PASS
	// Load slice & map
	// NOTE that a rel.path (".") is necessary
	// here or else really weird errors occur.
	e := fs.WalkDir(pCFS.FS, ".", wfnBuildContentityTree)
	if e != nil {
		L.L.Panic("mcfile.newContentityFS: " + e.Error())
	}
	L.L.Okay("FS walked OK: %d nords: %s", len(pCFS.asSlice), path)

	// DEBUG
	for _, pp := range pCFS.asSlice {
		L.L.Dbg("%s ", pp.MarkupType())
	}
	L.L.Dbg(" END")

	// SECOND PASS
	// Go down slice to identify parent nords and link together.
	for i, n := range pCFS.asSlice {
		if i == 0 {
			continue
		}
		// Is child of root ?
		if !S.Contains(n.Path(), FU.PathSep) {
			pCFS.rootNord.AddKid(n)
		} else {
			itsDir := FP.Dir(n.Path())
			// println(n.Path, "|cnex2|", itsDir)
			var par *Contentity
			var ok bool
			if par, ok = pCFS.asMap[itsDir]; !ok {
				L.L.Error("findParInMap: failed for: " + itsDir + " of " + n.Path())
				panic(n.Path)
			}
			if itsDir != par.Path() {
				panic(itsDir + " != " + par.Path())
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
	// pCFS.rootNord.PrintAll(os.Stdout)
	return pCFS
}

package mcfile

import (
	"io/fs"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	FSU "github.com/fbaube/fsutils"
	L "github.com/fbaube/mlog"
)

var pCFS *ContentityFS

// NewContentityFS is duh.
func NewContentityFS(path string, okayFilexts []string) *ContentityFS {

	pCFS = new(ContentityFS)
	// Initialize embedded baseFS
	pCFS.BaseFS = *(FSU.NewBaseFS(path))
	// println("FSU.newContentityFS:", FU.Tildotted(pCFS.BaseFS.RootAbsPath()))
	// Initialize slice & map
	pCFS.asSlice = make([]*Contentity, 0)
	pCFS.asMap = make(map[string]*Contentity)

	// IF path IS A FILE, THIS WILL ALL FAIL !!!!

	// FIRST PASS
	// Load slice & map
	if nil == pCFS.InputFS() {
		panic("OH SHIT in ctyfsnew")
	}
	e := fs.WalkDir(pCFS.InputFS(), ".", wfnBuildContentityTree)
	if e != nil {
		// panic("mcfile.newContentityFS: " + e.Error())
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

package mcfile

import (
	"fmt"
	"io/fs"
	"os"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	"github.com/fbaube/fss"
)

var pCFS *ContentityFS

// NewContentityFS is duh.
func NewContentityFS(path string, okayFilexts []string) *ContentityFS {
	pCFS = new(ContentityFS)
	// Initialize embedded baseFS
	pCFS.BaseFS = *(fss.NewBaseFS(path))
	println("fss.newContentityFS:", pCFS.BaseFS.RootAbsPath())
	// Initialize slice & map
	pCFS.asSlice = make([]*ContentityNord, 0)
	pCFS.asMap = make(map[string]*ContentityNord)

	// FIRST PASS
	// Load slice & map
	e := fs.WalkDir(pCFS.InputFS(), ".", wfnBuildContentityTree)
	if e != nil {
		panic("fss.newContentityFS: " + e.Error())
	}
	fmt.Printf("fss.newContentityFS: got %d nords \n", len(pCFS.asSlice))

	// SECOND PASS
	// Go down slice to identify parent nords and link together.
	for i, n := range pCFS.asSlice {
		if i == 0 {
			continue
		}
		// Is child of root ?
		if !S.Contains(n.Path(), FU.PathSep) {
			pCFS.rootNord.AddKid(n)
			// ON.AddKid2(pFTFS.rootNord, n)
		} else {
			itsDir := FP.Dir(n.Path())
			// println(n.Path, "|cnex2|", itsDir)
			var par *ContentityNord
			var ok bool
			if par, ok = pCFS.asMap[itsDir]; !ok {
				panic(n.Path)
			}
			if itsDir != par.Path() {
				panic(itsDir + " != " + par.Path())
			}
			par.AddKid(n)
			/*
				plk := par.LastKid()
				plk2 := plk.(*ON.Nord)
				if uintptr(unsafe.Pointer(n)) == uintptr(unsafe.Pointer(plk2)) {
					println("EQUAL!!??")
				}
				plk = n.Parent()
				plk2 = plk.(*ON.Nord)
				if uintptr(unsafe.Pointer(par)) == uintptr(unsafe.Pointer(plk2)) {
					println("EQUAL!!??")
				}
				if par.LastKid() == n {
					fmt.Printf("**** OK LINK 1??? %p,%p \n", n, par.LastKid())
				}
				if par.LastKid() != n {
					fmt.Printf("**** FAILED LINK 1??? %p,%p \n", n, par.LastKid())
				}
				if n.Parent() == par {
					fmt.Printf("**** OK LINK 2??? %p,%p \n", par, n.Parent())
				}
				if n.Parent() != par {
					fmt.Printf("**** FAILED LINK 2??? %p,%p \n", par, n.Parent())
				}
			*/
		}
	}
	/*
		println("DUMP LIST")
		for _, n := range pFTFS.asSlice {
			println(n.LinePrefixString(), n.LineSummaryString())
		}
		println("DUMP MAP")
		for k, v := range pFTFS.asMap {
			fmt.Printf("%s\t:: %s %s \n", k, v.LinePrefixString(), v.LineSummaryString())
		}
	*/
	println("=== TREE ===")
	pCFS.rootNord.PrintAll(os.Stdout)
	return pCFS
}

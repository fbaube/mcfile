package mcfile

import (
	"io/fs"
	S "strings"

	"github.com/fbaube/fss"
)

type ContentityFS struct {
	fss.BaseFS
	rootNord *Contentity
	asSlice  []*Contentity
	asMap    map[string]*Contentity // string is Rel.Path
}

// Open is a dummy function, just here to satisfy an interface.
func (p *ContentityFS) Open(path string) (fs.File, error) {
	return p. /*inputFS. */ Open(path)
}

func (p *ContentityFS) Size() int {
	return len(p.asSlice)
}

func (p *ContentityFS) RootContentity() *Contentity {
	return p.rootNord
}

func mustInitRoot() bool {
	var needsInit, didDoInit bool
	needsInit = (len(pCFS.asSlice) == 0 && len(pCFS.asMap) == 0)
	didDoInit = (len(pCFS.asSlice) > 0 && len(pCFS.asMap) > 0)
	if !(needsInit || didDoInit) {
		panic("mustInitRoot: illegal state")
	}
	return needsInit
}

// wfnBuildContentityTree is
// type WalkDirFunc func(path string, d DirEntry, err error) error
func wfnBuildContentityTree(path string, d fs.DirEntry, err error) error {
	var p *Contentity
	// ROOT NODE ?
	if mustInitRoot() {
		if path != "." {
			println("wfnBuildContentityTree: root path is not dot but instead:", path)
		}
		p = NewRootContentityNord(pCFS.RootAbsPath())
		pCFS.rootNord = p
		println("wfnBuildContentityTree: root node abs.FP:", p.AbsFP())
	} else {
		// Filter out hidden (esp'ly .git) and emacs backup.
		// Note that "/" is assumed, not os.Sep
		if S.HasPrefix(path, ".") || S.Contains(path, "/.") ||
			S.HasSuffix(path, "~") || S.Contains(path, "/.git/") {
			if !S.Contains(path, "/.git/") {
				println("Path rejected:", path)
			}
			return nil
		}
		p = NewContentity(path)
	}
	pCFS.asSlice = append(pCFS.asSlice, p)
	pCFS.asMap[path] = p
	// println("Path OK:", pN.AbsFilePath)
	return nil
}

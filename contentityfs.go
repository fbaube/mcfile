package mcfile

import (
	"errors"
	"io/fs"
	S "strings"

	FSU "github.com/fbaube/fsutils"
	L "github.com/fbaube/mlog"
)

type ContentityFS struct {
	FSU.BaseFS
	rootNord      *RootContentity
	asSlice       []*Contentity
	asMap         map[string]*Contentity // string is Rel.Path
	nFiles, nDirs int
}

/* Open is a dummy function, just here to satisfy an interface.
func (p *ContentityFS) Open(path string) (fs.File, error) {
	return p. /*inputFS. * / Open(path)
} */

func (p *ContentityFS) Size() int {
	return len(p.asSlice)
}

func (p *ContentityFS) DirCount() int {
	return p.nDirs
}

func (p *ContentityFS) FileCount() int {
	return p.nFiles
}

func (p *ContentityFS) RootContentity() *RootContentity {
	return p.rootNord
}

func (p *ContentityFS) AsSlice() []*Contentity {
	var z []*Contentity
	z = p.asSlice
	return z // p.AsSlice
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

func (p *ContentityFS) DoForEvery(stgprocsr ContentityStage) {

}

// wfnBuildContentityTree is
// type WalkDirFunc func(path string, d DirEntry, err error) error
func wfnBuildContentityTree(path string, d fs.DirEntry, err error) error {
	var p *Contentity
	var e error
	// ROOT NODE ?
	if mustInitRoot() {
		var r *RootContentity
		if path != "." {
			println("wfnBuildContentityTree: "+
				"root path is not dot but instead:", path)
		}
		if pCFS.RootAbsPath() == "" {
			panic("wfnBuildContentityTree: nil ROOT")
		}
		r, e = NewRootContentity(pCFS.RootAbsPath())
		if e != nil || r == nil {
			// panic("wfnBuildContentityTree mustInitRoot NewRootContentityNord FAILED")
			return errors.New("wfnBuildContentityTree mustInitRoot NewRootContentityNord L77")
		}
		pCFS.rootNord = r
		r.MimeType = "dir"
		r.MType = "dir"
		// println("wfnBuildContentityTree: root node abs.FP:\n\t", p.AbsFP())
		var p *Contentity
		p = ((*Contentity)(r))
		pCFS.asSlice = append(pCFS.asSlice, p)
		pCFS.asMap[path] = p
		// println("ADDED TO MAP:", path)
		pCFS.nDirs = 1
		pCFS.nFiles = 0
		return nil
	}
	// Filter out hidden, esp'ly .git & .git* ;
	// note that "/" is assumed, not os.Sep
	if S.Contains(path, "/.git") {
		return nil
	}
	if S.HasPrefix(path, ".") || S.HasPrefix(path, "_") ||
		S.Contains(path, "/.") || S.Contains(path, "/_") ||
		S.HasSuffix(path, "gtk") || S.HasSuffix(path, "gtr") ||
		S.HasSuffix(path, "~") {
		L.L.Dbg("Path rejected: " + path)
		return nil
	}
	p, e = NewContentity(path) // FP.Join(pCFS.RootAbsPath(), path))
	if p == nil || e != nil {
		// panic("nil Contentity")
		L.L.Error("Skipping this item!")
		return nil
	}
	L.L.Dbg("Directory traverser: FileType: " + p.FileType())
	nxtIdx := len(pCFS.asSlice)
	pCFS.asSlice = append(pCFS.asSlice, p)
	p.logIdx = nxtIdx
	pCFS.asMap[path] = p
	if p.IsDir() {
		pCFS.nDirs++
		// println("================DIR ========")
		p.MimeType = "dir"
		p.MType = "dir"
	} else {
		pCFS.nFiles++
	}
	// println("ADDED TO MAP:", path)
	// println("Path OK:", pN.AbsFilePath)
	return nil
}